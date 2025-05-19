package location

import (
	locationModel "sim_data_gen/models/location"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

func GetValidLocationsByIds(ids []string, minUtc int64, tx ...*gorm.DB) ([]locationModel.Location, error) {
	list := []locationModel.Location{}
	db := GetDb(tx...).Table("location").Where("token in (?) and Utc > ?", ids, minUtc).Find(&list)
	return list, db.Error
}
