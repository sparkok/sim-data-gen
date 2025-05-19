package services

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"os"
	weighLoggerDao "sim_data_gen/daos/weighLogger"
	boomGroupModel "sim_data_gen/models/boomGroup"
	"sim_data_gen/models/weighLogger"
	. "sim_data_gen/utils"
	"time"
)

/*
*
  - 自动调度:

1.根据挖机的位置,自动确定它正在开采那个配矿单元。
2.根据卡车的位置,自动确定它正在和哪个挖机进行搭配。
3.根据目前的开采情况,自动确定使用了那个采矿方案。
*/
type WeightBridgeBoomGroupScanner struct {
	mineProductEnableAsActive int64 // 1 - 使能状态重定义了代表的其实是激活的状态, 0 - 未重定义
	lorryId2Name              map[string]string
	diggerIds                 []string
	lorryIdsInSystem          []string
	lorryNamesInSystem        []string
	boomGroupIds              []string
	boomGroupNames            []string
	boomGroupId2Info          map[string]boomGroupModel.BoomGroup
	vehicleFrom               *SystemBaseInfo
}

func (t *WeightBridgeBoomGroupScanner) Scan(productName string, dateFlag string, excelPath string) {
	t.vehicleFrom = NewMineInfo("mineFrom", productName, t.mineProductEnableAsActive)
	Logger.Info("WeightBridgeBoomGroupScanner.Execute")
	//allSql := strings.Builder{}
	//allSql.WriteString(t.makeClearSql(dateFlag))
	//allSql.WriteString(";\n")

	//打开excel文件
	var excelFile *excelize.File
	var err error
	if _, err = os.Stat(excelPath); !os.IsNotExist(err) {
		//存在则先删除
		os.Remove(excelPath)
	}
	if _, err = os.Stat(excelPath); os.IsNotExist(err) {
		excelFile = excelize.NewFile()
	}
	defer excelFile.Close()

	weightBridges, err := t.queryWeightBridges(productName, dateFlag)
	if err != nil {
		Logger.Error("queryWeightBridges", zap.Error(err))
		return
	}

	//create a shell
	// 检查是否存在名为 "Sheet1" 的工作表
	sheetWeightLogger := "Sheet1"
	t.makeSureSheet(excelFile, sheetWeightLogger)
	lineInWeighLogger := 1
	excelFile.SetColWidth(sheetWeightLogger, "A", "A", 40)
	excelFile.SetColWidth(sheetWeightLogger, "B", "B", 18)
	excelFile.SetColWidth(sheetWeightLogger, "C", "C", 40)
	excelFile.SetColWidth(sheetWeightLogger, "D", "D", 40)
	excelFile.SetColWidth(sheetWeightLogger, "E", "E", 20)
	excelFile.SetColWidth(sheetWeightLogger, "F", "F", 20)
	excelFile.SetColWidth(sheetWeightLogger, "G", "G", 25)
	excelFile.SetColWidth(sheetWeightLogger, "H", "H", 40)
	excelFile.SetColWidth(sheetWeightLogger, "I", "I", 40)

	t.WriteRow(excelFile, sheetWeightLogger, lineInWeighLogger, []interface{}{
		"Id",
		"时间",
		"卡车Id",
		"卡车",
		"方向",
		"净重",
		"记录时间",
		"挖机Id",
		"挖机编号",
		"配矿单元Id",
		"配矿单元",
	})
	lineInWeighLogger += 1

	//将所以目前的绑定设置为终止
	err = excelFile.SetSheetName("Sheet1", "地磅数据")
	if err != nil {
		Logger.Error("重命名失败:", zap.Error(err))
		return
	}
	t.vehicleFrom.loadInfo(dateFlag)
	diggerName2Id := t.vehicleFrom.CreateDiggerName2Id()
	diggerId2Name := t.vehicleFrom.CreateDiggerId2Name()
	var lorryId string
	var diggerId, diggerName string
	var boomGroupId, boomGroupName string
	for _, weightBridge := range weightBridges {
		lorryId, err = t.GetLorryId(*weightBridge.VehicleNo, diggerName2Id)
		if err != nil {
			Logger.Error("GetLorryId", zap.Error(err))
			return
		}
		diggerId, boomGroupId, err = t.CalcMineInfo(lorryId, weightBridge.CheckTime)
		diggerName = diggerId2Name[diggerId]
		boomGroupName = *t.boomGroupId2Info[boomGroupId].Name
		if err != nil {
			Logger.Error("CalcDiggerInfo", zap.Error(err))
			return
		}

		//获取地磅数据
		t.WriteRow(excelFile, sheetWeightLogger, lineInWeighLogger, []interface{}{
			"Id",
			"时间",
			lorryId,
			weightBridge.VehicleNo,
			*weightBridge.Direction,
			*weightBridge.NetWeight,
			*weightBridge.CheckTime,
			diggerId,
			diggerName,
			boomGroupId,
			boomGroupName,
		})
		lineInWeighLogger += 1

	}

	// 设置活动工作表
	excelFile.SetActiveSheet(1)
	// 保存文件
	if err := excelFile.SaveAs(excelPath); err != nil {
		Logger.Error("failed to save excel", zap.Error(err))
	}
	//os.WriteFile(strings.Replace(excelPath, ".xlsx", ".sql", 1), []byte(allSql.String()), os.FileMode.Perm(0777))
}
func (t *WeightBridgeBoomGroupScanner) makeSureSheet(f *excelize.File, sheetName string) (error, int) {
	index, err := f.GetSheetIndex(sheetName)
	if err != nil {
		Logger.Error("GetSheetIndex", zap.Error(err))
		return err, -1
	}
	if index == -1 {
		// 如果不存在，则创建工作表
		index, err = f.NewSheet(sheetName)
	}
	return err, index
}

func (t *WeightBridgeBoomGroupScanner) WriteRow(f *excelize.File, sheetName string, row int, values []interface{}) {
	// 定义要写入的数据
	// 写入一行数据
	cell := fmt.Sprintf("A%d", row)
	err := f.SetSheetRow(sheetName, cell, &values)
	if err != nil {
		Logger.Error("failed to set sheet row", zap.Error(err))
		return
	}
}

func (t *WeightBridgeBoomGroupScanner) queryWeightBridges(siteCode, dateStr string) ([]weighLogger.WeighLogger, error) {
	startFrom, endTo, err := t.GetDateAsStrBeginAndEnd(dateStr)
	if err != nil {
		return nil, err
	}
	return weighLoggerDao.ListObjDuring(siteCode, startFrom, endTo)
}

func NewWeightBridgeBoomGroupScanner(varNamePrx string) *WeightBridgeBoomGroupScanner {
	lorrySwitchScanner := new(WeightBridgeBoomGroupScanner)
	lorrySwitchScanner.mineProductEnableAsActive = GetConfig().Int64("productions.mineProductEnableAsActive", 1)
	return lorrySwitchScanner
}

func (t *WeightBridgeBoomGroupScanner) GetDateAsStrBeginAndEnd(dateAsStr string) (time.Time, time.Time, error) {
	//获取指定日期的起始和结束时间戳
	var beginTime, endTime time.Time
	var err error
	beginTime, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", dateAsStr), SystemZone)
	if err != nil {
		return beginTime, endTime, err
	}
	endTime, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 23:59:59", dateAsStr), SystemZone)
	if err != nil {
		return beginTime, endTime, err
	}
	return beginTime, endTime, nil
}

func (t *WeightBridgeBoomGroupScanner) GetLorryId(vehicleNo string, lorryName2Id map[string]string) (string, error) {
	if id, ok := lorryName2Id[vehicleNo]; ok {
		return id, nil
	} else {
		return "", nil
	}
}

//	func (t *WeightBridgeBoomGroupScanner) makeClearSql(dateAsStr string) string {
//		beginFrom, endTo, err := GetDateAsStrBeginAndEndUtc(dateAsStr)
//		if err != nil {
//			Logger.Error("GetDateAsStrBeginAndEndUtc", zap.Error(err))
//			return ""
//		}
//		sql := GetDb().ToSQL(func(tx *gorm.DB) *gorm.DB {
//			condition := "check_time >= ? and check_time <= ?"
//			return tx.Where(condition, beginFrom, endTo).Delete("weigh_logger")
//		})
//		return sql
//	}
//
//	func (t *WeightBridgeBoomGroupScanner) makeCreateSql(bindingLog lorryDiggerBindingLog.LorryDiggerBindingLog) string {
//		sql := GetDb().ToSQL(func(tx *gorm.DB) *gorm.DB {
//			return tx.Create(bindingLog)
//		})
//		return sql
//	}
func (t *WeightBridgeBoomGroupScanner) CalcMineInfo(lorryId string, checkTime *string) (string, string, error) {
	var diggerId, boomGroupId string
	var found bool
	checkTimeObj, err := time.ParseInLocation("2006-01-02 15:04:05", *checkTime, SystemZone)
	if err != nil {
		Logger.Error("failed to parse bridge time", zap.Error(err))
		return "", "", err
	}
	diggerId, boomGroupId, found = t.vehicleFrom.calcBoomGroupDugByDigLogger(lorryId, checkTimeObj.Unix())
	if found != true {
		Logger.Error("CalcMineInfo return false")
		return "", "", errors.New("CalcMineInfo return nothing")
	}

	return diggerId, boomGroupId, nil
}
