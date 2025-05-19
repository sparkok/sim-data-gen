package mineProduct

import (
	mineProductModel "sim_data_gen/models/mineProduct"
	. "sim_data_gen/utils"
	"fmt"
	"gorm.io/gorm"
)

func ListObjSortByName(tx ...*gorm.DB) ([]mineProductModel.MineProduct, error) {
	list := []mineProductModel.MineProduct{}
	db := GetDb(tx...).Table("mine_product").Order("name asc").Find(&list)
	return list, db.Error
}

// 1 - 卸货 ,2 - 装载
func FindBoomGroupInfoByVehicleNo(vehicleNo string, bridgeCheckTime string, tx ...*gorm.DB) string {
	//从调度日志中查
	//sql := `SELECT boom_group.token,boom_group.name,carry_logger.digger_id,carry_logger.logger_time FROM carry_logger
	//LEFT JOIN boom_group ON carry_logger.boom_group_id = boom_group.token
	//WHERE carry_logger.lorry_id = '%s' AND (carry_logger.logger_type = 2 OR carry_logger.logger_type = 1)
	//ORDER BY carry_logger.logger_time`
	//从称重日志中查
	sql := "SELECT boom_group_id,digger_id from load_goods_logger" +
		"LEFT JOIN weigh_logger ON load_goods_logger.weight_time = weigh_logger.check_time" +
		fmt.Sprintf("WHERE load_goods_logger.weight_time = '"+bridgeCheckTime+"' AND vehicle_no = '"+vehicleNo+"'")
	record := make(map[string]string)
	db := GetDb(tx...).Raw(sql).Take(&record)
	if db.Error != nil {
		Logger.Error(db.Error.Error())
		return ""
	} else {
		return record["boom_group_id"]
	}

}

// 根据name获取产品
func GetObjByName(name *string, tx ...*gorm.DB) (mineProductModel.MineProduct, error) {
	mineProduct := mineProductModel.MineProduct{Name: name}
	result := mineProductModel.MineProduct{}
	db := GetDb(tx...).Where(&mineProduct).Take(&result)
	return result, db.Error
}
