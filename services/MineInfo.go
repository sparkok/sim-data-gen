package services

import (
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"math"
	"regexp"
	boomGroupDao "sim_data_gen/daos/boomGroup"
	"sim_data_gen/daos/diggerProductBinding"
	diggerSwitchBoomGroupLogDao "sim_data_gen/daos/diggerSwitchBoomGroupLog"
	locationGnssDataDao "sim_data_gen/daos/locationGnssData"
	lorryDao "sim_data_gen/daos/lorry"
	lorryNearbyTargetSpanDao "sim_data_gen/daos/lorryNearbyTargetSpan"
	"sim_data_gen/daos/mineProduct"
	"sim_data_gen/daos/unloadSite"
	. "sim_data_gen/entity"
	boomGroupModel "sim_data_gen/models/boomGroup"
	diggerProductBindingModel "sim_data_gen/models/diggerProductBinding"
	locationGnssDataModel "sim_data_gen/models/locationGnssData"
	lorryNearbyTargetSpanModel "sim_data_gen/models/lorryNearbyTargetSpan"
	. "sim_data_gen/utils"
	"strconv"
	"strings"
	"time"
)

type SystemBaseInfo struct {
	productName                    string
	accessDiggerThreshold          float64
	accessUnloadSiteThreshold      float64
	mineProductEnableAsActive      int64
	lorryId2Name                   map[string]string
	diggerId2Name                  map[string]string
	diggerNames                    []string
	diggerIds                      []string
	boomGroupId2Info               map[string]boomGroupModel.BoomGroup
	boomGroupNames                 []string
	boomGroupIds                   []string
	lorryNames                     []string
	lorryIds                       []string
	diggerCacheTime                int64
	boomGroupCacheTime             int64
	judgeBoomGroupMethod           int //0 - (默认)通过配矿单元和挖机的对应关系判断配矿单元,1 - 是通过地理信息来判断配矿单元
	polygonSearchTree              *GaeaPolygonSearchTree
	lorryCacheTime                 int64
	diggerToWeightBridgeMaxSeconds int64 //从挖机到地磅的最大秒数
	diggerLocationCache            *cache.Cache
	diggerLocationCacheTtl         time.Duration
	unloadSiteCacheTime            int64
	UnloadSiteLocations            UnloadSiteLocation
	cleanDiggerCacheTtl            time.Duration
	diggerLoggerCacheTime          int64
	diggerLoggers                  map[string]*DiggerLogger
}

func NewMineInfo(varName, productName string, mineProductEnableAsActive int64) *SystemBaseInfo {
	diggerLocationCacheTtl := time.Duration(GetConfig().Int(varName+".diggerCacheTime.ttl", 60)) * time.Minute
	cleanDiggerCacheTtl := time.Duration(GetConfig().Int(varName+".diggerCacheTime.cleanUpInterval", 10)) * time.Minute
	diggerToWeightBridgeMaxSeconds := GetConfig().Int64(varName+".diggerToWeightBridgeMaxSeconds", 600) //卡车从挖机到地磅所需的最长时间(秒数)
	lorryAccessDiggerThreshold := GetConfig().Float(varName+".lorryAccessDiggerThreshold", 15.0)        //卡车距离挖机在lorryAccessDiggerThreshold距离内视为接近(米)
	return &SystemBaseInfo{
		productName:                    productName,
		accessDiggerThreshold:          lorryAccessDiggerThreshold,
		mineProductEnableAsActive:      mineProductEnableAsActive,
		diggerToWeightBridgeMaxSeconds: diggerToWeightBridgeMaxSeconds,
		diggerLocationCache:            cache.New(diggerLocationCacheTtl, cleanDiggerCacheTtl),
		diggerLocationCacheTtl:         diggerLocationCacheTtl,
		cleanDiggerCacheTtl:            cleanDiggerCacheTtl,
	}
}

func (t *SystemBaseInfo) getDiggerInfosOfProduct(ProductName string) (ids []string, names []string, id2Name map[string]string) {
	ids = make([]string, 0)
	names = make([]string, 0)
	id2Name = make(map[string]string)
	product, err := mineProduct.GetObjByName(&ProductName)
	if err != nil {
		Logger.Error("mineProduct.GetObjByName", zap.Error(err))
		return
	}
	includeNoBoundDigger := false
	if t.mineProductEnableAsActive == 1 && (product.Status == nil || *product.Status == PRODUCT_STATUS_ENABLED) {
		includeNoBoundDigger = true
		//0 激活的产品, 1 其他的产品
	} else if product.Status == nil || *product.Status != PRODUCT_STATUS_ENABLED {
		//0 正常, 1 关闭
		return
	}
	var bindings []diggerProductBindingModel.DiggerProductBindingFully1
	//这里并没有支持缺省产品,因为这个自动设置器未来可能会弃用,因此简单实现
	bindings, err = diggerProductBinding.ListDiggerOfMineProductExt(*product.Token, includeNoBoundDigger)
	if len(bindings) == 0 {
		Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
		return
	}
	for _, binding := range bindings {
		ids = append(ids, *binding.TokenOfDigger)
		names = append(names, *binding.DiggerDesp)
		id2Name[*binding.TokenOfDigger] = *binding.DiggerDesp
	}
	return
}
func (t *SystemBaseInfo) getBoomGroupInfos() (ids []string, names []string, id2NameMap map[string]boomGroupModel.BoomGroup) {
	Logger.Info("Fetching boomGroup IDs...")
	list, err := boomGroupDao.ListObj()
	if err != nil || len(list) == 0 {
		Logger.Error("boomGroupDao.ListObj", zap.Error(err))
		return
	}

	ids = make([]string, len(list))
	names = make([]string, len(list))
	id2NameMap = make(map[string]boomGroupModel.BoomGroup)
	for i, boomGroup := range list {
		ids[i] = *boomGroup.Token
		names[i] = *boomGroup.Name
		id2NameMap[*boomGroup.Token] = boomGroup
	}
	return
}
func (t *SystemBaseInfo) getLorryInfos() (ids []string, names []string, id2NameMap map[string]string) {
	Logger.Info("Fetching lorry IDs...")
	list, err := lorryDao.ListObj()
	if err != nil {
		Logger.Error("lorryDiggerBindingLogDao.ListObj", zap.Error(err))
		return
	}
	if len(list) == 0 {
		Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
		return
	}

	ids = make([]string, len(list))
	names = make([]string, len(list))
	id2NameMap = make(map[string]string)
	for i, lorry := range list {
		ids[i] = *lorry.Token
		names[i] = *lorry.Name
		id2NameMap[*lorry.Token] = *lorry.Name
	}
	return
}

/**
 * 这个函数应该被定时器反复调用,以便为地磅溯源数据做准备
 */
func (t *SystemBaseInfo) CacheData4SourceWeightBridge(dateFlag string, someTime int64) {
	Logger.Info("SimulateDataMaker.Execute")
	//获取产品相关的挖机
	if !t.loadInfo(dateFlag) {
		return
	}

	var minUtc = someTime - t.diggerToWeightBridgeMaxSeconds
	diggerLocationsList := t.loadDiggerLocations(minUtc, someTime)
	if diggerLocationsList == nil {
		return
	}
	//var lorryBoundByTheDigger = make(map[string]*DiggerNearby)
	for _, lorryId := range t.lorryIds {
		//获取在某个时间点之前已经被计算过的最大时间
		err := t.cacheInDB(dateFlag, someTime, lorryId, minUtc, diggerLocationsList)
		if err != nil {
			return
		}
		//至此数据已经基本添加到缓存,随时可以为地磅数据溯源
	}
}

/*
* 计算卡车在某个时间点在干啥
* returns: "" - 在路上或者其他地方;UnloadSiteX - 在卸矿点; 字符串长度>0  - 在挖机附近
 */
func (t *SystemBaseInfo) WhereIsLorry4SomeTime(lorryId string, dateFlag string, someTime int64) (string, error) {
	if !t.loadInfo(dateFlag) {
		return "", errors.New("loadInfo error")
	}
	var minUtc = someTime - t.diggerToWeightBridgeMaxSeconds
	diggerLocationsList := t.loadDiggerLocations(minUtc, someTime)
	if diggerLocationsList == nil {
		return "", errors.New("no diggers")
	}
	err := t.cacheInDB(dateFlag, someTime, lorryId, minUtc, diggerLocationsList)
	if err != nil {
		return "", errors.New("cacheInDB error")
	}
	//为卡车生成溯源数据:
	// 1.根据sometime这个时间点,卡车不再任何目标附近;卡车在卸矿点附近;卡车在挖机附近;
	// 2.查询maxSpanBeforeSomeTime这个时间点之前的一个时间段,因为时间段的安装状态划分的所以前面一个时间段必定和maxSpanBeforeSomeTime时间段状态不同
	lastNearbySpans, err := lorryNearbyTargetSpanDao.ListLatestNUntil(lorryId, dateFlag, someTime, 5)
	if err != nil {
		Logger.Error("ListLorryDuring", zap.Error(err))
		return "", errors.New("ListLorryDuring error")
	}
	Logger.Info("lastNearbySpans", zap.Any("lastNearbySpans", lastNearbySpans))
	for _, nearbySpan := range lastNearbySpans {
		switch *nearbySpan.ObjType {
		case 0: //无目标,在路上或者其他地方
			{
				return "", nil
			}
		case 1: //挖机
			{
				return *nearbySpan.NearbyObj, nil
			}
		case 2: //破碎站
			{
				return *nearbySpan.NearbyObj, nil
			}
		case 3: //失联
			{
				return "", errors.New("Offline")
			}
		case -1: //出错了
			{
				return "", errors.New("Data Error")
			}
		}
	}
	return "", errors.New("Data Error")
}

func (t *SystemBaseInfo) SourceVehicleComeFromDigger(lorryId string, dateFlag string, weighBridgeTimePoint int64) (string, error) {
	if !t.loadInfo(dateFlag) {
		return "", errors.New("loadInfo error")
	}
	var minUtc = weighBridgeTimePoint - t.diggerToWeightBridgeMaxSeconds
	diggerLocationsList := t.loadDiggerLocations(minUtc, weighBridgeTimePoint)
	if diggerLocationsList == nil {
		return "", errors.New("no diggers")
	}
	err := t.cacheInDB(dateFlag, weighBridgeTimePoint, lorryId, minUtc, diggerLocationsList)
	if err != nil {
		return "", errors.New("cacheInDB error")
	}
	//为卡车生成溯源数据:
	// 1.根据sometime这个时间点,卡车不再任何目标附近;卡车在卸矿点附近;卡车在挖机附近;
	// 2.查询maxSpanBeforeSomeTime这个时间点之前的一个时间段,因为时间段的安装状态划分的所以前面一个时间段必定和maxSpanBeforeSomeTime时间段状态不同
	lastNearbySpans, err := lorryNearbyTargetSpanDao.ListLatestNUntil(lorryId, dateFlag, weighBridgeTimePoint, 5)
	if err != nil {
		Logger.Error("ListLorryDuring", zap.Error(err))
		return "", errors.New("ListLorryDuring error")
	}
	Logger.Info("lastNearbySpans", zap.Any("lastNearbySpans", lastNearbySpans))
	for _, nearbySpan := range lastNearbySpans {
		switch *nearbySpan.ObjType {
		case 0: //无目标,在路上或者其他地方
			{
				continue
			}
		case 1: //挖机
			{
				return *nearbySpan.NearbyObj, nil
			}
		case 2: //破碎站
			{
				//这个时间点无法溯源
				return "", errors.New("Not found")
			}
		case 3: //失联
			{
				continue
			}
		case -1: //出错了
			{
				return "", errors.New("Data Error")
			}
		}
	}
	return "", errors.New("Data Error")
}

func (t *SystemBaseInfo) loadDiggerLocations(minUtc int64, someTime int64) []DiggerLocations {
	var diggerLocationsList = make([]DiggerLocations, 0)
	for i, diggerId := range t.diggerIds {
		//此次要加入缓存
		locations := t.ListDiggerDuring(diggerId, minUtc, someTime)
		diggerLocationsList = append(diggerLocationsList, NewDiggerLocations(diggerId, t.diggerNames[i], locations))
	}
	return diggerLocationsList
}

func (t *SystemBaseInfo) cacheInDB(dateFlag string, someTime int64, lorryId string, minUtc int64, diggerLocationsList []DiggerLocations) error {
	_, maxSpanBeforeSomeTime, err := t.FetchLatestNearbyInfoBeforeSomeTime(lorryId, dateFlag, someTime)
	if err != nil {
		//数据库故障
		return err
	}
	var lorryLocations []locationGnssDataModel.LocationGnssData
	if maxSpanBeforeSomeTime != someTime {
		//对于maxSpanBeforeSomeTime == someTime 则说明someTime已经被计算过了,不需要再计算了
		//获取之前没有计算过的数据
		lorryLocations, err = t.ListLorryDuring(lorryId, int64(math.Max(float64(minUtc), float64(maxSpanBeforeSomeTime))), someTime)
		if err != nil {
			Logger.Error("ListLorryDuring", zap.Error(err))
			return err
		}
	}
	lorryLocationLength := len(lorryLocations)
	if lorryLocationLength > 0 {
		var lorryLocation locationGnssDataModel.LocationGnssData
		var nearbyDiggerIds []string
		var distances []float64
		var utcs []int
		var boomGroupId string
		var found bool
		provLocationTime := 0
		followedLocationTime := 0
		provNearbyTarget := "" //"Empty"(NO_NEARBY_TARGET) - 表示不接近任何挖机, "UnloadSiteN" 接近第N个破碎站,其他字符串 - 表示接近某个挖机,
		var nearbyTarget string
		for i := 0; i < len(lorryLocations); i++ {
			lorryLocation = lorryLocations[i]
			if followedLocationTime > 0 && *lorryLocation.Utc > followedLocationTime+MAX_BREAK_SECONDS {
				//如果中间数据丢失,则视为在 followedLocationTime+DEFAULT_LOCATION_BREAK_SECONDS 和 *lorryLocation.Utc-DEFAULT_LOCATION_BREAK_SECONDS 两个时间点各收到两个失联状态数据包
				//认为 [provLocationTime,followedLocationTime+DEFAULT_LOCATION_BREAK_SECONDS]是provNearbyTarget这个状态
				provLocationTime, provNearbyTarget = t.UpdateNearbyStatus(*lorryLocation.Tid, dateFlag, provLocationTime, followedLocationTime+DEFAULT_LOCATION_BREAK_SECONDS, LORRY_IS_OFFLINE, provNearbyTarget, i == len(lorryLocations)-1)
				//这其中失联 [followedLocationTime+DEFAULT_LOCATION_BREAK_SECONDS,*lorryLocation.Utc-DEFAULT_LOCATION_BREAK_SECONDS]是失联状态
				provLocationTime, provNearbyTarget = t.UpdateNearbyStatus(*lorryLocation.Tid, dateFlag, provLocationTime, *lorryLocation.Utc-DEFAULT_LOCATION_BREAK_SECONDS, LORRY_IS_OFFLINE, provNearbyTarget, i == len(lorryLocations)-1)
			}
			//-1 表示不限制最大数目
			nearbyDiggerIds, distances, _, utcs = t.FindNearDiggerIdsByLocation(diggerLocationsList, lorryLocation, t.accessDiggerThreshold, -1, 180)
			if len(nearbyDiggerIds) == 0 {
				nearbyTarget = t.FindNearUnloadSiteIdsByLocation(lorryLocation, *lorryLocation.Tid, t.accessUnloadSiteThreshold)
				if nearbyTarget == "" {
					//卡车不接近任何挖机,也不接近任何破碎站
					provLocationTime, provNearbyTarget = t.UpdateNearbyStatus(*lorryLocation.Tid, dateFlag, provLocationTime, *lorryLocation.Utc, NO_NEARBY_TARGET, provNearbyTarget, i == len(lorryLocations)-1)
				} else {
					//卡车接近某个破碎站
					provLocationTime, provNearbyTarget = t.UpdateNearbyStatus(*lorryLocation.Tid, dateFlag, provLocationTime, *lorryLocation.Utc, nearbyTarget, provNearbyTarget, i == len(lorryLocations)-1)
				}
			} else {
				Logger.Debug("lorry is by diggers", zap.String("lorryId", lorryId), zap.Strings("diggerIds", nearbyDiggerIds))
				//when code reaches here,there are only two conditions: 1. only one digger is near lorry; 2. two diggers are near lorry.
				if len(distances) == 1 {
					_, boomGroupId, found = t.calcBoomGroupDugByDigLogger(nearbyDiggerIds[0], int64(utcs[0]))
					if !found {
						Logger.Error("calcBoomGroupDugByDigLogger", zap.String("diggerId", nearbyDiggerIds[0]), zap.Int("utc", utcs[0]))
						boomGroupId = "ErrorBoomGroupId" //写一个错误的id方便查看
					}
					//卡车要么接近一个挖机要么不接近任何挖机,其他情况都是异常情况
					provLocationTime, provNearbyTarget = t.UpdateNearbyStatus(*lorryLocation.Tid, dateFlag, provLocationTime, *lorryLocation.Utc, boomGroupId, provNearbyTarget, i == len(lorryLocations)-1)
				} else if len(distances) > 1 {
					_, boomGroupId, found = t.calcBoomGroupDugByDigLogger(nearbyDiggerIds[0], int64(utcs[0]))
					if !found {
						Logger.Error("calcBoomGroupDugByDigLogger", zap.String("diggerId", nearbyDiggerIds[0]), zap.Int("utc", utcs[0]))
						boomGroupId = "ErrorBoomGroupId" //写一个错误的id方便查看
					}
					//卡车接近两个挖机,这种情况是异常情况,为了方便就选第一个,但是要打印错误日志
					provLocationTime, provNearbyTarget = t.UpdateNearbyStatus(*lorryLocation.Tid, dateFlag, provLocationTime, *lorryLocation.Utc, boomGroupId, provNearbyTarget, i == len(lorryLocations)-1)
					Logger.Error("many diggers are near lorry", zap.String("lorryId", lorryId), zap.Strings("diggerIds", nearbyDiggerIds))
				}
			}
			followedLocationTime = *lorryLocation.Utc
		}
	}
	return nil
}

const NO_NEARBY_TARGET = "Empty"
const LORRY_IS_OFFLINE = "Offline"

func (t *SystemBaseInfo) loadBUnloadSites(loadTime int64) bool {
	if t.UnloadSiteCacheTimeOut(loadTime) {
		unloadSiteList, err := unloadSite.ListObj()
		if err != nil {
			Logger.Error("ListLorryDuring", zap.Error(err))
			return false
		}
		t.UnloadSiteLocations = NewUnloadSiteLocation(unloadSiteList)
		t.unloadSiteCacheTime = loadTime
	}
	return t.UnloadSiteLocations.Valid()
}

func (t *SystemBaseInfo) loadBoomGroups(loadTime int64) bool {
	if t.BoomGroupCacheTimeOut(loadTime) {
		t.boomGroupIds, t.boomGroupNames, t.boomGroupId2Info = t.getBoomGroupInfos()
		if t.judgeBoomGroupMethod == 1 {
			//要通过经纬度来判断配矿单元
			digBuffer := 10.0
			t.polygonSearchTree = NewSearchTreeByPolygons()
			for boomGroupId, boomGroupInfo := range t.boomGroupId2Info {
				if boomGroupInfo.Geom == nil {
					Logger.Error("boomGroupInfo.Geom is nil", zap.String("boomGroupId", boomGroupId))
					continue
				}
				t.polygonSearchTree.InsertWktAsPolygon(boomGroupId, *boomGroupInfo.Name, boomGroupInfo.Geom.Wkt, digBuffer)
			}
		}
		t.boomGroupCacheTime = loadTime
	}
	return len(t.boomGroupIds) > 0
}
func (t *SystemBaseInfo) loadInfo(dateFlag string) bool {
	loadTime := time.Now().Unix()
	if !t.loadDiggers(loadTime) {
		Logger.Error("failed to loadDiggers!")
		return false
	}
	if !t.loadDiggerLogIfNecessary(dateFlag, loadTime) {
		Logger.Error("digger logger invalid!")
		return false
	}
	if !t.loadBoomGroups(loadTime) {
		Logger.Error("failed to load boomGroups!")
		return false
	}

	if !t.loadLorryIds(loadTime) {
		Logger.Error("No lorries found!")
		return false
	}
	Logger.Debug("boomGroupNames", zap.Any("boomGroupNames", t.boomGroupNames))
	Logger.Debug("diggerId2Name", zap.Any("diggerId2Name", t.diggerId2Name))
	Logger.Debug("boomGroupId2Info", zap.Any("boomGroupId2Info", t.boomGroupId2Info))
	return true
}

func (t *SystemBaseInfo) loadLorryIds(loadTime int64) bool {
	if t.LorryCacheTimeOut(loadTime) {
		t.lorryIds, t.lorryNames, t.lorryId2Name = t.getLorryInfos()
		t.lorryCacheTime = loadTime
	}
	return len(t.lorryIds) > 0
}

func (t *SystemBaseInfo) loadDiggers(loadTime int64) bool {
	if t.DiggerCacheTimeOut(loadTime) {
		t.diggerIds, t.diggerNames, t.diggerId2Name = t.getDiggerInfosOfProduct(t.productName)
		t.diggerCacheTime = loadTime
	}
	if len(t.diggerIds) == 0 {
		Logger.Error("No diggerIds found!")
		return false
	}
	return true
}
func (t *SystemBaseInfo) DiggerCacheTimeOut(loadTime int64) bool {
	return loadTime-t.diggerCacheTime > 5.0*60
}
func (t *SystemBaseInfo) DiggerLoggerCacheTimeOut(now int64) bool {
	return now-t.diggerLoggerCacheTime > 5.0*60
}
func (t *SystemBaseInfo) BoomGroupCacheTimeOut(loadTime int64) bool {
	return loadTime-t.boomGroupCacheTime > 5.0*60
}
func (t *SystemBaseInfo) UnloadSiteCacheTimeOut(loadTime int64) bool {
	return loadTime-t.unloadSiteCacheTime > 5.0*60
}

func (t *SystemBaseInfo) LorryCacheTimeOut(loadTime int64) bool {
	return loadTime-t.lorryCacheTime > 5.0*60
}
func (t *SystemBaseInfo) ListDiggerDuring(lorryId string, beginUtc int64, endUtc int64) []locationGnssDataModel.LocationGnssData {
	var ret []locationGnssDataModel.LocationGnssData
	for i := t.calcBlockIndex(beginUtc); i <= t.calcBlockIndex(endUtc); i++ {
		ret = append(ret, t.getBlockByIndex(lorryId, i, beginUtc, endUtc)...)
	}
	return ret
}
func (t *SystemBaseInfo) ListLorryDuring(id string, utc int64, now int64) ([]locationGnssDataModel.LocationGnssData, error) {
	//这里不用加缓存,因为以前计算过的不会再载入了
	return locationGnssDataDao.ListObjDuring(id, utc, now)
}
func (t *SystemBaseInfo) FindNearDiggerIdsByLocation(diggers []DiggerLocations, lorry locationGnssDataModel.LocationGnssData, rByMeter float64, maxSize int, maxBreakSeconds int) ([]string, []float64, [][2]float64, []int) {
	ids := make([]string, 0)
	distances := make([]float64, 0)
	locations := make([][2]float64, 0)
	diggerUTCs := make([]int, 0)
	for _, digger := range diggers {
		diggerLocation, found := digger.FindLatestLocation(int64(*lorry.Utc))
		if !found {
			continue
		}
		if int(*lorry.Utc-*diggerLocation.Utc) > maxBreakSeconds {
			continue
		}

		distance := CalculateDistanceAsMeter(*diggerLocation.X, *diggerLocation.Y, *lorry.X, *lorry.Y)

		if distance < rByMeter {
			//Logger.Info(fmt.Sprintf("lorry(%s) is near digger(%s)", *lorry.Token, digger.DiggerId))
			ids = append(ids, digger.DiggerId)
			distances = append(distances, distance)
			diggerUTCs = append(diggerUTCs, *diggerLocation.Utc)
			locations = append(locations, [2]float64{*lorry.X, *lorry.Y})
		}
	}
	//对distances按从小到大排序,同时调整对应索引的ids
	if len(distances) > 1 {
		for i := 0; i < len(distances)-1; i++ {
			for j := 0; j < len(distances)-1-i; j++ {
				if distances[j] > distances[j+1] {
					distances[j], distances[j+1] = distances[j+1], distances[j]
					ids[j], ids[j+1] = ids[j+1], ids[j]
					locations[j], locations[j+1] = locations[j+1], locations[j]
					diggerUTCs[j], diggerUTCs[j+1] = diggerUTCs[j+1], diggerUTCs[j]
				}
			}
		}
	}
	if maxSize > 0 && len(ids) > maxSize {
		ids = ids[:maxSize]
		distances = distances[:maxSize]
		locations = locations[:maxSize]
		diggerUTCs = diggerUTCs[:maxSize]
	}
	return ids, distances, locations, diggerUTCs
}

/*
*
  - 函数返回的是这个状态的起始时间,如果中间状态是连续的这不改变这个时间
    forceUpdateDb - 强制更新数据库,一般来说这是最后一包才会这样
*/
// 假定一般情况下10秒一包数据
const DEFAULT_LOCATION_BREAK_SECONDS = 10
const MAX_BREAK_SECONDS = DEFAULT_LOCATION_BREAK_SECONDS * 20

func (t *SystemBaseInfo) UpdateNearbyStatus(lorryId string, dateAsStr string, provUtc int, utc int, target, provTarget string, flushData bool) (int, string) {
	if target == provTarget {
		if flushData {
			t.UpdateNearbySpan(lorryId, dateAsStr, target, provUtc, utc)
		}
		return provUtc, provTarget
	}
	//如果之前的状态和现在的状态不一样,则更新前面一个状态
	t.UpdateNearbySpan(lorryId, dateAsStr, target, provUtc, utc)
	return utc, target
}

func (t *SystemBaseInfo) FindNearUnloadSiteIdsByLocation(location locationGnssDataModel.LocationGnssData, tid string, thresholdByMeter float64) string {
	//UnloadSite0~UnloadSiteN 第一个到第一个卸矿点
	unloadSiteObj := t.UnloadSiteLocations.FindNearby(*location.X, *location.Y, thresholdByMeter)
	if unloadSiteObj != nil {
		return t.convertName2Target(*unloadSiteObj.Name)
	}
	return ""
}

func (t *SystemBaseInfo) FetchLatestNearbyInfoBeforeSomeTime(lorryId string, dateAsStr string, someTime int64) (target string, timeValue int64, err error) {
	//1.先查是否有时间段含盖someTime,如果有则返回这个时间段的target和min(someTime,时间段最大值)
	spanIncludeSomeTime, err := lorryNearbyTargetSpanDao.FindObjIncludeTime(lorryId, dateAsStr, "", int(someTime))
	if err != nil {
		return "", 0, err
	}
	if spanIncludeSomeTime != nil {
		return *spanIncludeSomeTime.NearbyObj, int64(math.Min(float64(*spanIncludeSomeTime.EndUtc), float64(someTime))), nil
	}
	//2.如果1没有查到结果,则找小于someTime的最大时间点,返回这个时间点的target和这个时间段的最大值
	spanIncludeSomeTime, err = lorryNearbyTargetSpanDao.FindLatestObjBeforeTime(lorryId, dateAsStr, int(someTime))
	if err != nil {
		return "", 0, err
	}
	if spanIncludeSomeTime != nil {
		return *spanIncludeSomeTime.NearbyObj, int64(*spanIncludeSomeTime.EndUtc), nil
	}
	return "", 0, nil
}

func (t *SystemBaseInfo) UpdateNearbySpan(lorryId string, dateAsStr string, target string, beginUtc int, endUtc int) {
	//1.取beginUtc和endUtc是否落在某个时间段内
	//2.如果落在某个时间段内,并且前后target相同,则需要延长这个时间段再进行更新
	//3.如果落在某个时间段内,并且前后target不同同,则直接增加一个新时段
	//4.如果不在某个时间段内,则直接增加一个新时段
	var spanBegin, spanEnd *lorryNearbyTargetSpanModel.LorryNearbyTargetSpan
	var err error
	spanBegin, err = lorryNearbyTargetSpanDao.FindObjIncludeTime(lorryId, dateAsStr, target, beginUtc)
	if err != nil {
		Logger.Error("FindObjIncludeTime", zap.Error(err))
	}
	spanEnd, err = lorryNearbyTargetSpanDao.FindObjIncludeTime(lorryId, dateAsStr, target, endUtc)
	if err != nil {
		Logger.Error("FindObjIncludeTime", zap.Error(err))
	}
	if spanBegin != nil && *spanBegin.NearbyObj == target {
		if spanEnd != nil && *spanEnd.NearbyObj == target {
			// 如果开始时间和结束时间都有时间段包含则三个时间段合并成一个,也就是要
			//1.更新spanBegin所在时间段;
			spanBegin.EndUtc = RefInt(*spanEnd.EndUtc)
			lorryNearbyTargetSpanDao.UpdateObj(spanBegin)
			//2.不保存本时间段;
			//3.要删除spanEnd所在时间段.
			lorryNearbyTargetSpanDao.DeleteObj(spanEnd.Token)
		} else {
			//开始时间被包含,结束时间段没有被包含
			spanBegin.EndUtc = RefInt(endUtc)
			lorryNearbyTargetSpanDao.UpdateObj(spanBegin)
		}
	} else if spanEnd != nil && *spanEnd.NearbyObj == target {
		//提前spanEnd所在时间段并更新
		spanEnd.BeginUtc = RefInt(beginUtc)
		lorryNearbyTargetSpanDao.UpdateObj(spanEnd)
	} else {
		//开始时间段和结束时间段都不被包含
		span := &lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{
			BeginUtc:  &beginUtc,
			LorryId:   &lorryId,
			EndUtc:    &endUtc,
			DateFlag:  RefString(t.calcDateFlag(int64(beginUtc))),
			NearbyObj: &target,
			ObjType:   RefInt(t.calcObjType(target)),
		}
		lorryNearbyTargetSpanDao.CreateObj(span)
	}

}

func (t *SystemBaseInfo) calcDateFlag(utc int64) string {
	return time.Unix(utc, 0).In(SystemZone).Format("2006-01-02")
}

func (t *SystemBaseInfo) calcObjType(target string) int {
	//0:'无目标',1:'挖机',2:'破碎站',3:'离线',-1 : '出错了‘
	if target == NO_NEARBY_TARGET {
		return 0
	} else if strings.HasPrefix(target, "UnloadSite") {
		return 2
	} else if len(target) > 0 {
		return 1
	} else if target == LORRY_IS_OFFLINE {
		return 3
	} else {
		Logger.Error("unknown calcObjType", zap.String("target", target))
		return -1
	}
}

func (t *SystemBaseInfo) convertName2Target(name string) string {
	// 定义正则表达式：匹配"破碎站"后跟单个数字（0-9）
	re := regexp.MustCompile(`^破碎站(\d)$`)
	if re.MatchString(name) {
		// 提取匹配的数字部分
		matches := re.FindStringSubmatch(name)
		if len(matches) > 1 {
			return "UnloadSite" + matches[1]
		}
	}
	// 默认返回UnloadSite0
	return "UnloadSite0"
}

func (t *SystemBaseInfo) genDiggerLocationKey(id string, utcFrom int64, blockSize int64) string {
	return id + "_" + strconv.FormatInt(utcFrom, 10) + "_" + strconv.FormatInt(blockSize, 10)
}

const BLOCK_SIZE = 60

func (t *SystemBaseInfo) calcBlockIndex(value int64) int64 {
	return value / BLOCK_SIZE
}

func (t *SystemBaseInfo) getBlockByIndex(lorryId string, blockIndex int64, minUtc int64, maxUtc int64) []locationGnssDataModel.LocationGnssData {
	var cacheBlock []locationGnssDataModel.LocationGnssData
	storeKey := t.genDiggerLocationKey(lorryId, blockIndex*BLOCK_SIZE, BLOCK_SIZE)
	//try to get data from cache at first
	obj, found := t.diggerLocationCache.Get(storeKey)
	if found {
		if obj != nil {
			cacheBlock = obj.([]locationGnssDataModel.LocationGnssData)
			return t.filterLocationGnssDataByTime(cacheBlock, int(minUtc), int(maxUtc))
		} else {
			return nil
		}
	}

	var err error
	//query from db
	cacheBlock, err = locationGnssDataDao.ListObjDuring(lorryId, blockIndex*BLOCK_SIZE, BLOCK_SIZE)
	if err != nil {
		Logger.Error("locationGnssDataDao.ListObjDuring", zap.Error(err))
		return nil
	}
	if cacheBlock == nil {
		cacheBlock = make([]locationGnssDataModel.LocationGnssData, 0)
	}
	t.diggerLocationCache.Set(storeKey, cacheBlock, cache.DefaultExpiration)
	//filter by time
	return t.filterLocationGnssDataByTime(cacheBlock, int(minUtc), int(maxUtc))
}

func (t *SystemBaseInfo) filterLocationGnssDataByTime(block []locationGnssDataModel.LocationGnssData, minUtc int, maxUtc int) []locationGnssDataModel.LocationGnssData {
	var result []locationGnssDataModel.LocationGnssData
	for _, v := range block {
		if *v.Utc >= minUtc && *v.Utc <= maxUtc {
			result = append(result, v)
		}
	}
	return result
}

func (t *SystemBaseInfo) CreateDiggerName2Id() map[string]string {
	ret := make(map[string]string)
	for i, v := range t.diggerIds {
		ret[v] = t.diggerNames[i]
	}
	return ret
}

/**
 * 通过挖机日志计算当时挖机正在挖那个配矿单元
 */
func (t *SystemBaseInfo) calcBoomGroupDugByDigLogger(diggerId string, unixTime int64) (string, string, bool) {
	dateFlag := t.calcDateFlag(unixTime)
	key := t.calcDiggerDateKey(dateFlag, diggerId)
	diggerLogger, ok := t.diggerLoggers[key]
	if !ok {
		return "", "", false
	}
	diggerSwitchBoomGroupLog, found := diggerLogger.FindLatest(unixTime)
	if !found {
		return "", "", false
	}
	return *diggerSwitchBoomGroupLog.DiggerId, *diggerSwitchBoomGroupLog.BoomGroupId, true
}

func (t *SystemBaseInfo) loadDiggerLogIfNecessary(dateFlag string, loadTime int64) bool {
	if t.DiggerLoggerCacheTimeOut(loadTime) {
		if t.diggerLoggers == nil {
			t.diggerLoggers = make(map[string]*DiggerLogger)
		}
		for diggerId, name := range t.diggerId2Name {
			if diggerLogger, ok := t.diggerLoggers[t.calcDiggerDateKey(dateFlag, diggerId)]; ok {
				items, err := diggerSwitchBoomGroupLogDao.ListObjBeyond(diggerId, dateFlag, diggerLogger.MaxUtc)
				if err != nil {
					Logger.Error("ListObjBeyond", zap.Error(err))
					continue
				}
				if len(items) > 0 {
					diggerLogger.AppendItemsBeyondTime(items)
				}
			} else {
				items, err := diggerSwitchBoomGroupLogDao.ListObjBeyond(diggerId, dateFlag, 0)
				if err != nil {
					Logger.Error("ListObjBeyond", zap.Error(err))
					continue
				}
				if len(items) > 0 {
					t.diggerLoggers[t.calcDiggerDateKey(dateFlag, diggerId)] = NewDiggerLogger(diggerId, name, items)
				}
			}
		}
		t.diggerLoggerCacheTime = loadTime
	}
	return len(t.diggerLoggers) > 0
}

func (t *SystemBaseInfo) calcDiggerDateKey(dateFlag string, id string) string {
	return fmt.Sprintf("%s-%s", dateFlag, id)
}

func (t *SystemBaseInfo) CreateDiggerId2Name() map[string]string {
	ret := make(map[string]string)
	for i, v := range t.diggerIds {
		ret[v] = t.diggerNames[i]
	}
	return ret
}
