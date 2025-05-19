package diggerSwitchBoomGroupLog

import (
	diggerSwitchBoomGroupLogModel "sim_data_gen/models/diggerSwitchBoomGroupLog"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 列出 挖机配矿单元切换
func ListObjBeyond(diggerId, dateFlag string, beyondUtc int, tx ...*gorm.DB) ([]diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog, error) {
	list := []diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{}
	db := GetDb(tx...).Table("digger_switch_boom_group_log").Where("( apply_utc > ? and digger_id = ? and date_flag = ? and status = 2 )", beyondUtc, diggerId, dateFlag).Find(&list)
	return list, db.Error
}
