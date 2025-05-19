package services

import (
	"sim_data_gen/daos/diggerProductBinding"
	locationDao "sim_data_gen/daos/location"
	lorryDao "sim_data_gen/daos/lorry"
	"sim_data_gen/daos/mineProduct"
	diggerProductBindingModel "sim_data_gen/models/diggerProductBinding"
	locationModel "sim_data_gen/models/location"
	. "sim_data_gen/utils"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"math"
	"strings"
	"time"
)

/*
*
  - 自动调度:

1.根据挖机的位置,自动确定它正在开采那个配矿单元。
2.根据卡车的位置,自动确定它正在和哪个挖机进行搭配。
3.根据目前的开采情况,自动确定使用了那个采矿方案。
*/
type DispatchStatusAutoSetter struct {
	ProductName       string
	cache             *cache.Cache
	diggerIds         []string
	lorryIdsInSystem  []string
	ttl               time.Duration
	targetScheduleUrl string
	accessThreshold   float64
	queryServerToBind int //是否调用真实服务进行卡车和挖机的绑定
}
type DispatchStatusAutoSetterList struct {
	SetterList map[string]*DispatchStatusAutoSetter
	Crontab    string
	Cron       *cron.Cron
}

func (t *DispatchStatusAutoSetterList) Handle() {
	for _, setter := range t.SetterList {
		setter.Handle()
	}
}

func (t *DispatchStatusAutoSetter) Handle() {
	Logger.Info("DispatchStatusAutoSetter.Execute")
	//获取产品相关的挖机
	t.diggerIds = t.getDiggerIds()
	t.lorryIdsInSystem = t.getLorryIds()
	if len(t.diggerIds) == 0 {
		Logger.Error("No digger IDs found!")
		return
	}
	if len(t.lorryIdsInSystem) == 0 {
		Logger.Error("No lorry IDs found!")
		return
	}
	//认为装车时间不会超过15分钟
	minUtc := time.Now().Add(-15 * time.Minute).Unix()
	//获取卡车位置
	lorryLocations := t.GetValidLocationsByIds(t.lorryIdsInSystem, minUtc)
	diggerLocations := t.GetValidLocationsByIds(t.diggerIds, minUtc)
	if len(lorryLocations) == 0 {
		Logger.Error("No lorry locations found!")
		return
	}

	//自动根据卡车和挖机的位置确定绑定关系
	binds := t.FindNearDiggersByLorries(lorryLocations, diggerLocations, t.accessThreshold) // 距离挖机15米则认为卡车和挖机绑定
	lorryIds := make([]string, len(binds))
	diggerIds := make([]string, len(lorryIds))
	for _, bind := range binds {
		lorryIds = append(lorryIds, bind.LorryId)
		diggerIds = append(diggerIds, bind.DiggerId)
	}
	//更新挖机的绑定关系
	if t.queryServerToBind == 1 {
		t.BindLorryDigger(lorryIds, diggerIds)
	}
	//如果卡车停留在挖机附近然后不再有数据,说明卡车在等待载货

}

type LorryAccessDigger struct {
	Distance float64
	DiggerId string
	LorryId  string
}

func (t *DispatchStatusAutoSetter) FindNearDiggersByLorries(diggers []locationModel.Location, lorries []locationModel.Location, rByMeter float64) map[string]LorryAccessDigger {
	searcher := NewThingSearch(8)
	for _, location := range lorries {
		//这里的坐标单位是百万分之一度
		searcher.AddPoint(*location.Token, *location.X, *location.Y)
	}
	nearest := make(map[string]LorryAccessDigger)
	for _, digger := range diggers {
		//这里的坐标单位是百万分之一度,米*9转换成百万分之一度
		results, _ := searcher.SearchCircle(*digger.X, *digger.Y, rByMeter*9)
		if len(results) > 0 {
			for _, result := range results {
				if HeightDiff(result, digger) {
					continue
				}
				if element, found := nearest[result.ID]; !found {
					nearest[result.ID] = LorryAccessDigger{
						Distance: result.Distance,
						DiggerId: *digger.Token,
						LorryId:  result.ID,
					}
				} else if result.Distance < element.Distance {
					nearest[result.ID] = LorryAccessDigger{
						Distance: result.Distance,
						DiggerId: *digger.Token,
						LorryId:  result.ID,
					}
				}
			}
		}
	}
	return nearest
}

func (t *DispatchStatusAutoSetter) BindLorryDigger(lorryIds, diggerIds []string) bool {
	for i, lorryId := range lorryIds {
		Logger.Info("BindLorryDigger", zap.String("lorry", lorryId), zap.String("digger", diggerIds[i]))
	}

	//如果原来已经绑定则不必重新绑定
	//进行绑定
	//POST http://localhost:8086/api/target_schedule/lorrydiggerbinding/updateObj
	//{"Enabled":1,"Level":0,"Name":"卡车 卡车1号 绑定 挖机 挖机2号","Token":"e8e899be-9b52-11ef-a660-0242ac120016","DiggerDesp":"挖机2号","TokenOfDigger":"d6b04032-9b52-11ef-a660-0242ac120016","LorryDesp":"卡车1号","TokenOfLorry":"e26ee941-9b52-11ef-a660-0242ac120016"}
	url := fmt.Sprintf("%s/api/target_schedule/lorrydiggerbinding/bindLorriesAndDiggers", t.targetScheduleUrl)
	data := map[string]interface{}{
		"lorryIds":  lorryIds,
		"diggerIds": diggerIds,
	}
	resp, _ := PostRequest(url, data, "json")
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		Logger.Error("BindLorryDigger", zap.String("url", url), zap.String("resp", resp.Status))
		return false
	}
	return true
}
func HeightDiff(result *SearchResult, digger locationModel.Location) bool {
	if result.Elevation == nil || digger.Elevation == nil {
		return false
	}
	if math.Abs(*result.Elevation-*digger.Elevation) < 7 {
		return false
	}
	return true
}

func (t *DispatchStatusAutoSetter) getDiggerIds() []string {
	if x, found := t.cache.Get("diggerIds"); found {
		return x.([]string)
	}

	Logger.Info("Fetching digger IDs...")
	ids := t.queryDiggersOfProductName()
	t.cache.Set("diggerIds", ids, t.ttl)
	return ids
}
func (t *DispatchStatusAutoSetter) getLorryIds() []string {
	if x, found := t.cache.Get("lorryIdsInSystem"); found {
		return x.([]string)
	}

	ids := t.queryLorryIds()
	t.cache.Set("lorryIdsInSystem", ids, t.ttl)
	return ids
}

func (t *DispatchStatusAutoSetter) queryDiggersOfProductName() []string {
	product, err := mineProduct.GetObjByName(&t.ProductName)
	if err != nil {
		Logger.Error("mineProduct.GetObjByName", zap.Error(err))
		return []string{}
	}
	var bindings []diggerProductBindingModel.DiggerProductBindingFully1
	//这里并没有支持缺省产品,因为这个自动设置器未来可能会弃用,因此简单实现
	bindings, err = diggerProductBinding.ListDiggerOfMineProductExt(*product.Token, false)
	if len(bindings) == 0 {
		Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
		return []string{}
	}
	Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
	var diggerIds []string
	diggerIds = make([]string, len(bindings))
	for i, binding := range bindings {
		diggerIds[i] = *binding.TokenOfDigger
	}
	return diggerIds
}

func (t *DispatchStatusAutoSetter) queryLorryIds() []string {

	Logger.Info("Fetching lorry IDs...")
	list, err := lorryDao.ListObj()
	if err != nil {
		Logger.Error("lorryDiggerBindingLogDao.ListObj", zap.Error(err))
		return []string{}
	}
	if len(list) == 0 {
		Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
		return []string{}
	}
	Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
	var ids []string
	ids = make([]string, len(list))
	for i, lorry := range list {
		ids[i] = *lorry.Token
	}
	return ids
}

func (t *DispatchStatusAutoSetter) GetValidLocationsByIds(ids []string, minUtc int64) []locationModel.Location {
	var locations []locationModel.Location
	var err error
	locations, err = locationDao.GetValidLocationsByIds(ids, minUtc)
	if err != nil {
		Logger.Error("locationDao.GetLocationsByIds", zap.Error(err))
		return locations
	}
	return locations
}

func NewDispatchStatusAutoSetterList(varNamePrx string) *DispatchStatusAutoSetterList {
	dispatchSetterList := new(DispatchStatusAutoSetterList)
	productNames := GetConfig().String(varNamePrx+".productNames", "")
	if len(productNames) == 0 {
		Logger.Info("productNames is empty,skips!")
		return dispatchSetterList
	}
	queryServerToBind := GetConfig().Int(varNamePrx+".queryServerToBind", 0)
	var targetScheduleUrl string
	if queryServerToBind == 1 {
		targetScheduleUrl = GetConfig().String(varNamePrx+".targetScheduleUrl", "")
		if targetScheduleUrl == "" {
			Logger.Info("targetScheduleUrl is empty,skips")
			return dispatchSetterList
		}
		if !strings.HasSuffix(targetScheduleUrl, "/") {
			targetScheduleUrl = targetScheduleUrl + "/"
		}
	}
	setterList := make(map[string]*DispatchStatusAutoSetter)
	accessThreshold := GetConfig().Float(varNamePrx+".accessDiggerThreshold", 15.0)
	for _, productName := range strings.Split(productNames, ",") {
		if _, ok := setterList[productName]; ok {
			Logger.Error("product name duplicated: " + productName)
			continue
		}
		setter := NewDispatchStatusAutoSetter(productName, queryServerToBind, targetScheduleUrl, accessThreshold)
		setterList[productName] = setter
	}
	dispatchSetterList.SetterList = setterList
	crontabEnabled := GetConfig().Int(varNamePrx+".crontabEnabled", 0)
	if crontabEnabled == 1 {
		if dispatchSetterList.Crontab = GetConfig().String(varNamePrx+".crontab", "*/1 * * * *"); len(dispatchSetterList.Crontab) == 0 {
			Logger.Info("Skip dispatchSetterList!")
		}
		dispatchSetterList.Cron = cron.New()
		if entityId, err := dispatchSetterList.Cron.AddFunc(dispatchSetterList.Crontab, dispatchSetterList.Handle); err != nil {
			Logger.Error("Skip StartCleanDataService for error!", zap.Error(err))
		} else {
			Logger.Info(fmt.Sprintf("StartCleanDataService with Id = %d", entityId))
		}
	}
	return dispatchSetterList
}
func NewDispatchStatusAutoSetter(productName string, queryServerToBind int, targetScheduleUrl string, accessThreshold float64) *DispatchStatusAutoSetter {
	AutoSetter := new(DispatchStatusAutoSetter)
	AutoSetter.ProductName = productName
	AutoSetter.targetScheduleUrl = targetScheduleUrl
	AutoSetter.accessThreshold = accessThreshold
	AutoSetter.ttl = time.Duration(GetConfig().Int("cache.ttl", 5)) * time.Minute
	AutoSetter.queryServerToBind = queryServerToBind
	return AutoSetter
}
