package diggerSwitchBoomGroupLog
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	diggerSwitchBoomGroupLogModel "sim_data_gen/models/diggerSwitchBoomGroupLog"
	. "sim_data_gen/utils"
)

// 列出 挖机配矿单元切换
func ListObj(tx ... *gorm.DB)([]diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog,error){
	list := []diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{}
	db := GetDb(tx ...).Table("digger_switch_boom_group_log").Find(&list)
	return list,db.Error
}

/**
* diggerSwitchBoomGroupLog数据库操作类 
*/
func CreateObj(diggerSwitchBoomGroupLog *diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(diggerSwitchBoomGroupLog)
	return db.RowsAffected,db.Error
}

//  更新挖机配矿单元切换
func UpdateObj(diggerSwitchBoomGroupLog *diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(diggerSwitchBoomGroupLog)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(diggerSwitchBoomGroupLog *diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(diggerSwitchBoomGroupLog)
	return db.RowsAffected,db.Error
}

//  更新挖机配矿单元切换
func SaveObj(diggerSwitchBoomGroupLog *diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(diggerSwitchBoomGroupLog)
	return db.RowsAffected,db.Error
}


// 根据ID获取挖机配矿单元切换
func GetObjById(token *string,tx ... *gorm.DB)(diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog,error){
	diggerSwitchBoomGroupLog := diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{Token:token}
	result := diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{}
	db := GetDb(tx ...).Where(&diggerSwitchBoomGroupLog).Take(&result)
	return result,db.Error
}
//  用分页方式列出 挖机配矿单元切换
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog,error){
	sql := `select digger_switch_boom_group_log.apply_utc as apply_utc,digger_switch_boom_group_log.boom_group_id as boom_group_id,digger_switch_boom_group_log.date_flag as date_flag,digger_switch_boom_group_log.digger_id as digger_id,digger_switch_boom_group_log.name as name,digger_switch_boom_group_log.status as status,digger_switch_boom_group_log.submit_utc as submit_utc,digger_switch_boom_group_log.token as token from digger_switch_boom_group_log `
	list := []diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from digger_switch_boom_group_log `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 挖机配矿单元切换
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{Token:token})
	return db.RowsAffected,db.Error
}
