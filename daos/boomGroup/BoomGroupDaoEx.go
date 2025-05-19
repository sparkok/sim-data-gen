package boomGroup

import (
	boomGroupModel "sim_data_gen/models/boomGroup"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 根据ID获取采矿单元
func GetObjByName(name *string, tx ...*gorm.DB) (boomGroupModel.BoomGroup, error) {
	boomGroup := boomGroupModel.BoomGroup{Name: name}
	result := boomGroupModel.BoomGroup{}
	db := GetDb(tx...).Where(&boomGroup).Take(&result)
	return result, db.Error
}
