package services

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"math/rand"
	. "sim_data_gen/calculate"
	boomGroupDao "sim_data_gen/daos/boomGroup"
	boomGroupInfoDao "sim_data_gen/daos/boomGroupInfo"
	lorryDao "sim_data_gen/daos/lorry"
	lorryDiggerBindDao "sim_data_gen/daos/lorryDiggerBinding"
	lorryDiggerBindingLogDao "sim_data_gen/daos/lorryDiggerBindingLog"
	mineChangeStatusDao "sim_data_gen/daos/mineChangeStatus"
	mineProductDao "sim_data_gen/daos/mineProduct"
	"sim_data_gen/daos/noAssignMineSchedule"
	productAndBoomGroupDao "sim_data_gen/daos/productAndBoomGroup"
	setOfBoomGroupsDao "sim_data_gen/daos/setOfBoomGroups"
	weighLoggerDao "sim_data_gen/daos/weighLogger"
	yAnalyserDao "sim_data_gen/daos/yAnalyser"
	. "sim_data_gen/entity"
	. "sim_data_gen/meterial"
	"sim_data_gen/mine_assignment"
	"sim_data_gen/models/boomGroup"
	boomGroupInfoModel "sim_data_gen/models/boomGroupInfo"
	lorryModel "sim_data_gen/models/lorry"
	lorryDiggerBindingModel "sim_data_gen/models/lorryDiggerBinding"
	"sim_data_gen/models/lorryDiggerBindingLog"
	mineChangeStatusModel "sim_data_gen/models/mineChangeStatus"
	noAssignMineScheduleModel "sim_data_gen/models/noAssignMineSchedule"
	productAndBoomGroupModel "sim_data_gen/models/productAndBoomGroup"
	setOfBoomGroupsModel "sim_data_gen/models/setOfBoomGroups"
	weighLoggerModel "sim_data_gen/models/weighLogger"
	yAnalyserModel "sim_data_gen/models/yAnalyser"
	. "sim_data_gen/utils"
	"strings"
	"sync"
	"time"
)

type CalculatorType int

const (
	CalculatorByTotalMass      CalculatorType = iota //根据中子仪的皮带秤和地磅的总质量相同的原理来推算
	CalculatorByLorryDetecting                       //根据一车矿石的品味相近的策略来推算
)

/*
*
根据中子仪的读数,估计当前配单单元的含量
*/
type ProductionContent struct {
	Name               string
	TokenOfMineProduct string
	//是否允许使用未绑定产品的配矿单元,当系统中同时只有一个产品处于激活状态时,将不指定配矿单元和产品的绑定关系,这时候为了能让产品使用配矿单元,则必须设置这个选项为true,否则就要手动绑定
	//无主的资源都默认属于该产品
	IsDefault                       bool
	BoomGroups                      map[string]*boomGroup.BoomGroup
	productsMd5                     string
	dataLoc                         sync.Mutex
	SampleMinutesInMineChangeStatus int //进行采样的时间间隔分钟
	WeigherUseInOrOutSite           int //地磅是使用 1 - 空车入场时间 还是 0 - 重车出场时间
	Calculator                      CalculatorType
	MaterialAttribs                 []MaterialAttrib //物质索引,例如 1,2,3
	ContentLimits                   []mine_assignment.ContentPercentLimitReq
	prevTimeOfNeutronData           int    //上一次获取中子仪数据的时间
	prevTimeOfBridgeData            string //上一次获取皮带秤数据的时间
	primaryMaterialName             string //cao的物质索引
	Vehicle2BoomGroupCalcMode       CalcBoomGroupMode
	VehicleMineFrom                 *SystemBaseInfo
}

func NewProductionContent(token string, name string, matAttribsTxt string, contentLimitsTxt string, statusSampleSeconds int, weigherUseInOrOutSite int, primaryMaterialName string, vehicle2BoomGroupCalcMode CalcBoomGroupMode, mineProductEnableAsActive int64) *ProductionContent {
	product := new(ProductionContent)
	product.Name = name
	product.WeigherUseInOrOutSite = weigherUseInOrOutSite
	product.primaryMaterialName = primaryMaterialName
	product.TokenOfMineProduct = token
	product.SampleMinutesInMineChangeStatus = statusSampleSeconds
	product.Vehicle2BoomGroupCalcMode = vehicle2BoomGroupCalcMode
	if product.Vehicle2BoomGroupCalcMode == ByDiggerBoomGroupLogger {
		product.VehicleMineFrom = NewMineInfo("vehicleMineFrom", product.Name, mineProductEnableAsActive)
		Logger.Info("use ByDiggerBoomGroupLogger")
	}

	materialAttribs := make([]MaterialAttrib, 0)
	if err := json.Unmarshal([]byte(matAttribsTxt), &materialAttribs); err != nil {
		Logger.Error("json.Unmarshal matAttribsTxt error", zap.Error(err))
		return nil
	}
	product.MaterialAttribs = materialAttribs
	ContentLimits := make([]mine_assignment.ContentPercentLimitReq, 0)
	if err := json.Unmarshal([]byte(contentLimitsTxt), &ContentLimits); err != nil {
		Logger.Error("json.Unmarshal contentLimitsTxt error", zap.Error(err))
		return nil
	}
	product.ContentLimits = ContentLimits
	return product
}

/**
* 配矿单元
 */
func (this *ProductionContent) containBoomGroup(boomGroupId string) bool {
	return true
}

/*
从数据库中刷新当前配单信息
*/
//func (this *ProductionContent) LoadBoomGroups() {
//	var request request.SearchInfo
//	request.TokenOfMineProduct = &this.TokenOfMineProduct
//	productAndBoomGroupList, err := productAndBoomGroup.ListByFilter(request.GetConditions(), "order by boom_group_desp asc")
//	if err != nil {
//		Logger.Error(fmt.Sprintf("ListByFilter %s", err.Error()))
//		return
//	}
//	md5Sum := calculateBoomGroupMd5(productAndBoomGroupList)
//	if len(md5Sum) == 0 {
//		return
//	}
//	if md5Sum == this.productsMd5 {
//		Logger.Info("the boom groups of product no change!")
//		return
//	}
//	this.productsMd5 = md5Sum
//	for _, productAndBoomGroupObj := range productAndBoomGroupList {
//		this.LoadBoomGroup(*productAndBoomGroupObj.TokenOfBoomGroup)
//	}
//}

func calculateBoomGroupMd5(list []productAndBoomGroupModel.ProductAndBoomGroupFully) string {
	bytes, err := json.Marshal(list)
	if err != nil {
		Logger.Error("failed to calculateBoomGroupMd5", zap.Error(err))
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum(bytes))
}

func (this *ProductionContent) LoadBoomGroup(tokenOfBoomGroup string) {
	boomGroupObj, err := boomGroupDao.GetObjById(&tokenOfBoomGroup)
	if err != nil {
		Logger.Error(fmt.Sprintf("ListByFilter %s", err.Error()))
		return
	}
	if this.BoomGroups == nil {
		this.BoomGroups = make(map[string]*boomGroup.BoomGroup)
	}
	this.BoomGroups[tokenOfBoomGroup] = &boomGroupObj
}

/*
* 计算配矿方案,并吧结果保存到数据库中
 */
func (this *ProductionContent) predictAndSaveMineAssignments4Product(mineAssignmentApiUrl string, dateStrOfFetchData, dateStrOfSavePlan string, diggers []mine_assignment.DiggerReq, minProbability int, forcePredict bool, isAsync bool, predictTime time.Time, calcBoomGroupMode CalcBoomGroupMode, snapName string, maxTimeAsUtc *int, debugProductName string, outputExcel *string, cacheInDb bool) {
	//如果有新数据则进行推断,并返回变化的配矿单元信息
	neutronDataTime, bridgeDataTime, hasNewData := this.ifNewDataArrive(dateStrOfFetchData, maxTimeAsUtc)
	if !forcePredict && !hasNewData {
		this.Cache4SourceWeightBridgeToday(predictTime)
		Logger.Info("no new data arrive,skip predict")
		return
	}
	startPointName := fmt.Sprintf("%s_%s_%s", this.Name, time.Unix(int64(neutronDataTime), 0).Format("150405"), bridgeDataTime)
	resultList, err := this.predict(minProbability, dateStrOfFetchData, startPointName, predictTime, snapName, maxTimeAsUtc, cacheInDb)
	if err != nil {
		Logger.Info("skip for error", zap.Error(err))
		return
	}
	materialMap4Debug := make(map[string][2]float64)
	this.saveBoomGroupInfos(resultList, func(boomGroupId string, contents [2]float64) {
		materialMap4Debug[boomGroupId] = contents
	})
	if len(snapName) > 0 {
		//SnapSvr.BindResults4Task(snapName, resultList)
		SnapSvr.SaveResult(snapName, &resultList)
		excelFilePath, _ := SnapSvr.GenerateMineReport(snapName, materialMap4Debug)
		if excelFilePath != "" && outputExcel != nil {
			*outputExcel = excelFilePath
		}
	}
	if isAsync {
		//重算配矿方案
		go func() {
			//重算配矿计划
			this.calcAllMineAssignment(mineAssignmentApiUrl, diggers, dateStrOfSavePlan)
			//更新计算时间
			this.keepDataArriveTime(neutronDataTime, bridgeDataTime)
		}()
	} else {
		//重算配矿方案
		this.calcAllMineAssignment(mineAssignmentApiUrl, diggers, dateStrOfSavePlan)
		//更新计算时间
		this.keepDataArriveTime(neutronDataTime, bridgeDataTime)
	}
}

func (this *ProductionContent) Cache4SourceWeightBridgeToday(predictTime time.Time) {
	if this.VehicleMineFrom != nil {
		dateFlag := predictTime.In(SystemZone).Format("2006-01-02")
		this.VehicleMineFrom.CacheData4SourceWeightBridge(dateFlag, predictTime.Unix())
	}
}

/*
* 给氧化钙的标准差计算数据准确度的百分比
 */
func CaoStdDev2Accuracy(meanAndStdDev [2]float64) float64 {
	ValueOfTenTime := meanAndStdDev[1]
	return math.Abs(100 - ValueOfTenTime*100)
}
func (this *ProductionContent) saveBoomGroupInfos(resultList ContentResultList, materialContentCallBack func(string, [2]float64)) {
	//更新数据到boomGroupInfo表中,从后向前更新,如何后面的单元在boomGroupExisted已经存在,这前面的就不必更新了
	resultList.EnumerateLastLorryCarries(func(boomGroupId string, results []ContentResult, contentContentAndPBeyondMinP [2]float64) {
		results = FilterLowPurityResult(results)
		if len(results) > 0 {
			//只用较高纯度的结果来计算品味
			this.saveBoomGroupInfo(boomGroupId, results, materialContentCallBack)
		}
	})
}

func (this *ProductionContent) saveBoomGroupInfo(boomGroupToken string, correctResultArray []ContentResult, materialContentCallBack func(string, [2]float64)) *boomGroupInfoModel.BoomGroupInfo {
	//1.获取dateFlag是今天的,时间值最大的数据
	//bridgeCheckTime := (*correctResultArray)[0].BridgeData.CheckTime
	//boomGroupToken := this.calcBoomGroupIdByVehicleAndTime(vehicleNo, vehicleToken, bridgeCheckTime, calcBoomGroupMode)
	//if boomGroupToken == "" {
	//	Logger.Warn("failed to find BoomGroup by vehicleNo", zap.String("vehicleNo", vehicleNo))
	//	return nil
	//}
	materialContents := make([][2]float64, len(this.MaterialAttribs))
	boomGroupInfo := &boomGroupInfoModel.BoomGroupInfo{
		TokenOfBoomGroup: RefString(boomGroupToken),
		Token:            Uuid(),
		CreatedAt:        RefTime(time.Now()),
	}

	for i := 0; i < len(this.MaterialAttribs); i++ {
		//materialContents 均值和方差
		materialContents[i] = CalculateContent(correctResultArray, this.MaterialAttribs[i])
		//if strings.ToLower(this.MaterialAttribs[i].Name) == strings.ToLower(this.primaryMaterialName) {
		//	//氧化钙的物质索引号,对氧化钙要从表里面计算概率,这样可调整的余地比较大
		//	materialContents[i][1] = CaoStdDev2Accuracy(materialContents[i])
		//}
		FillMaterialContentInResult(boomGroupInfo, this.MaterialAttribs[i], materialContents[i])
		materialContentCallBack(*boomGroupInfo.TokenOfBoomGroup, materialContents[i])
	}
	var err error
	var effectedCount int64
	var prevBoomGroupInfo boomGroupInfoModel.BoomGroupInfo
	needToReCalculateAssignPlans := false
	if prevBoomGroupInfo, err = boomGroupInfoDao.GetObjByBoomGroupToken(&boomGroupToken); err == nil {
		//进行数值比较,如果数值发生较大变化则离开触发配矿
		//为了防止品位缓慢变化，导致最终变化很大,但是在小尺度上看不出变化,因此采用根据变化值取整的方法,来判断变化是否发生
		if this.needReCalculateAssignPlan(prevBoomGroupInfo, boomGroupInfo) {
			needToReCalculateAssignPlans = true
		}
	}
	defer this.dataLoc.Unlock()
	this.dataLoc.Lock()
	effectedCount, err = boomGroupInfoDao.SaveObj(boomGroupInfo)
	if err != nil {
		needToReCalculateAssignPlans = false
		Logger.Info("failed to UpdateMaterialContentBoomGroup", zap.Any("boomGroupInfo", boomGroupInfo), zap.Int64("effectedCount", effectedCount))
		return nil
	}
	Logger.Info("succeed to UpdateMaterialContentBoomGroup", zap.Any("boomGroupInfo", boomGroupInfo))
	if !needToReCalculateAssignPlans {
		//不必重新计算配矿计划
		return nil
	}
	return boomGroupInfo
}

type CalcBoomGroupMode int

const (
	ByLoadGoodsLogger CalcBoomGroupMode = iota
	ByBindLogger
	ByDiggerBoomGroupLogger
)

func (this *ProductionContent) calcBoomGroupIdByVehicleAndTime(dateFlag string, vehicleNo string, bridgeCheckTime *string, calcBoomGroupMode CalcBoomGroupMode) string {
	switch calcBoomGroupMode {
	case ByLoadGoodsLogger:
		{
			if len(vehicleNo) == 0 {
				Logger.Warn("calcBoomGroupIdByVehicleAndTime need vehicleNo")
				return ""
			}
			//根据载货日志确定
			return mineProductDao.FindBoomGroupInfoByVehicleNo(vehicleNo, *bridgeCheckTime)
		}
	case ByBindLogger:
		{
			if len(vehicleNo) == 0 {
				Logger.Warn("calcBoomGroupIdByVehicleAndTime need vehicleToken")
				return ""
			}
			return this.calcBoomGroupIdByByBindLogger(vehicleNo, dateFlag, *bridgeCheckTime)
		}
	case ByDiggerBoomGroupLogger:
		{
			if len(vehicleNo) == 0 {
				Logger.Warn("calcBoomGroupIdByVehicleAndTime need vehicleNo")
				return ""
			}
			return this.calcBoomGroupIdByDiggerBoomGroupLogger(vehicleNo, dateFlag, *bridgeCheckTime)
		}
	default:
		Logger.Error("calcBoomGroupIdByVehicleAndTime unknown calcBoomGroupMode", zap.Int("calcBoomGroupMode", int(calcBoomGroupMode)))
		return ""
	}
}
func (this *ProductionContent) calcBoomGroupIdByDiggerBoomGroupLogger(vehicleNo, dateFlag string, bridgeCheckTime string) (boomGroupId string) {
	if this.VehicleMineFrom == nil {
		Logger.Error("calcBoomGroupIdByDiggerBoomGroupLogger need SystemBaseInfo")
		return ""
	}
	var timeAsUtc int
	var err error
	if timeAsUtc, err = timeStr2Utc(bridgeCheckTime); err != nil {
		Logger.Warn("timeStr2Utc error", zap.Error(err))
		return ""
	}
	var lorryObj lorryModel.Lorry
	if lorryObj, err = lorryDao.GetObjByName(&vehicleNo); err != nil {
		Logger.Error("lorryDao.GetObjByName error", zap.Error(err))
	}
	boomGroupId, err = this.VehicleMineFrom.SourceVehicleComeFromDigger(*lorryObj.Token, dateFlag, int64(timeAsUtc))
	if err != nil {
		Logger.Error("calcBoomGroupIdByDiggerBoomGroupLogger error", zap.Error(err))
		return ""
	}
	return boomGroupId
}
func timeStr2Utc(timeAsStr string) (int, error) {
	if timeObj, err := time.ParseInLocation("2006-01-02 15:04:05", timeAsStr, SystemZone); err != nil {
		return 0, err
	} else {
		return int(timeObj.Unix()), nil
	}
}

func randomSelectBoomGroup(boomGroups []productAndBoomGroupModel.ProductAndBoomGroupInDetailFully) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	index := r.Intn(len(boomGroups))
	return *boomGroups[index].TokenOfBoomGroup
}
func (this *ProductionContent) ExtractFromBoomGroupInfo(boomGroupInfo boomGroupInfoModel.BoomGroupInfo) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, materialAttrib := range this.MaterialAttribs {
		if materialAttrib.Required != nil && !*materialAttrib.Required {
			continue
		}
		val := GetMaterialContentInBoomGroupInfo(materialAttrib, boomGroupInfo)
		if val < 0 {
			Logger.Error("Material Content is error!")
			return nil
		}
		ret[strings.ToLower(materialAttrib.Name)] = val
	}
	ret["name"] = *boomGroupInfo.TokenOfBoomGroup
	return ret
}
func (this *ProductionContent) ExtractFromBoomGroup(boomGroup productAndBoomGroupModel.ProductAndBoomGroupInDetailFully) map[string]interface{} {
	ret := make(map[string]interface{})
	//配矿单元
	for _, materialAttrib := range this.MaterialAttribs {
		if materialAttrib.Required != nil && !*materialAttrib.Required {
			continue
		}
		val := GetMaterialContentInBoomGroup(materialAttrib, boomGroup)
		if val < 0 {
			Logger.Error("Material Content is error!")
			return nil
		}
		//进行配矿时要求物质名要求必须消息
		ret[strings.ToLower(materialAttrib.Name)] = val
	}
	ret["min"] = 0
	if boomGroup.BoomGroup.High == nil {
		Logger.Error("High is nil!")
		return nil
	}
	if boomGroup.BoomGroup.Used == nil {
		Logger.Error("Used is nil!")
		return nil
	}
	if *boomGroup.BoomGroup.High < *boomGroup.BoomGroup.Used {
		return nil
	}
	ret["max"] = *boomGroup.BoomGroup.High - *boomGroup.BoomGroup.Used
	ret["name"] = *boomGroup.TokenOfBoomGroup
	return ret
}

/*
*
分析每分钟的矿
weightRange 本区间的代表了从总重量 weightRange[0]吨到weightRange[1]吨的范围
根据数据发现一分钟基本上是16吨左右,所以有可能是1辆车或两辆车的混合
*/
func (this *ProductionContent) calcMineEveryMinute(startMineStatus *mineChangeStatusModel.MineChangeStatus, weightRange []float64, record yAnalyserModel.YAnalyser) error {
	//startMineStatus.DataTimeOfBridge是上次计算的地磅数据时间, startMineStatus.DataTimeOfNeutron 是上次计算的中子仪时间
	var weightLoggers []weighLoggerModel.WeighLogger
	prevWeightOfBridge := *startMineStatus.TotalMassOfBridge
	for _, item := range weightLoggers {
		prevWeightOfBridge += *item.NetWeight
	}
	return nil
}

//	func (this *ProductionContent) updateRecentMineChangeInDb(lastTime *time.Time, prevTotalMass float64) bool {
//		dataListFromNetronProbe, err := yAnalyserDao.ListObjSince(this.Name, lastTime.Unix())
//		if err != nil || len(dataListFromNetronProbe) <= 1 {
//			Logger.Error("the length of data list < 2")
//			return true
//		}
//		//中子仪每分钟一条数据,所以要搜集这些数据计算最后一条数据所代表的时刻的总质量,但是
//		//要记住这个时刻的数据是统计了一分钟的平均数据,同时单位是吨/小时
//		totalMass := prevTotalMass
//		for _, yAnalyserRecord := range dataListFromNetronProbe {
//			totalMass += *yAnalyserRecord.Flux / 60.0
//		}
//
//		return false
//	}
var ERROR_NULL_TO_STR = fmt.Errorf("converting NULL to string is unsupported")

/**
 * 从数据库中最近一个符合质量标准的状态记录
 */
func (this *ProductionContent) getLastMineChangeByDateStr(predictTime time.Time, dateStr string, minProbability int, cacheInDb bool) (changeStatus *mineChangeStatusModel.MineChangeStatus, firstOfTheDay bool) {
	var err error
	var dataTime *time.Time
	var mineChangeStatus *mineChangeStatusModel.MineChangeStatus
	if cacheInDb {
		//注意第一条记录一点是一个高质量记录,因为一天的开始没有数据干扰
		if dataTime, err = mineChangeStatusDao.GetLastTimeOfDay(this.Name, dateStr, minProbability); err == nil {
			//数据库今天有最后记录,则获取记录
			mineChangeStatus, _ = mineChangeStatusDao.GetObjByDataTime(this.Name, dateStr, dataTime)
			return mineChangeStatus, false
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			Logger.Error("getLastMineChangeByDateStr is error", zap.Error(err))
			return nil, false
		}
	}
	//如果没有之前的记录,说明今天一一条记录都没有,这是后要做一个假数据,这个数据比第一个有效的中子仪和地磅数据都早1秒钟,这样用这个状态去(xxx,xxx)查就能查到全部的有效数据,并且不会查到无效的数据
	//create it if the first record don't exist.
	if mineChangeStatus = this.createProductFirstChangeStatus(predictTime, dateStr, cacheInDb); mineChangeStatus == nil {
		//尚不满足创建当天首条记录的条件,返回nil
		return nil, false
	}
	return mineChangeStatus, true
}

//func (this *ProductionContent) produceMineChangeByDateStr(dateStr string) *mineChangeStatusModel.MineChangeStatus {
//	var err error
//	var dataTime *time.Time
//	var prevMineChangeStatus *mineChangeStatusModel.MineChangeStatus = nil
//	//1.获取最近一次的mineChangeStatus数据
//	if dataTime, err = mineChangeStatusDao.GetLastTimeOfDay(this.Name, dateStr); err == nil {
//		//数据库今天有最后记录,则获取记录
//		prevMineChangeStatus, _ = mineChangeStatusDao.GetObjByDataTime(dataTime)
//	}
//	//如果没有之前的记录,说明今天一一条记录都没有
//	this.productMineChangeStatus(prevMineChangeStatus, dateStr)
//	return prevMineChangeStatus
//}

//func (this *ProductionContent) productMineChangeStatus(prevMineChangeStatus *mineChangeStatusModel.MineChangeStatus, dateStr string) *mineChangeStatusModel.MineChangeStatus {
//	//productName 和 analyserNum 有相同的值
//	var err error
//	var neutronStartWeight float64
//	var bridgeStartWeight float64
//
//	var neutronTimes [2]time.Time
//	var bridgeTimes [2]time.Time
//	var onPredictFunc = func(result *ContentResult) {
//		Logger.Info(fmt.Sprintf("calcContents4ProductOfTheDay: the purity of mine carried by lorry(%s) is %f", result.VehicleNo, result.Purity))
//	}
//	neutronTimes[0], _ = time.Parse("2006-01-02", dateStr)
//	neutronTimes[1] = neutronTimes[0].Add(time.Hour*24 - time.Second*1)
//	neutronTimes[0] = neutronTimes[0].Add(time.Hour * 5)
//	bridgeTimes[0] = neutronTimes[0]
//	bridgeTimes[1] = neutronTimes[1]
//	neutronStartWeight = 0
//	bridgeStartWeight = 0
//	if prevMineChangeStatus != nil {
//		//如果没有前一次的查询时间,则从今天5点算起,当天24点结束
//		neutronTimes[0] = *prevMineChangeStatus.DataTimeOfNeutron
//		bridgeTimes[0] = *prevMineChangeStatus.DataTimeOfBridge
//		neutronStartWeight = *prevMineChangeStatus.TotalMassOfNeutron
//		bridgeStartWeight = *prevMineChangeStatus.TotalMassOfBridge
//	}
//
//	var yAnalyserRecords []yAnalyserModel.YAnalyser
//	yAnalyserRecords, err = yAnalyserDao.GetDataListDuring(this.Name, int(neutronTimes[0].Unix()), int(neutronTimes[1].Unix()))
//	if err != nil {
//		Logger.Error("failed to GetEarliestTimeExistedMineOfDay from yAnalyserDao", zap.Error(err))
//		return nil
//	}
//
//	var weighLoggers []weighLoggerModel.WeighLogger
//	weighLoggers, err = weighLoggerDao.GetDataListDuring(this.Name, &bridgeTimes[0], &bridgeTimes[1])
//	if err != nil {
//		Logger.Error("failed to GetEarliestTimeExistedMineOfDay from weighLoggerDao", zap.Error(err))
//		return nil
//	}
//	yAnalyserRecordSize := len(yAnalyserRecords)
//	weighLoggersSize := len(weighLoggers)
//	if yAnalyserRecordSize <= 0 || weighLoggersSize <= 0 {
//		Logger.Error("yAnalyserRecordSize <= 0 || weighLoggersSize <= 0")
//		return nil
//	}
//	var neutronWeightIndex int = 0
//	var bridgeWeightIndex int = 0
//
//	if neutronWeightIndex >= yAnalyserRecordSize || bridgeWeightIndex >= weighLoggersSize {
//		return nil
//	}
//	//把中子仪和地磅数据看成放在一起的已经排序了的数组,然后依次读取
//	weighAscentArray := NewTotalMassAssignmentCalculator(&yAnalyserRecords, neutronStartWeight, &weighLoggers, bridgeStartWeight)
//	var lastMineChangeStatusRecord *mineChangeStatusModel.MineChangeStatus
//	//TODO: 这里有一个问题,随着时间的推移,这种质量累积,误差必然越来越大,因此最终要找到一种根据质量的特性进行总质量对齐的方案
//	// 比如每隔1个小时,根据品位的序列,重新估计一下当前的原点开始的位置
//	results, lastMineChangeStatusRecord := weighAscentArray.DealWithData()
//	weighAscentArray.Predict()
//	if len(results) < 1 {
//		return nil
//	}
//	for i, result := range results {
//		lastMineChangeStatusRecord = this.saveMineChangeStatusRecord(result, dateStr, (*weighAscentArray.NeutronDataBoost)[i], (*weighAscentArray.BridgeDataBoost)[i])
//		onPredictFunc(&result)
//	}
//	return lastMineChangeStatusRecord
//}

//	func (this *ProductionContent) saveMineChangeStatusRecord(result ContentResult, dateStr string, neutronBoost *NeutronBoost, bridgeBoost *WeighLoggerBoost) *mineChangeStatusModel.MineChangeStatus {
//		var mineChangeStatusRecord = new(mineChangeStatusModel.MineChangeStatus)
//		mineChangeStatusRecord.DataTimeOfNeutron = RefTime(time.Unix(int64(*result.NeutronData.TestAt), 0))
//		mineChangeStatusRecord.CreatedAt = RefTime(time.Now())
//		mineChangeStatusRecord.Token = Uuid()
//		mineChangeStatusRecord.Changing = RefString("")
//		mineChangeStatusRecord.TotalMassOfNeutron = RefFloat64(neutronBoost.TotalMass)
//		mineChangeStatusRecord.Name = RefString(result.Name)
//		mineChangeStatusRecord.Changing = RefString("")
//		mineChangeStatusRecord.DataTimeOfBridge = RefTime(EverOk(time.Parse("2006-01-02 15:04:05", *result.BridgeData.CheckTime)))
//		//因为时间是一分钟的,所以需要除以60
//		mineChangeStatusRecord.TotalMassOfBridge = RefFloat64(bridgeBoost.TotalMass)
//		mineChangeStatusRecord.ProductName = RefString(this.Name)
//		mineChangeStatusRecord.DateFlag = &dateStr
//		if _, err := mineChangeStatusDao.CreateObj(mineChangeStatusRecord); err != nil {
//			Logger.Error("failed to CreateObj")
//			return nil
//		}
//		return mineChangeStatusRecord
//	}

func calcContents4ProductOfTheDay(calcStartPoint *mineChangeStatusModel.MineChangeStatus, productName string, dateStr string, maxTimeAsUtc *int) (results []ContentResult, status *mineChangeStatusModel.MineChangeStatus, err error) {
	//here, we need to create more records.
	//search data in yAnalyser from (dataTime,now)
	var lastTodayAnalyserRecord yAnalyserModel.YAnalyser
	var lastTodayWebLogger weighLoggerModel.WeighLogger
	var analyserRecords []yAnalyserModel.YAnalyser
	var bridgeRecords []weighLoggerModel.WeighLogger
	//获取最近一次的中子仪数据,如果需要用截至时间来限制则需要传入截止时间
	lastTodayAnalyserRecord, err = yAnalyserDao.GetLastTimeExistedMineOfDay(productName, dateStr, maxTimeAsUtc)
	if err != nil {
		return
	}
	//获取今天最近的一包地磅数据
	lastTodayWebLogger, err = weighLoggerDao.GetLatestTimeExistedMineOfDay(productName, dateStr, maxTimeAsUtc)
	if err != nil {
		return
	}
	neutronDataChanged := int64(*lastTodayAnalyserRecord.TestAt) != (*calcStartPoint.DataTimeOfNeutron).Unix()
	bridgeDataChanged := *lastTodayWebLogger.CheckTime != calcStartPoint.DataTimeOfBridge.Format("2006-01-02 15:04:05")
	if !calculable(neutronDataChanged, bridgeDataChanged) {
		results = make([]ContentResult, 0)
		Logger.Info("the condition of calculation is not matched!")
		return
	}

	analyserRecords, err = yAnalyserDao.GetDataListDuring(productName, int((*(calcStartPoint.DataTimeOfNeutron)).Unix()), *lastTodayAnalyserRecord.TestAt, false)
	if err != nil {
		Logger.Info("yAnalyserDao.GetDataListDuring error")
		return
	}

	var lastTimeOfBridge time.Time
	lastTimeOfBridge, err = time.ParseInLocation("2006-01-02 15:04:05", *lastTodayWebLogger.CheckTime, SystemZone)
	if err != nil {
		Logger.Error("failed to parse bridge time", zap.Error(err), zap.String("time", *lastTodayWebLogger.CheckTime))
		return
	}
	bridgeRecords, err = weighLoggerDao.GetDataListDuring(productName, calcStartPoint.DataTimeOfBridge, &lastTimeOfBridge, false)
	if err != nil {
		Logger.Info("weighLoggerDao.GetDataListDuring error")
		return
	}
	//存储数据中间状态
	mixList := NewTotalMassAssignmentCalculator(&analyserRecords, *calcStartPoint.TotalMassOfNeutron, &bridgeRecords, *calcStartPoint.TotalMassOfBridge)
	return mixList.Predict()
}

func calculable(neutronDataChanged bool, bridgeDataChanged bool) bool {
	if neutronDataChanged && bridgeDataChanged {
		return true
	}
	return false
}

/*
获取认为已经判断准确的矿石品位数据
产品名称 - productName
采样时间间隔 - sampleMinutesInMineChangeStatus
数据的最小可信度 - minProbability
是否保存到模拟数据表,而不是正式表 saveInMockTable
*/
//(correctResults map[string]*[]ContentResult, err error)
func (this *ProductionContent) predict(minProbability int, dateStrOfFetchData string, startPointName string, predictTime time.Time, snapName string, maxTimeAsUtc *int, cacheInDb bool) (correctResults ContentResultList, err error) {
	//获取今天字符串表示
	var mineChange *mineChangeStatusModel.MineChangeStatus
	//获取当天的最近一个大于等于minProbability可信度的起始状态。
	//这里的MineChange是指一个可信的可以用于计算的起始状态
	//我们计算时不需要再考虑这个起点状态之前的数据,这个起始状态(MineChange)是某次计算留下来的。
	//针对不同的计算要求,会考虑使用不同的MineChange作为起点,进行高可信度的计算就必须从一个高可信度的起点开始计算。
	//minProbability 就是对这个起点可信度的最低要求,可信度存储在dataType字段中 0 是最低可信度,90是最大的可信度
	if mineChange, _ = this.getLastMineChangeByDateStr(predictTime, dateStrOfFetchData, minProbability, cacheInDb); mineChange == nil {
		//如果一包都不存在,则说明数据条件还不满足
		//当天的第一包必须设置成高质量数据,一个是因为这时数据无干扰,另外一个是必须能在各种质量下都至少能找到一包数据
		err = fmt.Errorf("No Data")
		return
	}
	var results []ContentResult
	correctResults = ContentResultList{}
	//根据某天的数据进行品位推测
	results, mineChange, err = calcContents4ProductOfTheDay(mineChange, this.Name, dateStrOfFetchData, maxTimeAsUtc)
	carryId := 0
	latestVehicleNo := ""
	latestGroupFinished := false
	mineInALorry := make([]ContentResult, 0)
	if err != nil {
		return correctResults, err
	}
	primaryMaterialAttrib := this.GetPrimaryMaterialAttrib()
	var mainMaterialContentAndProbability [2]float64
	for i, result := range results {
		carryId, latestVehicleNo, latestGroupFinished = CalcCarryIdByVehicleNo(latestVehicleNo, *result.VehicleNo, carryId)
		result.CarryId = carryId
		result.Index = i
		boomGroupId := this.calcBoomGroupIdByVehicleAndTime(dateStrOfFetchData, *result.VehicleNo, result.BridgeData.CheckTime, this.Vehicle2BoomGroupCalcMode)
		if len(boomGroupId) == 0 {
			Logger.Warn("failed to calcBoomGroupIdByVehicleAndTime", zap.String("vehicleNo", *result.VehicleNo), zap.String("checkTime", *result.BridgeData.CheckTime))
			boomGroupId = "unknown"
		}
		correctResults.AddResult(boomGroupId, result)
		mineInALorry = append(mineInALorry, result)
		if latestGroupFinished || i == len(results)-1 {
			//计算mineInALorry的加权平均品位和概率
			if primaryMaterialAttrib != nil {
				mainMaterialContentAndProbability = this.CalculatePrimaryContent(*primaryMaterialAttrib, mineInALorry)
			} else {
				mainMaterialContentAndProbability = [2]float64{-1, -1}
			}
			correctResults.AddWholeLorryMine(boomGroupId, *result.VehicleNo, mineInALorry, mainMaterialContentAndProbability)
			mineInALorry = make([]ContentResult, 0)
		}
	}
	if mineChange != nil {
		//如果并没有出现一个新的车牌号,则那么我们不能假定当前这车矿已经检测完,这时候mineChange返回的是null
		mineChange.Token = RefString(startPointName)
		mineChange.CreatedAt = RefTime(predictTime)
		mineChange.Name = RefString(fmt.Sprintf("%d辆车品位确定", correctResults.Size()))
		mineChange.ProductName = RefString(this.Name)
		mineChange.Changing = RefString(correctResults.Results2Txt())
		mineChange.DateFlag = RefString(dateStrOfFetchData)
		if cacheInDb {
			if _, err = mineChangeStatusDao.SaveObj(mineChange); err != nil {
				Logger.Error("failed to CreateObj")
				return
			}
		}
	}
	return
}

/**
 * 第三个参数表示是否这个上组结果已经完毕,可以计算
 */
func CalcCarryIdByVehicleNo(latestVehicleNo string, currentVehicleNo string, lastCarryId int) (int, string, bool) {
	if latestVehicleNo != currentVehicleNo {
		lastGroupEnd := len(latestVehicleNo) > 0
		return lastCarryId + 1, currentVehicleNo, lastGroupEnd
	}
	return lastCarryId, currentVehicleNo, false
}

func (this *ProductionContent) createProductFirstChangeStatus(predictTime time.Time, dateStr string, cacheInDb bool) *mineChangeStatusModel.MineChangeStatus {
	//当天的第一个状态,要用今天最早的中子仪和地磅数据创建
	//productName 和 analyserNum 有相同的值
	var err error
	var yAnalyserRecord yAnalyserModel.YAnalyser
	yAnalyserRecord, err = yAnalyserDao.GetEarliestTimeExistedMineOfDay(this.Name, dateStr)
	if err != nil {
		if err.Error() != "nothing" {
			Logger.Error("failed to GetEarliestTimeExistedMineOfDay from yAnalyserDao", zap.Error(err))
		}
		return nil
	}

	var weighLogger weighLoggerModel.WeighLogger
	weighLogger, err = weighLoggerDao.GetEarliestTimeExistedMineOfDay(this.Name, dateStr)
	if err != nil {
		if err.Error() != "nothing" {
			Logger.Error("failed to GetEarliestTimeExistedMineOfDay from weighLoggerDao", zap.Error(err))
		}
		return nil
	}
	var mineChangeStatusRecord mineChangeStatusModel.MineChangeStatus
	//今天第一包有质量的中子仪的时间
	const firstDataOfTodayBeforeSeconds int = 1
	//中子仪的时间也向前60秒
	mineChangeStatusRecord.DataTimeOfNeutron = RefTime(time.Unix(int64(*yAnalyserRecord.TestAt-firstDataOfTodayBeforeSeconds), 0))
	var dataTimeOfBridge time.Time
	dataTimeOfBridge, err = time.ParseInLocation("2006-01-02 15:04:05", *weighLogger.CheckTime, SystemZone)
	if err != nil {
		Logger.Error("failed to parse weighLogger.CheckTime", zap.Error(err))
		return &mineChangeStatusRecord
	}

	//每天第一包是为了查询下一包做准备的,所以要特殊处理,将实际的第一包地磅数据和中子仪数据都向前1分钟
	mineChangeStatusRecord.DataTimeOfBridge = RefTime(dataTimeOfBridge.Add(time.Second * time.Duration(-firstDataOfTodayBeforeSeconds)))
	mineChangeStatusRecord.CreatedAt = RefTime(predictTime)
	mineChangeStatusRecord.Token = RefString(fmt.Sprintf("%s_FIRST_%s_%s_%s", this.Name, dateStr, mineChangeStatusRecord.DataTimeOfNeutron.Format("150405"), dataTimeOfBridge.Format("150405")))
	mineChangeStatusRecord.ProductName = RefString(this.Name)
	mineChangeStatusRecord.DateFlag = &dateStr
	mineChangeStatusRecord.Name = RefString("最早")
	mineChangeStatusRecord.Changing = RefString("")
	mineChangeStatusRecord.DataType = RefInt(mineChangeStatusModel.HIGH_PROBABILITY)
	//因为时间是一分钟的,所以需要除以60
	mineChangeStatusRecord.TotalMassOfNeutron = RefFloat64(0) //RefFloat64(*yAnalyserRecord.Flux / 60.0)
	mineChangeStatusRecord.TotalMassOfBridge = RefFloat64(0)
	if cacheInDb {
		if _, err = mineChangeStatusDao.CreateObj(&mineChangeStatusRecord); err != nil {
			Logger.Error("failed to CreateObj")
			return nil
		}
	}
	return &mineChangeStatusRecord
}

func (this *ProductionContent) needReCalculateAssignPlan(prevBoomGroupInfo boomGroupInfoModel.BoomGroupInfo, boomGroupInfo *boomGroupInfoModel.BoomGroupInfo) bool {
	if prevBoomGroupInfo.CreatedAt == nil || boomGroupInfo.CreatedAt == nil {
		return true
	}
	//如果超过10分钟也会重新配矿
	if boomGroupInfo.CreatedAt.Sub(*prevBoomGroupInfo.CreatedAt) > time.Minute*10 {
		//这样就可以消除因为微小变化导致的累计误差
		return true
	}
	for i := 0; i < len(this.MaterialAttribs); i++ {
		prevContent := boomGroupInfoModel.GetMaterialOfBoomGroupInfo(this.MaterialAttribs[i].Index, prevBoomGroupInfo)
		currentContent := boomGroupInfoModel.GetMaterialOfBoomGroupInfo(this.MaterialAttribs[i].Index, *boomGroupInfo)
		//有任何一个物质触发了重新配矿阈值,都会导致重新配矿
		if math.Abs(prevContent-currentContent) > this.MaterialAttribs[i].Threshold {
			return true
		}
	}
	return false
}

func (this *ProductionContent) calcAllMineAssignment(url string, diggers []mine_assignment.DiggerReq, dateStrOfSavePlan string) {
	diggerCount := len(diggers)
	if diggerCount == 0 {
		Logger.Warn("no diggers,so no result.")
		return
	}
	//1.获取boomGroup表中的所有配矿单元的原始化验信息,这个信息基本不变
	//TODO: 要检查数据是否和想象的一致
	boomGroupFirstContents, err := productAndBoomGroupDao.ListEnabledObjByProductToken(this.TokenOfMineProduct, this.IsDefault)
	if err != nil {
		return
	}

	if diggerCount > len(boomGroupFirstContents) {
		Logger.Warn("diggers outnumber boomGroupFirstContents,so no result.")
		return
	}

	boomGroupIds := make([]string, len(boomGroupFirstContents))
	for i, boomGroupInfo := range boomGroupFirstContents {
		Logger.Info("boomGroupInfo", zap.Any("name", boomGroupInfo.BoomGroupDesp))
		if this.IsDefault && boomGroupInfo.TokenOfMineProduct == nil {
			//如果是默认方案,则需要把 mine_product.token 和 mine_product.name 的值进行修复
			boomGroupInfo.TokenOfMineProduct = RefString(this.TokenOfMineProduct)
			boomGroupInfo.MineProductDesp = RefString("")
			boomGroupFirstContents[i] = boomGroupInfo
		}
		boomGroupIds[i] = *boomGroupInfo.TokenOfBoomGroup
	}
	//获取动态数据
	//3.根据id获取爆堆的品位
	dynamicContents, err := boomGroupInfoDao.ListLastContentsByIds(boomGroupIds)
	if err != nil {
		//此时一条记录也没有
		return
	}
	Logger.Info("dynamicContents", zap.Any("dynamicContents", dynamicContents))

	//2.查看这些爆堆是否有配矿单元相关的动态信息
	NewProductionContentCalculator(boomGroupFirstContents, dynamicContents).
		Calculate(url, this.ContentLimits, diggers,
			func(boomGroup productAndBoomGroupModel.ProductAndBoomGroupInDetailFully) map[string]interface{} {
				return this.ExtractFromBoomGroup(boomGroup)
			},
			func(boomGroupInfo boomGroupInfoModel.BoomGroupInfo) map[string]interface{} {
				return this.ExtractFromBoomGroupInfo(boomGroupInfo)
			}, func(result bool, boomGroupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, Results []mine_assignment.ResultData) {
				//先获取当前正在执行的方案,理论上最多只能有一个正在执行的方案
				selectedAssignMineResult := setOfBoomGroupsDao.QuerySelectedAssignMineResult(&this.TokenOfMineProduct)
				Logger.Info("calcAllMineAssignment", zap.Bool("result", result), zap.Any("Results", Results))
				now := time.Now()
				now = Change2TheDay(now, dateStrOfSavePlan)
				for i, assignMineResult := range Results {
					boomGroupIds := extractBoomGroupIds(assignMineResult)
					diggerIds := extractDiggerIds(diggers)
					token := genToken(this.TokenOfMineProduct, boomGroupIds, diggerIds, dateStrOfSavePlan)
					if selectedAssignMineResult != nil && *selectedAssignMineResult.Token == token {
						//如果是被选中的方案
						setOfBoomGroups := &setOfBoomGroupsModel.SetOfBoomGroups{
							BoomGroupIds:       RefString(boomGroupIds),
							Diggers:            RefString(diggerIds), //为了查询方便,直接将所有的id拼接起来,并且收尾都有 | 符号
							MatContents:        RefString(extractContents(assignMineResult.ContentPercents)),
							TokenOfMineProduct: RefString(this.TokenOfMineProduct),
							Name:               RefString(fmt.Sprintf("%s-%d", this.Name, i)),
							Nt:                 RefString(genRemark(boomGroupsMap, assignMineResult, diggers, 1)),
							Token:              RefString(token),
							UpdateAt:           RefTime(now),
							DateFlag:           &dateStrOfSavePlan,
						}
						_, err = setOfBoomGroupsDao.UpdateObjItem(setOfBoomGroups)
						continue
					} else {
						//fmt.Printf("result %v", result)
						//status 0:'已选',1:'未选',2:'作废'
						setOfBoomGroups := &setOfBoomGroupsModel.SetOfBoomGroups{
							BoomGroupIds:       RefString(boomGroupIds),
							CreatedAt:          RefTime(now),
							Diggers:            RefString(diggerIds), //为了查询方便,直接将所有的id拼接起来,并且收尾都有 | 符号
							MatContents:        RefString(extractContents(assignMineResult.ContentPercents)),
							TokenOfMineProduct: RefString(this.TokenOfMineProduct),
							Name:               RefString(fmt.Sprintf("%s-%d", this.Name, i)),
							Nt:                 RefString(genRemark(boomGroupsMap, assignMineResult, diggers, 1)),
							Token:              RefString(token),
							Status:             RefInt(1),
							UpdateAt:           RefTime(now),
							DateFlag:           &dateStrOfSavePlan,
						}
						_, err = setOfBoomGroupsDao.SaveObj(setOfBoomGroups)
					}
					if err != nil {
						Logger.Error("failed to SaveObj", zap.Error(err))
					}
				}
			})
}
func (this *ProductionContent) keepDataArriveTime(timeOfNeutronData int, prevTimeOfBridgeData string) {
	this.prevTimeOfNeutronData = timeOfNeutronData
	this.prevTimeOfBridgeData = prevTimeOfBridgeData
}
func (this *ProductionContent) ifNewDataArrive(dateStr string, maxTimeAsUtc *int) (int, string, bool) {
	//here, we need to create more records.
	//search data in yAnalyser from (dataTime,now)
	var lastTodayAnalyserRecord yAnalyserModel.YAnalyser
	var lastTodayWebLogger weighLoggerModel.WeighLogger
	var err error
	//获取最近一次的中子仪数据
	lastTodayAnalyserRecord, err = yAnalyserDao.GetLastTimeExistedMineOfDay(this.Name, dateStr, maxTimeAsUtc)
	if err != nil {
		return -1, "", false
	}
	//获取今天最近的一包地磅数据
	lastTodayWebLogger, err = weighLoggerDao.GetLatestTimeExistedMineOfDay(this.Name, dateStr, maxTimeAsUtc)
	if err != nil {
		return -1, "", false
	}
	//如果是prevMineChange当天的第一条数据,那么创建它的时候使用的DataTimeOfNeutron时间是今天最早的中子仪数据时间
	//如果和当前的最后一天时间一样,说明中子仪数据还没有更新
	neutronDataChanged := *lastTodayAnalyserRecord.TestAt > this.prevTimeOfNeutronData
	bridgeDataChanged := *lastTodayWebLogger.CheckTime > this.prevTimeOfBridgeData
	if neutronDataChanged == false && bridgeDataChanged == false {
		return *lastTodayAnalyserRecord.TestAt, *lastTodayWebLogger.CheckTime, false
	}
	return *lastTodayAnalyserRecord.TestAt, *lastTodayWebLogger.CheckTime, true
}

func (this *ProductionContent) CalculatePrimaryContent(materialAttrib MaterialAttrib, correctResultArray []ContentResult) (ret [2]float64) {
	correctResultArray = FilterLowPurityResult(correctResultArray)
	ret = CalculateContent(correctResultArray, materialAttrib)
	//ret[1] = CaoStdDev2Accuracy(ret)
	return
}

func (this *ProductionContent) GetPrimaryMaterialAttrib() *MaterialAttrib {
	var ret MaterialAttrib
	for _, materialAttrib := range this.MaterialAttribs {
		if strings.ToLower(materialAttrib.Name) == strings.ToLower(this.primaryMaterialName) {
			ret = materialAttrib
			break
		}
	}
	return &ret
}

// 返回:未被任何产品关联的绑定的配矿单元(包括enabled=false的也算绑定)
func (this *ProductionContent) calcUnboundBoomGroupIds() []string {

	boomGroups, err1 := boomGroupDao.ListObj()
	if err1 != nil {
		return nil
	}
	unboundBoomGroupIds := make(map[string]bool, len(boomGroups))
	for _, item := range boomGroups {
		unboundBoomGroupIds[*item.Token] = true
	}

	productAndBoomGroups, err2 := productAndBoomGroupDao.ListObj()
	if err2 != nil {
		return nil
	}
	for _, item := range productAndBoomGroups {
		unboundBoomGroupIds[*item.Token] = false
	}
	ret := make([]string, 0)
	for k, v := range unboundBoomGroupIds {
		if v == true {
			ret = append(ret, k)
		}
	}
	return ret
}

func (this *ProductionContent) calcBoomGroupIdByByBindLogger(vehicleNo, dateFlag string, bridgeCheckTime string) string {
	//根据绑定日志确定
	//1.先查询历史日志
	var bindingLog lorryDiggerBindingLog.LorryDiggerBindingLog
	var err error
	var timeAsUtc int
	if timeAsUtc, err = timeStr2Utc(bridgeCheckTime); err != nil {
		Logger.Warn("timeStr2Utc error", zap.Error(err))
		return ""
	}
	bindingLog, err = lorryDiggerBindingLogDao.FindObjBySpan4Name(dateFlag, &vehicleNo, timeAsUtc)
	if err == nil {
		return *bindingLog.TokenOfBoomGroup
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		Logger.Warn("FindObjBySpan error", zap.Error(err))
		return ""
	}

	//2.如果历史日志没有，则根据绑定日志确定,只是这时候要忽略enable标志,因为要管理enable标志那就太麻烦了,先简单处理
	var LorryDiggerBinding lorryDiggerBindingModel.LorryDiggerBindingFully
	LorryDiggerBinding, err = lorryDiggerBindDao.SearchMatchedDiggerByLorryName(vehicleNo)
	if err != nil {
		Logger.Warn("SearchMatchedDiggerByToken error", zap.Error(err))
		return ""
	}
	var schedule noAssignMineScheduleModel.NoAssignMineSchedule

	schedule, err = noAssignMineSchedule.GetObjByDiggerId(LorryDiggerBinding.TokenOfDigger)
	if err != nil {
		Logger.Warn("failed to fin schedule", zap.Error(err))
		return ""
	}
	if schedule.TokenOfBoomGroup == nil {
		Logger.Warn("schedule.TokenOfBoomGroup is nil")
		return ""
	}
	return *schedule.TokenOfBoomGroup
}

/*
*
mode 0 - 文本模式,1 - json模式,2 - vue模式
*/
func genRemark(groupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, resultData mine_assignment.ResultData, diggers []mine_assignment.DiggerReq, mode int) string {
	switch mode {
	case 0: //文本模式
		return genTextRemark(groupsMap, resultData, diggers)
	case 2: //vue 模式
		return genVueRemark(groupsMap, resultData, diggers)
	default:
		return genJsonRemark(groupsMap, resultData, diggers)
	case 1: //json模式
		return genJsonRemark(groupsMap, resultData, diggers)
	}
}

func genJsonRemark(groupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, resultData mine_assignment.ResultData, diggers []mine_assignment.DiggerReq) string {
	sb := strings.Builder{}
	if len(diggers) == 2 {
		return Pairs2JsonStr(
			DiggerBoomGroupPair{diggers[0],
				[]mine_assignment.ElementBoomGroupOccupation{
					resultData.BoomGroupOccupations[0],
				}, []string{extractGroupName(groupsMap, resultData.BoomGroupOccupations[0].Group)}},
			DiggerBoomGroupPair{diggers[1],
				[]mine_assignment.ElementBoomGroupOccupation{
					resultData.BoomGroupOccupations[1],
					resultData.BoomGroupOccupations[2],
				}, []string{extractGroupName(groupsMap, resultData.BoomGroupOccupations[1].Group),
					extractGroupName(groupsMap, resultData.BoomGroupOccupations[2].Group)}},
		)
	}
	return sb.String()
}

func extractGroupName(groupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, boomGroupToken string) string {
	if item, ok := groupsMap[boomGroupToken]; ok {
		return *item.BoomGroupDesp
	}
	return ""
}

type DiggerBoomGroupPair struct {
	DiggerItem     mine_assignment.DiggerReq
	BoomGroupItems []mine_assignment.ElementBoomGroupOccupation
	BoomGroupNames []string
}

//	func DiggerWithBoomGroup1(diggerItem mine_assignment.DiggerReq,boomGroupItems []mine_assignment.ElementBoomGroupOccupation) JsonPair{
//		return DiggerBoomGroupPair{
//			DiggerItem:     diggerItem,
//			BoomGroupItems: boomGroupItems,
//		}
//	}
func Pairs2JsonStr(pairs ...DiggerBoomGroupPair) string {
	theMap := map[string]any{}
	for _, pair := range pairs {
		theMap[pair.DiggerItem.Token] = map[string]interface{}{
			"name":       pair.DiggerItem.Name,
			"boomGroups": genBoomGroups(pair.BoomGroupItems, pair.BoomGroupNames),
		}
	}
	if bin, err := json.Marshal(theMap); err == nil {
		return string(bin)
	}
	return ""
}

func genBoomGroups(boomGroupOccupation []mine_assignment.ElementBoomGroupOccupation, names []string) map[string]interface{} {
	boomGroupMap := map[string]interface{}{}
	for i, item := range boomGroupOccupation {
		boomGroupMap[item.Group] = map[string]interface{}{
			"name":       names[i],
			"occupation": item.Occupation,
			"token":      item.Group,
		}
	}
	return boomGroupMap
}

func genVueRemark(groupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, data mine_assignment.ResultData, diggers []mine_assignment.DiggerReq) string {
	return ""
}

func genTextRemark(groupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, resultData mine_assignment.ResultData, diggers []mine_assignment.DiggerReq) string {
	changeLine := "\n"
	sb := strings.Builder{}
	sb.WriteString("配矿单元:")
	sb.WriteString(changeLine)
	if len(groupsMap) == 3 && len(diggers) == 2 {
		//2个文件,3个爆堆的情况:是指第一个挖机使用第一个爆堆,第二个挖机使用第二和第三个爆堆
		sb.WriteString(diggers[0].Name)
		sb.WriteString("开采:")
		sb.WriteString(changeLine)
		sb.WriteString(calcGroupRemark(groupsMap, diggers[0].Name, resultData.BoomGroupOccupations[0]))
		sb.WriteString(changeLine)
		sb.WriteString(changeLine)
		sb.WriteString(diggers[1].Name)
		sb.WriteString("开采:")
		sb.WriteString(fmt.Sprintf("比例: %2.f : %2.f", resultData.BoomGroupOccupations[1].Occupation, resultData.BoomGroupOccupations[2].Occupation))
		sb.WriteString(changeLine)
		sb.WriteString(calcGroupRemark(groupsMap, diggers[1].Name, resultData.BoomGroupOccupations[1]))
		sb.WriteString(changeLine)
		sb.WriteString(calcGroupRemark(groupsMap, diggers[1].Name, resultData.BoomGroupOccupations[2]))
		sb.WriteString(changeLine)
	}
	return sb.String()
}

func calcGroupRemark(groupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, diggerName string, occupation mine_assignment.ElementBoomGroupOccupation) string {
	if boomGroupInfo, ok := groupsMap[occupation.Group]; ok {
		return fmt.Sprintf("%s(%02.f吨)", *boomGroupInfo.BoomGroupDesp, occupation.Occupation)
	} else {
		Logger.Error("boomGroupInfo not found", zap.String("group", occupation.Group))
		return fmt.Sprintf("%s[出错](%02.f吨)", occupation.Group, occupation.Occupation)
	}
}

func extractBoomGroupIds(result mine_assignment.ResultData) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i, boomGroupOccupation := range result.BoomGroupOccupations {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString("\"")
		sb.WriteString(boomGroupOccupation.Group)
		sb.WriteString("\"")
	}
	sb.WriteString("]")
	return sb.String()
}

func genToken(productToken, groupIds, diggerIds, dateStr string) string {
	sb := strings.Builder{}
	sb.WriteString(productToken)
	sb.WriteString("-")
	sb.WriteString(groupIds)
	sb.WriteString("-")
	sb.WriteString(diggerIds)
	sb.WriteString("-")
	sb.WriteString(dateStr)
	return fmt.Sprintf("%x", md5.Sum(([]byte)(sb.String())))
}
func extractContents(contentPercents []mine_assignment.ElementContentPercent) string {
	rst, err := json.Marshal(contentPercents)
	if err != nil {
		Logger.Error("extractContents", zap.Error(err))
		return ""
	}
	return string(rst)
}
func extractDiggerIds(diggers []mine_assignment.DiggerReq) string {
	sb := strings.Builder{}
	for _, digger := range diggers {
		sb.WriteString("|")
		sb.WriteString(digger.Token)
	}
	sb.WriteString("|")
	return sb.String()
}
func extractDiggerNames(diggers []mine_assignment.DiggerReq) string {
	sb := strings.Builder{}
	for _, digger := range diggers {
		sb.WriteString("|")
		sb.WriteString(digger.Name)
	}
	sb.WriteString("|")
	return sb.String()
}
