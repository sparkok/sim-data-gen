package weighLogger

import (
	"sim_data_gen/entity"
	weighLoggerModel "sim_data_gen/models/weighLogger"
	. "sim_data_gen/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 列出从 startFrom以后的过磅记录
func ListObjSince(siteCode string, startFrom *time.Time, tx ...*gorm.DB) ([]weighLoggerModel.WeighLogger, error) {
	list := []weighLoggerModel.WeighLogger{}
	db := GetDb(tx...).Table("weigh_logger").Where("( check_time > ? ) and ( site_code = ? )", startFrom.Format("2006-01-02 15:04:05"), siteCode).Find(&list)
	return list, db.Error
}
func ListObjDuring(siteCode string, startFrom time.Time, endTo time.Time, tx ...*gorm.DB) ([]weighLoggerModel.WeighLogger, error) {
	list := []weighLoggerModel.WeighLogger{}
	db := GetDb(tx...).Table("weigh_logger").Where("( check_time >= ? and check_time <= ? ) and ( site_code = ? )", startFrom.Format("2006-01-02 15:04:05"), endTo.Format("2006-01-02 15:04:05"), siteCode).Find(&list)
	return list, db.Error
}
func GetDataListDuring(siteCode string, startFrom *time.Time, endTo *time.Time, leftIncluded bool, tx ...*gorm.DB) ([]weighLoggerModel.WeighLogger, error) {
	list := []weighLoggerModel.WeighLogger{}
	var db *gorm.DB
	if leftIncluded {
		db = GetDb(tx...).Table("weigh_logger").Where("( check_time >= ? and check_time <= ? ) and ( site_code = ? )", startFrom.Format("2006-01-02 15:04:05"), endTo.Format("2006-01-02 15:04:05"), siteCode).Find(&list)
	} else {
		db = GetDb(tx...).Table("weigh_logger").Where("( check_time > ? and check_time <= ? ) and ( site_code = ? )", startFrom.Format("2006-01-02 15:04:05"), endTo.Format("2006-01-02 15:04:05"), siteCode).Find(&list)
	}
	return list, db.Error
}

/**
 * 获取今天第一个过磅记录
 */
func GetEarliestTimeExistedMineOfDay(siteCode string, dateFlag string, tx ...*gorm.DB) (weighLoggerModel.WeighLogger, error) {
	//获取今天第一个地磅数据过磅时间
	sql := `select COALESCE(min(check_time),'-1')  from weigh_logger `
	conditions := "where ( direction = '1' and check_time > '" + fmt.Sprintf("%s 05:00:00", dateFlag) + "') and ( site_code = '" + siteCode + "')"
	var checkTime string
	db := GetDb(tx...).Raw(sql + conditions).Take(&checkTime)
	result := weighLoggerModel.WeighLogger{}
	if db.Error != nil {
		return result, db.Error
	}
	if checkTime == "-1" {
		return result, fmt.Errorf("nothing")
	}
	yAnalyser := weighLoggerModel.WeighLogger{CheckTime: &checkTime, SiteCode: &siteCode}
	db = GetDb(tx...).Where(&yAnalyser).Take(&result)
	return result, db.Error
}

/**
 * 获取最近的一条地磅数据
 */
func GetLatestTimeExistedMineOfDay(siteCode string, dateFlag string, maxTime *int, tx ...*gorm.DB) (weighLoggerModel.WeighLogger, error) {
	//获取今天第一个地磅数据过磅时间
	sql := `select COALESCE(max(check_time),'-1')  from weigh_logger `
	conditions := "where ( direction = '1' and check_time > '" + fmt.Sprintf("%s 05:00:00", dateFlag) + "') and ( site_code = '" + siteCode + "')"
	if maxTime != nil {
		conditions += fmt.Sprintf(" and check_time <= '%s'", entity.UnixTime2StrTime(maxTime))
	}

	var checkTime string
	db := GetDb(tx...).Raw(sql + conditions).Take(&checkTime)
	result := weighLoggerModel.WeighLogger{}
	if db.Error != nil {
		return result, db.Error
	}
	if checkTime == "-1" {
		return result, fmt.Errorf("nothing")
	}
	yAnalyser := weighLoggerModel.WeighLogger{CheckTime: &checkTime, SiteCode: &siteCode}
	db = GetDb(tx...).Where(&yAnalyser).Take(&result)
	return result, db.Error
}

func ClearDataOfDay(siteCode string, dateFlag string, tx ...*gorm.DB) (int64, error) {
	from := time2Str(someDayBegin(dateFlag))
	end := time2Str(someDayEnd(dateFlag))
	db := GetDb(tx...).Where("( check_time >= ? and check_time <= ? ) and ( site_code = ? )", from, end, siteCode).Delete(&weighLoggerModel.WeighLogger{})
	return db.RowsAffected, db.Error
}
func time2Str(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}

func someDayBegin(someDay string) time.Time {
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", someDay), SystemZone)
	if err != nil {
		return theTime
	}
	return theTime
}
func someDayEnd(someDay string) time.Time {
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 23:59:59", someDay), SystemZone)
	if err != nil {
		return theTime
	}
	return theTime
}
