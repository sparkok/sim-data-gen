package test

import (
	. "sim_data_gen/calculate"
	. "sim_data_gen/models/weighLogger"
	. "sim_data_gen/models/yAnalyser"
	. "sim_data_gen/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

var PRODUCT_NAME = RefString("cp1")
var CREATE_AT = RefTime(time.Date(2024, 11, 20, 6, 0, 0, 0, time.Local))
var TEST_FROM = time.Date(2024, 11, 22, 6, 0, 0, 0, time.Local)

var neutronIndex = 0

func NewYAnalyser(seconds int, fluxPerMinute float64, cao float64) YAnalyser {
	data := YAnalyser{
		Token:       RefString(fmt.Sprintf("neu-%d", neutronIndex)),
		AnalyserNum: PRODUCT_NAME,
		CreatedAt:   CREATE_AT,
		TestAt:      RefInt(int(TEST_FROM.Add(time.Duration(seconds) * time.Second).Unix())),
		Flux:        RefFloat64(fluxPerMinute * 60),
		Load:        RefFloat64(30),
		Mat1:        RefFloat64(cao),
	}
	neutronIndex++
	return data
}
func NewYAnalyserWithRealTime(dateTime time.Time, flux float64, load float64, cao float64) (YAnalyser, error) {
	var data YAnalyser
	if flux > 0.00001 && cao < 0.00001 {
		return data, errors.New("Invalid flux and cao values")
	}
	data = YAnalyser{
		Token:       RefString(fmt.Sprintf("neu-%d", neutronIndex)),
		AnalyserNum: PRODUCT_NAME,
		CreatedAt:   CREATE_AT,
		TestAt:      RefInt(int(dateTime.Unix())),
		Flux:        RefFloat64(flux),
		Load:        RefFloat64(load),
		Mat1:        RefFloat64(cao),
	}
	neutronIndex++
	return data, nil
}
func Reset() {
	bridgeIndex = 0
	neutronIndex = 0
}

var bridgeIndex = 0

func NewBridgeData(seconds int, netWeight float64, vehicleNo string) WeighLogger {
	data := WeighLogger{
		Token:     RefString(fmt.Sprintf("bri-%d", bridgeIndex)),
		SiteCode:  PRODUCT_NAME,
		CheckTime: RefString(TEST_FROM.Add(time.Duration(seconds) * time.Second).Format("2006-01-02 15:04-05")),
		VehicleNo: RefString(vehicleNo),
		NetWeight: RefFloat64(netWeight),
	}
	bridgeIndex++
	return data
}
func NewBridgeDataWithRealTime(dataTime time.Time, netWeight float64, vehicleNo string) WeighLogger {
	data := WeighLogger{
		Token:     RefString(fmt.Sprintf("bri-%d", bridgeIndex)),
		SiteCode:  PRODUCT_NAME,
		CheckTime: RefString(dataTime.Format("2006-01-02 15:04-05")),
		VehicleNo: RefString(vehicleNo),
		NetWeight: RefFloat64(netWeight),
	}
	bridgeIndex++
	return data
}
func TestCondition2(t *testing.T) {
	//		//情况2:如果最后一条中子仪数据累计完,没有第一条地磅的数据的初始重量大
	//		// 中子仪初始总重量16: |中子仪数据,16 |中子仪数据,16 |中子仪数据,16|
	//		// 地磅初始总重量70: |地磅,70|
	//		// 这是一种错误,正常情况下不会出现这种情况
	//		// 处理后 =>
	//		// 处理策略:跳过中子仪数据
	//		// 中子仪初始总重量 64: |中子仪数据空
	//		// 地磅初始总重量 70 : |地磅,70|
	//
	//		//情况3:如果最后一条地磅数据累计完,没有第一条中子仪的数据完成后质量大
	//		// 初始总重量0: |中子仪数据,16 |中子仪数据,16 |中子仪数据,16|
	//		// 初始总重量0: |地磅,70|
	//		// 处理掉重量小于第一条地磅数据之前的中子仪数据
	//
	// }

}
func TestCondition1(t *testing.T) {
	//情况1:如果最后一条中子仪数据累计完,没有第一条的地磅数据完成后质量大
	// 中子仪初始总重量0: |A 中子仪数据,18|B 中子仪数据,18 |C 中子仪数据,18|
	// 地磅初始总重量0: |地磅,70|
	// 处理掉重量小于第一条地磅数据之前的中子仪数据
	// 处理后 =>
	// 处理策略: 认定地磅所对应车辆的配矿单元的品位为 A,B,C的品均值,并检测A,B,C应该品位是否相当。
	// 中子仪初始总重量 48 |中子仪数据空
	// 地磅初始总重量0: |地磅,40|
	// 模拟中子仪数据
	neutronData := []YAnalyser{
		NewYAnalyser(60, 18.0, 48.0),
		NewYAnalyser(60, 18.0, 48.0),
		NewYAnalyser(60, 18.0, 47.0),
	}

	neutronDataInitSumMass := 0.0

	// 模拟地磅数据
	bridgeData := []WeighLogger{
		NewBridgeData(70, 40, "闽G382731"),
		NewBridgeData(70, 80, "闽G323293"),
	}
	bridgeDataInitSumMass := 0.0

	// 创建 TotalMassAssignmentCalculator 实例
	mixList := NewTotalMassAssignmentCalculator(&neutronData, neutronDataInitSumMass, &bridgeData, bridgeDataInitSumMass)
	_, _, err := mixList.Predict()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("test ok")
	}
}

func TestFromDataDir(t *testing.T) {
	fileName := "241110"
	neutronReader := WeightDataReader{SourceType: SourceNeutronGauge}
	neutronData := []YAnalyser{}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	Reset()
	if err := neutronReader.ProcessFile("../../weight-data-statistic", fmt.Sprintf("%s.xlsx", fileName)); err == nil {
		// 处理错误，例如打印错误信息或返回
		dataLen := len(neutronReader.NeutronData)
		fmt.Printf("%d\n", dataLen)
		for _, neutronItem := range neutronReader.NeutronData {
			dataTime, err := time.ParseInLocation("01/02/06 15:04", neutronItem.Time, loc)
			if err != nil {
				Logger.Error("failed to parse dataTime", zap.String("dataTime", neutronItem.Time))
				continue
			}
			analyser, err := NewYAnalyserWithRealTime(dataTime, neutronItem.Flow, neutronItem.Load, neutronItem.Mat1)
			if err != nil {
				continue
			}
			neutronData = append(neutronData, analyser)
		}
	} else {
		Logger.Error("Failed to process file: %v", zap.Error(err))
		return
	}
	bridgeReader := WeightDataReader{SourceType: SourceWeighBridge}
	bridgeData := []WeighLogger{}
	if err := bridgeReader.ProcessFile("../../weight-data-statistic", fmt.Sprintf("%s.xlsx", fileName)); err == nil {
		dataLen := len(bridgeReader.BridgeData)
		fmt.Printf("%d\n", dataLen)
		//获取今天5点钟的时间
		for _, bridgeItem := range bridgeReader.BridgeData {
			dataTime, err := time.ParseInLocation("2006/01/02 15:04:05", bridgeItem.TareTime, SystemZone)
			if err != nil {
				Logger.Error("failed to parse dataTime", zap.String("dataTime", bridgeItem.TareTime))
				continue
			}
			weighLogger := NewBridgeDataWithRealTime(dataTime, bridgeItem.NetWeight, bridgeItem.VehicleNo)
			bridgeData = append(bridgeData, weighLogger)
		}
	} else {
		// 处理错误，例如打印错误信息或返回
		Logger.Error("Failed to process file: %v", zap.Error(err))
		return
	}
	neutronDataInitSumMass := 0.0
	bridgeDataInitSumMass := 0.0

	// 创建 TotalMassAssignmentCalculator 实例
	mixList := NewTotalMassAssignmentCalculator(&neutronData, neutronDataInitSumMass, &bridgeData, bridgeDataInitSumMass)
	mixList.Predict()
}
