package lorryDiggerBindingLog

import (
	"sim_data_gen/models/common"
	lorryDiggerBindingLogModel "sim_data_gen/models/lorryDiggerBindingLog"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 列出 静态调度日志
func ListObj(tx ...*gorm.DB) ([]lorryDiggerBindingLogModel.LorryDiggerBindingLog, error) {
	list := []lorryDiggerBindingLogModel.LorryDiggerBindingLog{}
	db := GetDb(tx...).Table("lorry_digger_binding_log").Find(&list)
	return list, db.Error
}

/**
* lorryDiggerBindingLog数据库操作类
 */
func CreateObj(lorryDiggerBindingLog *lorryDiggerBindingLogModel.LorryDiggerBindingLog, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(lorryDiggerBindingLog)
	return db.RowsAffected, db.Error
}

// 更新静态调度日志
func UpdateObj(lorryDiggerBindingLog *lorryDiggerBindingLogModel.LorryDiggerBindingLog, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Updates(lorryDiggerBindingLog)
	return db.RowsAffected, db.Error
}

// 更新非null字段
func UpdateObjItem(lorryDiggerBindingLog *lorryDiggerBindingLogModel.LorryDiggerBindingLog, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).UpdateColumns(lorryDiggerBindingLog)
	return db.RowsAffected, db.Error
}

// 更新静态调度日志
func SaveObj(lorryDiggerBindingLog *lorryDiggerBindingLogModel.LorryDiggerBindingLog, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Save(lorryDiggerBindingLog)
	return db.RowsAffected, db.Error
}

// 根据ID获取静态调度日志
func GetObjByIdFully(token *string, tx ...*gorm.DB) (lorryDiggerBindingLogModel.LorryDiggerBindingLogFully, error) {
	sql := `select lorry_digger_binding_log.boom_group_name as boom_group_name,lorry_digger_binding_log.date_flag as date_flag,lorry_digger_binding_log.digger_name as digger_name,lorry_digger_binding_log.end_utc as end_utc,lorry_digger_binding_log.lorry_name as lorry_name,lorry_digger_binding_log.start_utc as start_utc,lorry_digger_binding_log.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp,digger0.token as token_of_digger,digger0.name as digger_desp,lorry0.token as token_of_lorry,lorry0.name as lorry_desp from lorry_digger_binding_log 
 		left join boom_group boom_group0 on (boom_group0.token = lorry_digger_binding_log.boom_group_id)
 		left join digger digger0 on (digger0.token = lorry_digger_binding_log.digger_id)
 		left join lorry lorry0 on (lorry0.token = lorry_digger_binding_log.lorry_id) where  lorry_digger_binding_log.token = ?`
	lorryDiggerBindingLog := lorryDiggerBindingLogModel.LorryDiggerBindingLogFully{}
	db := GetDb(tx...).Raw(sql, token).Take(&lorryDiggerBindingLog)
	return lorryDiggerBindingLog, db.Error
}

// 根据ID获取静态调度日志
func GetObjById(token *string, tx ...*gorm.DB) (lorryDiggerBindingLogModel.LorryDiggerBindingLog, error) {
	lorryDiggerBindingLog := lorryDiggerBindingLogModel.LorryDiggerBindingLog{Token: token}
	result := lorryDiggerBindingLogModel.LorryDiggerBindingLog{}
	db := GetDb(tx...).Where(&lorryDiggerBindingLog).Take(&result)
	return result, db.Error
}

// 用分页方式列出 静态调度日志
func PageObj(conditions string, order string, curPage int, pageSize int, tx ...*gorm.DB) ([]lorryDiggerBindingLogModel.LorryDiggerBindingLogFully, error) {
	sql := `select lorry_digger_binding_log.boom_group_name as boom_group_name,lorry_digger_binding_log.date_flag as date_flag,lorry_digger_binding_log.digger_name as digger_name,lorry_digger_binding_log.end_utc as end_utc,lorry_digger_binding_log.lorry_name as lorry_name,lorry_digger_binding_log.start_utc as start_utc,lorry_digger_binding_log.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp,digger0.token as token_of_digger,digger0.name as digger_desp,lorry0.token as token_of_lorry,lorry0.name as lorry_desp from lorry_digger_binding_log 
 		left join boom_group boom_group0 on (boom_group0.token = lorry_digger_binding_log.boom_group_id)
 		left join digger digger0 on (digger0.token = lorry_digger_binding_log.digger_id)
 		left join lorry lorry0 on (lorry0.token = lorry_digger_binding_log.lorry_id)`
	list := []lorryDiggerBindingLogModel.LorryDiggerBindingLogFully{}
	limitAndOffset := MakeLimitOffset(curPage, pageSize)
	db := GetDb(tx...).Raw(sql + conditions + order + " " + limitAndOffset).Find(&list)
	return list, db.Error
}

func Count4Page(conditions string, tx ...*gorm.DB) (int64, error) {
	sql := `select count(*) as Count from lorry_digger_binding_log 
 		left join boom_group boom_group0 on (boom_group0.token = lorry_digger_binding_log.boom_group_id)
 		left join digger digger0 on (digger0.token = lorry_digger_binding_log.digger_id)
 		left join lorry lorry0 on (lorry0.token = lorry_digger_binding_log.lorry_id)`
	count := common.Count{}
	db := GetDb(tx...).Raw(sql + conditions).Take(&count)
	return count.Count, db.Error
}

// 根据id删除 静态调度日志
func DeleteObj(token *string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&lorryDiggerBindingLogModel.LorryDiggerBindingLog{Token: token})
	return db.RowsAffected, db.Error
}

// 根据卡车id获取静态调度日志
func GetCurrentObjByLorryId(tokenOfLorry *string, dateFlag *string, tx ...*gorm.DB) (lorryDiggerBindingLogModel.LorryDiggerBindingLog, error) {
	lorryDiggerBindingLog := lorryDiggerBindingLogModel.LorryDiggerBindingLog{TokenOfLorry: tokenOfLorry, DateFlag: dateFlag}
	result := lorryDiggerBindingLogModel.LorryDiggerBindingLog{}
	db := GetDb(tx...).Where(&lorryDiggerBindingLog).Order("start_utc desc").Find(&result)
	return result, db.Error
}
