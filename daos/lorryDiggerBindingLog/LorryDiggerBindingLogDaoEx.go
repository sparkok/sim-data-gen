package lorryDiggerBindingLog

import (
	lorryDiggerBindingLogModel "sim_data_gen/models/lorryDiggerBindingLog"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 根据ID获取静态调度日志
func FindObjBySpan(tokenOfLorry *string, utc int, tx ...*gorm.DB) (lorryDiggerBindingLogModel.LorryDiggerBindingLog, error) {
	obj := lorryDiggerBindingLogModel.LorryDiggerBindingLog{}
	db := GetDb(tx...).Table("lorry_digger_binding_log").Where("lorry_id = ? and ? between start_utc and end_utc", tokenOfLorry, utc).Take(&obj)
	return obj, db.Error
}
func FindObjBySpan4Name(dateFlag string, tokenOfName *string, utc int, tx ...*gorm.DB) (lorryDiggerBindingLogModel.LorryDiggerBindingLog, error) {
	obj := lorryDiggerBindingLogModel.LorryDiggerBindingLog{}
	//如果时间落在这个范围内,则说明有切换记录,这个是挖机和卡车对应关系的历史表
	db := GetDb(tx...).Table("lorry_digger_binding_log").Where("( date_flag = ? and lorry_name = ? ) and ( ? between start_utc and end_utc )", dateFlag, tokenOfName, utc).Take(&obj)
	//如果今天从来没切换过,则要去绑定表里面查找最新的记录
	return obj, db.Error
}

func ClearDataOfDay(dateFlag string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&lorryDiggerBindingLogModel.LorryDiggerBindingLog{}, "date_flag = ?", dateFlag)
	return db.RowsAffected, db.Error
}
