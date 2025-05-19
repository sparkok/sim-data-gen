package lorry

import (
	lorryModel "sim_data_gen/models/lorry"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 根据ID获取卡车
func GetObjByName(name *string, tx ...*gorm.DB) (lorryModel.Lorry, error) {
	lorry := lorryModel.Lorry{Name: name}
	result := lorryModel.Lorry{}
	db := GetDb(tx...).Where(&lorry).Take(&result)
	return result, db.Error
}
