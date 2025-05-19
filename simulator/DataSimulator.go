package simulator

import (
	. "sim_data_gen/calculate"
	boomGroupDao "sim_data_gen/daos/boomGroup"
	lorryDiggerBindingLogDao "sim_data_gen/daos/lorryDiggerBindingLog"
	weighLoggerDao "sim_data_gen/daos/weighLogger"
	yAnalyserDao "sim_data_gen/daos/yAnalyser"
	"sim_data_gen/models/lorryDiggerBindingLog"
	. "sim_data_gen/models/weighLogger"
	. "sim_data_gen/models/yAnalyser"
	. "sim_data_gen/utils"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"
)

type DataSimulator struct {
	DataDir           string
	TimeRate          float64
	NeutronDataOffset int64
	WeightDataOffset  int64
	clearTodayData    bool   // 是否清除当天数据
	workMode          int    // 工作模式: 0:不模拟，1:定时发模拟数据, 2:载入测试数据进行一次性测试
	neutronDelay      int    //中子仪数据延迟
	bridgeDelay       int    //地磅数据延迟
	ProductName       string //产品名称
}

var dataSim *DataSimulator

func GetDataSimulator() *DataSimulator {
	return dataSim
}
func StartSimulator(varName string) bool {
	simulatorWorkMode := GetConfig().Int(varName+".workModel", 0)
	if simulatorWorkMode == 0 {
		return false
	}

	//作为读取并写入模拟数据的工具
	timeRate := GetConfig().Float(varName+".timeRate", 1)
	productName := GetConfig().String(varName+".ProductName", "cp1")
	dayStr := GetConfig().String(varName+".dayStr", "2024-11-10")
	neutronDataOffset := GetConfig().Int64(varName+".neutronDataOffset", 0)
	weightDataOffset := GetConfig().Int64(varName+".neutronDataOffset", 0)
	clearTodayData := GetConfig().Int64(varName+".clearTodayData", 1)
	dataDir := GetConfig().String(varName+".dataDir", "E:\\work\\gitee\\dev\\we-mine-digger-assistant\\weight-data-statistic\\")
	dataSim = &DataSimulator{
		DataDir:           dataDir,
		TimeRate:          timeRate,
		NeutronDataOffset: neutronDataOffset,
		WeightDataOffset:  weightDataOffset,
		clearTodayData:    clearTodayData == 1,
		workMode:          simulatorWorkMode,
		neutronDelay:      25,  //中子仪数据延迟
		bridgeDelay:       150, //地磅数据延迟
		ProductName:       productName,
	}
	go dataSim.Start(timeRate, dayStr)
	return true
}

// var PRODUCT_NAME = RefString("cp1")
var wg sync.WaitGroup

func (this *DataSimulator) Start(timeRate float64, dataStr string) {
	if this.workMode == 1 {
		neutronData := ProcessNeutronGaugeFile(this.ProductName, this.DataDir, convertShortDateFormat(dataStr))
		weighData, _ := ProcessWeighGaugeFile(this.ProductName, this.DataDir, convertShortDateFormat(dataStr))
		this.saveMockData2DbRepeatedly(this.ProductName, timeRate, neutronData, weighData)
	}
}

func saveBindingLoggers2Db(loggerMap map[string][]lorryDiggerBindingLog.LorryDiggerBindingLog) {
	for _, loggers := range loggerMap {
		for _, logger := range loggers {
			_, err := lorryDiggerBindingLogDao.CreateObj(&logger)
			if err != nil {
				Logger.Error("保存绑定记录日志失败", zap.Error(err))
			}
		}
	}
}

/*
*
用数据快照来填充数据库
*/
func (this *DataSimulator) FlushDbFromSnap(productToken string, dateStr string, startNeutronTime, endNeutronTime, startWeighTime, endWeighTime time.Time, productName string, clearDataFirst, isAsync bool) {
	//查到对应的模拟数据
	neutronData := ProcessNeutronGaugeFile(productName, this.DataDir, convertShortDateFormat(dateStr))
	//地磅数据里面要保存配矿单元的序号(需要测试人员手动填写,字段名叫做配矿单元序号)
	weighData, bindingLoggerMap := ProcessWeighGaugeFile(productName, this.DataDir, convertShortDateFormat(dateStr))
	//把符合条件的时间写入到数据库中
	this.saveAllMockData2DbAtOnce(productToken, dateStr, startNeutronTime, endNeutronTime, startWeighTime, endWeighTime, neutronData, weighData, bindingLoggerMap, productName, clearDataFirst, isAsync)
}

func (this *DataSimulator) saveMockData2DbRepeatedly(productName string, timeRate float64, neutronData []YAnalyser, weighData []WeighLogger) {
	today := time.Now().Format("2006-01-02")
	if this.clearTodayData {
		_, _ = yAnalyserDao.ClearDataOfDay(productName, today)
		_, _ = weighLoggerDao.ClearDataOfDay(productName, today)
		_, _ = lorryDiggerBindingLogDao.ClearDataOfDay(today)
	}
	wg.Add(1)
	go this.SaveNeutronData(today, neutronData, timeRate)
	wg.Add(1)
	go this.SaveWeighData(today, weighData, timeRate)
	wg.Wait()
}
func (this *DataSimulator) saveAllMockData2DbAtOnce(productToken, theDay string, startNeutronTime, endNeutronTime, startWeighTime, endWeighTime time.Time, neutronData []YAnalyser, weighData []WeighLogger, bindingLoggerMap map[string][]lorryDiggerBindingLog.LorryDiggerBindingLog, productName string, clearDataFirst, isAsync bool) {
	if clearDataFirst {
		_, _ = yAnalyserDao.ClearDataOfDay(productName, theDay)
		_, _ = weighLoggerDao.ClearDataOfDay(productName, theDay)
		//要清除被模拟那一天的绑定记录日志
		_, _ = lorryDiggerBindingLogDao.ClearDataOfDay(theDay)
	}
	if isAsync {
		var waitGroup sync.WaitGroup
		waitGroup.Add(3)
		go func() {
			defer waitGroup.Done()
			this.SaveSpanNeutronData(theDay, neutronData, startNeutronTime, endNeutronTime)
		}()
		go func() {
			defer waitGroup.Done()
			this.SaveSpanWeighData(theDay, weighData, startWeighTime, endWeighTime)
		}()
		go func() {
			defer waitGroup.Done()
			saveBindingLoggers2Db(bindingLoggerMap)
		}()

		//当异步任务执行结束后再返回
		waitGroup.Wait()
	} else {
		this.SaveSpanNeutronData(theDay, neutronData, startNeutronTime, endNeutronTime)
		this.SaveSpanWeighData(theDay, weighData, startWeighTime, endWeighTime)
		saveBindingLoggers2Db(bindingLoggerMap)
	}
}

func ProcessNeutronGaugeFile(productName, dataDir, fileName string) []YAnalyser {
	neutronReader := WeightDataReader{SourceType: SourceNeutronGauge}
	neutronData := []YAnalyser{}

	if err := neutronReader.ProcessFile(dataDir, fmt.Sprintf("%s.xlsx", fileName)); err == nil {
		for neutronIndex, neutronItem := range neutronReader.NeutronData {
			dataTime, err := time.ParseInLocation("01/02/06 15:04", neutronItem.Time, SystemZone)
			//dataTime, err := time.ParseInLocation("2006-01-02 15:04", neutronItem.Time, BeijingZone)
			if err != nil {
				Logger.Error("failed to parse dataTime", zap.String("dataTime", neutronItem.Time))
				continue
			}
			analyser, err := NewYAnalyserWithRealTime(productName, dataTime, neutronIndex, neutronItem)
			if err != nil {
				continue
			}
			neutronData = append(neutronData, analyser)
		}
	} else {
		Logger.Error("Failed to process neutron file: %v", zap.Error(err))
	}
	return neutronData
}
func GetTokenOfGroup(boomGroupName string, mapping map[string]*string) *string {
	mappingptr := mapping[boomGroupName]
	if mappingptr != nil {
		return mappingptr
	}
	item, err := boomGroupDao.GetObjByName(&boomGroupName)
	if err != nil {
		Logger.Error("failed to get boom group", zap.String("boomGroupName", boomGroupName))
		return nil
	}
	mapping[boomGroupName] = RefString(*item.Token)
	return item.Token
}
func AddNewMockDiggerBindingIfNeed(dataTime time.Time, boomGroupName, lorryName string, lorryName2Binding map[string][]lorryDiggerBindingLog.LorryDiggerBindingLog, name2TokenMap map[string]*string) {
	dateFlag := dataTime.Format("2006-01-02")
	//如果当前时间已经有绑定记录了，就不需要再添加了
	bindingLoggerObj := lorryDiggerBindingLog.LorryDiggerBindingLog{
		TokenOfBoomGroup: GetTokenOfGroup(boomGroupName, name2TokenMap),
		DateFlag:         RefString(dateFlag),
		BoomGroupName:    RefString(boomGroupName),
		TokenOfDigger:    nil,
		DiggerName:       RefString(""),
		EndUtc:           RefInt(int(dataTime.Unix())),
		TokenOfLorry:     nil,
		LorryName:        RefString(lorryName),
		StartUtc:         RefInt(int(dataTime.Unix())),
		Token:            Uuid(),
	}

	if bindings := lorryName2Binding[lorryName]; bindings == nil {
		//如果卡车在模拟数据中第一次切换则需要加一个开始时间
		lorryName2Binding[lorryName] = []lorryDiggerBindingLog.LorryDiggerBindingLog{bindingLoggerObj}
	} else if *bindingLoggerObj.BoomGroupName != *bindings[len(bindings)-1].BoomGroupName {
		//如果绑定的配矿单元变了则需要重新添加一个记录,并设置上一个记录的结束时间
		bindings[len(bindings)-1].EndUtc = RefInt((int)(dataTime.Unix() - 1))
		lorryName2Binding[lorryName] = append(lorryName2Binding[lorryName], bindingLoggerObj)
	} else {
		//如果绑定的配矿单元不变则需要重新设置最近一个记录的结束时间
		bindings[len(bindings)-1].EndUtc = RefInt((int)(dataTime.Unix()))
	}
}
func ProcessWeighGaugeFile(productName string, dataDir, fileName string) ([]WeighLogger, map[string][]lorryDiggerBindingLog.LorryDiggerBindingLog) {
	weighReader := WeightDataReader{SourceType: SourceWeighBridge}
	var weighData []WeighLogger
	bindingLoggers := map[string][]lorryDiggerBindingLog.LorryDiggerBindingLog{}
	var name2TokenMap = make(map[string]*string)
	if err := weighReader.ProcessFile(dataDir, fmt.Sprintf("%s.xlsx", fileName)); err == nil {
		for bridgeIndex, weighItem := range weighReader.BridgeData {
			dataTime, err := time.ParseInLocation("2006/01/02 15:04:05", weighItem.TareTime, SystemZone)
			if err != nil {
				Logger.Error("failed to parse dataTime", zap.String("dataTime", weighItem.TareTime))
				continue
			}
			AddNewMockDiggerBindingIfNeed(
				dataTime,
				weighItem.BoomGroupName,
				weighItem.VehicleNo,
				bindingLoggers,
				name2TokenMap,
			)
			weighLogger := NewBridgeDataWithRealTime(productName, dataTime, bridgeIndex, weighItem.NetWeight, weighItem.VehicleNo)
			weighData = append(weighData, weighLogger)
		}
	} else {
		Logger.Error("Failed to process weigh file: %v", zap.Error(err))
	}
	return weighData, bindingLoggers
}
func convertShortDateFormat(standardDateFormat string) string {
	dateVal, err := time.ParseInLocation("2006-01-02", standardDateFormat, SystemZone)
	if err != nil {
		return ""
	}
	return dateVal.Format("060102")
}

func (this *DataSimulator) SaveSpanNeutronData(theDay string, neutronData []YAnalyser, startFrom, endTo time.Time) {
	//这里假定中子仪的延迟时间20秒
	//数据时间就取原始时间但是日期要换掉,并且最后一个点设置为当前时间减去延迟时间,其他数据表向前延迟
	for _, data := range neutronData {
		testAt := time.Unix(int64(*data.TestAt), 0)
		//testAt = Change2TheDay(testAt, theDay)
		//filter out data which is out of range
		if testAt.Before(startFrom) || testAt.After(endTo) {
			continue
		}
		data.CreatedAt = RefTime(testAt.Add(time.Second * time.Duration(this.neutronDelay)))
		data.TestAt = RefInt(int(testAt.Unix()))
		Logger.Info("SaveNeutronData", zap.Time("testAt", testAt))
		if _, err := yAnalyserDao.SaveObj(&data); err != nil {
			Logger.Error(fmt.Sprintf("yAnalyser SaveObj %s", err.Error()))
		}
	}
}

func (this *DataSimulator) SaveSpanWeighData(theDay string, weighData []WeighLogger, startFrom, endTo time.Time) {
	const delaySeconds = 60 //这里假定地磅的延迟时间60秒
	//数据时间就取原始时间但是日期要换掉,并且最后一个点设置为当前时间减去延迟时间,其他数据表向前延迟
	//startFrom = startFrom.Add(-time.Minute*25 - time.Second*delaySeconds)
	//endTo = endTo.Add(-time.Minute*10 - time.Second*delaySeconds)
	for _, data := range weighData {
		checkTime, err := time.ParseInLocation("2006-01-02 15:04:05", *data.CheckTime, SystemZone)
		if err != nil {
			Logger.Error("failed to parse CheckTime", zap.String("CheckTime", *data.CheckTime))
			continue
		}
		checkTime = Change2TheDay(checkTime, theDay)
		if checkTime.Before(startFrom) || checkTime.After(endTo) {
			continue
		}
		data.CheckTime = RefString(checkTime.Format("2006-01-02 15:04:05"))
		data.UpdateAt = RefTime(checkTime.Add(time.Duration(delaySeconds) * time.Second))
		if _, err := weighLoggerDao.SaveObj(&data); err != nil {
			Logger.Error(fmt.Sprintf("weighLogger SaveObj %s", err.Error()))
		}
	}
}

//func calcDstSubSrcInDay(src time.Time, dst time.Time) time.Duration {
//	year, month, day := src.Date()
//	theStartOfSrc := time.Date(year, month, day, 0, 0, 0, 0, src.Location())
//	dst.Date()
//	year, month, day = dst.Date()
//	theStartOfDst := time.Date(year, month, day, 0, 0, 0, 0, dst.Location())
//	return theStartOfDst.Sub(theStartOfSrc)
//}

func (this *DataSimulator) SaveNeutronData(theDay string, neutronData []YAnalyser, timeRate float64) {
	var lastTime = 0
	var delaySeconds = 20 //这里假定中子仪的延迟时间20秒
	var offset int64 = 0
	for _, data := range neutronData {
		if offset < this.WeightDataOffset {
			offset++
			continue
		}
		currentTime := *data.TestAt
		Logger.Info("SaveNeutronData", zap.Int("currentTime", currentTime))
		if lastTime != 0 {
			interval := time.Duration(float64(currentTime-lastTime) / timeRate)
			Logger.Info("SaveNeutronData TimeInterval", zap.Duration("interval", interval))
			time.Sleep(interval * time.Second)
		}
		now := time.Now()
		data.CreatedAt = RefTime(now)
		data.TestAt = RefInt(int(now.Unix()) - delaySeconds)
		if _, err := yAnalyserDao.SaveObj(&data); err != nil {
			Logger.Error(fmt.Sprintf("yAnalyser SaveObj %s", err.Error()))
		}
		lastTime = currentTime
	}
	defer wg.Done()
}

func (this *DataSimulator) SaveWeighData(theDay string, weighData []WeighLogger, timeRate float64) {
	var lastTime = 0
	var delaySeconds = 60 //这里假定中子仪的延迟时间20秒
	var offset int64 = 0
	for _, data := range weighData {
		if offset < this.WeightDataOffset {
			offset++
			continue
		}
		checkTime, err := time.ParseInLocation("2006-01-02 15:04:05", *data.CheckTime, SystemZone)
		if err != nil {
			Logger.Error("failed to parse CheckTime", zap.String("CheckTime", *data.CheckTime))
			continue
		}
		currentTime := int(checkTime.Unix())
		Logger.Info("SaveWeighData", zap.Int("currentTime", currentTime))
		if lastTime != 0 {
			interval := time.Duration(float64(currentTime-lastTime) / timeRate)
			Logger.Info("SaveWeighData TimeInterval", zap.Duration("interval", interval))
			time.Sleep(interval * time.Second)
		}
		now := time.Now()
		now.Add(time.Duration(-1*delaySeconds) * time.Second)
		data.CheckTime = RefString(now.Format("2006-01-02 15:04:05"))
		if _, err := weighLoggerDao.SaveObj(&data); err != nil {
			Logger.Error(fmt.Sprintf("weighLogger SaveObj %s", err.Error()))
		}
		lastTime = currentTime
	}
	defer wg.Done()
}

func NewYAnalyserWithRealTime(productName string, dateTime time.Time, neutronIndex int, neutronItem NeutronGaugeData) (YAnalyser, error) {
	var data YAnalyser
	//var CREATE_AT = RefTime(time.Date(2024, 11, 20, 6, 0, 0, 0, time.Local))
	//if neutronItem.Flow > 0.00001 && neutronItem.Mat1 < 0.00001 {
	//	return data, errors.New("Invalid flux and cao values")
	//}
	data = YAnalyser{
		Token:       RefString(fmt.Sprintf("neu-%d", neutronIndex)),
		AnalyserNum: RefString(productName),
		CreatedAt:   RefTime(dateTime),
		TestAt:      RefInt(int(dateTime.Unix())),
		Speed:       RefFloat64(neutronItem.BeltSpeed),
		Flux:        RefFloat64(neutronItem.Flow),
		Load:        RefFloat64(neutronItem.Load),
		Mat1:        RefFloat64(neutronItem.Mat1),
		Mat2:        RefFloat64(neutronItem.Mat2),
		Mat3:        RefFloat64(neutronItem.Mat3),
		Mat4:        RefFloat64(neutronItem.Mat4),
		Mat5:        RefFloat64(neutronItem.Mat5),
		Mat6:        RefFloat64(neutronItem.Mat6),
		Mat7:        RefFloat64(neutronItem.Mat7),
		Mat8:        RefFloat64(neutronItem.Mat8),
		Mat9:        RefFloat64(neutronItem.Mat9),
		Mat10:       RefFloat64(neutronItem.Mat10),
		Mat11:       RefFloat64(neutronItem.Mat11),
		Mat12:       RefFloat64(neutronItem.Mat12),
		Mat13:       RefFloat64(neutronItem.Mat13),
		Mat14:       RefFloat64(neutronItem.Mat14),
		Mat15:       RefFloat64(neutronItem.Mat15),
		Mat16:       RefFloat64(neutronItem.Mat16),
		Mat17:       RefFloat64(neutronItem.Mat17),
		Mat18:       RefFloat64(neutronItem.Mat18),
		Mat19:       RefFloat64(neutronItem.Mat19),
		Mat20:       RefFloat64(neutronItem.Mat20),
	}
	return data, nil
}

func NewBridgeDataWithRealTime(productName string, dataTime time.Time, bridgeIndex int, netWeight float64, vehicleNo string) WeighLogger {
	data := WeighLogger{
		Token:     RefString(fmt.Sprintf("bri-%d", bridgeIndex)),
		SiteCode:  RefString(productName),
		CheckTime: RefString(dataTime.Format("2006-01-02 15:04:05")),
		VehicleNo: RefString(vehicleNo),
		NetWeight: RefFloat64(netWeight),
	}
	return data
}
