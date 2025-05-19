package locationGnssData

import (
	locationGnssDataModel "sim_data_gen/models/locationGnssData"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

func ListObjDuring(targetId string, beginUtc int64, endUtc int64, tx ...*gorm.DB) ([]locationGnssDataModel.LocationGnssData, error) {
	list := []locationGnssDataModel.LocationGnssData{}
	db := GetDb(tx...).Table("location_gnss_data").
		Where("location_gnss_data.tid = ? AND location_gnss_data.utc BETWEEN ? AND ?", targetId, beginUtc, endUtc).
		Order("location_gnss_data.utc").Find(&list)
	return list, db.Error
}
