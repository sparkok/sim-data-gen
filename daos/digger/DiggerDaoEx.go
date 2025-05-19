package digger

import (
	diggerModel "sim_data_gen/models/digger"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

func ListValidDiggers(tx ...*gorm.DB) ([]diggerModel.Digger, error) {
	list := []diggerModel.Digger{}
	whereObj := diggerModel.Digger{Status: RefInt(1)}
	db := GetDb(tx...).Table("digger").Where(&whereObj).Find(&list)
	return list, db.Error
}
