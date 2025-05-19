package yAnalyser

import (
	yAnalyserModel "sim_data_gen/models/yAnalyser"
	. "sim_data_gen/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

/*
* 获取今天皮带上最早的有矿记录
 */
func GetEarliestTimeExistedMineOfDay(crushingPlant string, date string, tx ...*gorm.DB) (yAnalyserModel.YAnalyser, error) {
	sql := `select COALESCE(min("test_at"),-1)  from y_analyser `

	conditions := "where (test_at > " + fmt.Sprintf("%d", todayFive(date).Unix()) + ") and flux > 0 and crushing_plant = '" + crushingPlant + "'"
	var dataTime int
	db := GetDb(tx...).Raw(sql + conditions).Take(&dataTime)
	result := yAnalyserModel.YAnalyser{}
	if db.Error != nil {
		return result, db.Error
	}
	if dataTime == -1 {
		//还没有今天的数据
		return result, fmt.Errorf("nothing")
	}
	yAnalyser := yAnalyserModel.YAnalyser{TestAt: &dataTime, CrushingPlant: &crushingPlant}
	db = GetDb(tx...).Where(&yAnalyser).Take(&result)
	return result, db.Error
}

func GetLastTimeExistedMineOfDay(crushingPlant string, date string, maxTime *int, tx ...*gorm.DB) (yAnalyserModel.YAnalyser, error) {
	sql := `select COALESCE(max("test_at"),-1)  from y_analyser `
	conditions := "where (test_at > " + fmt.Sprintf("%d", todayFive(date).Unix()) + ") and flux > 0 and crushing_plant = '" + crushingPlant + "'"
	if maxTime != nil {
		conditions += fmt.Sprintf(" and test_at <= %d", *maxTime)
	}
	var dataTime int
	db := GetDb(tx...).Raw(sql + conditions).Take(&dataTime)
	result := yAnalyserModel.YAnalyser{}
	if db.Error != nil {
		return result, db.Error
	}
	if dataTime == -1 {
		return result, fmt.Errorf("nothing")
	}
	yAnalyser := yAnalyserModel.YAnalyser{TestAt: &dataTime, CrushingPlant: &crushingPlant}
	db = GetDb(tx...).Where(&yAnalyser).Take(&result)
	return result, db.Error
}

func ListObjSince(crushingPlant string, startFrom int64, tx ...*gorm.DB) ([]yAnalyserModel.YAnalyser, error) {
	list := []yAnalyserModel.YAnalyser{}
	db := GetDb(tx...).Table("y_analyser").Where("crushing_plant = ? and test_at > ? ", crushingPlant, startFrom).Find(&list)
	return list, db.Error
}
func GetDataListDuring(crushingPlant string, startFrom int, endTo int, leftIncluded bool, tx ...*gorm.DB) ([]yAnalyserModel.YAnalyser, error) {
	list := []yAnalyserModel.YAnalyser{}
	var db *gorm.DB
	if leftIncluded {
		db = GetDb(tx...).Table("y_analyser").Where("(crushing_plant = ? ) and (test_at >= ? and test_at <= ? ) ", crushingPlant, startFrom, endTo).Find(&list)
	} else {
		db = GetDb(tx...).Table("y_analyser").Where("(crushing_plant = ? ) and (test_at > ? and test_at <= ? ) ", crushingPlant, startFrom, endTo).Find(&list)

	}

	return list, db.Error
}

func todayFive(someDay string) time.Time {

	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 05:00:00", someDay), SystemZone)
	if err != nil {
		return theTime
	}
	return theTime
}
func todayBegin(someDay string) time.Time {
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", someDay), SystemZone)
	if err != nil {
		return theTime
	}
	return theTime
}
func todayEnd(someDay string) time.Time {
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 23:59:59", someDay), SystemZone)
	if err != nil {
		return theTime
	}
	return theTime
}

func ClearDataOfDay(crushingPlant string, date string, tx ...*gorm.DB) (int64, error) {
	from := todayBegin(date).Unix()
	to := todayEnd(date).Unix()
	// bug修复：改为模型实例
	db := GetDb(tx...).Where("( test_at >= ? and test_at <= ? and crushing_plant = ? )", from, to, crushingPlant).Delete(&yAnalyserModel.YAnalyser{})
	return db.RowsAffected, db.Error
}
