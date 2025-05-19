package locationGnssData

import (
	"gorm.io/gorm"
	locationGnssDataModel "target_schedule/models/locationGnssData"
	. "target_schedule/utils"
)

func CreateObj(locationGnssData *locationGnssDataModel.LocationGnssData, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(locationGnssData)
	return db.RowsAffected, db.Error
}
