package weighLogger
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	weighLoggerModel "sim_data_gen/models/weighLogger"
	. "sim_data_gen/utils"
)

// 列出 过磅记录
func ListObj(tx ... *gorm.DB)([]weighLoggerModel.WeighLogger,error){
	list := []weighLoggerModel.WeighLogger{}
	db := GetDb(tx ...).Table("weigh_logger").Find(&list)
	return list,db.Error
}

/**
* weighLogger数据库操作类 
*/
func CreateObj(weighLogger *weighLoggerModel.WeighLogger,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(weighLogger)
	return db.RowsAffected,db.Error
}

//  更新过磅记录
func UpdateObj(weighLogger *weighLoggerModel.WeighLogger,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(weighLogger)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(weighLogger *weighLoggerModel.WeighLogger,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(weighLogger)
	return db.RowsAffected,db.Error
}

//  更新过磅记录
func SaveObj(weighLogger *weighLoggerModel.WeighLogger,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(weighLogger)
	return db.RowsAffected,db.Error
}


// 根据ID获取过磅记录
func GetObjById(token *string,tx ... *gorm.DB)(weighLoggerModel.WeighLogger,error){
	weighLogger := weighLoggerModel.WeighLogger{Token:token}
	result := weighLoggerModel.WeighLogger{}
	db := GetDb(tx ...).Where(&weighLogger).Take(&result)
	return result,db.Error
}
//  用分页方式列出 过磅记录
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]weighLoggerModel.WeighLogger,error){
	sql := `select weigh_logger.busi_no as busi_no,weigh_logger.check_time as check_time,weigh_logger.direction as direction,weigh_logger.gross_weight as gross_weight,weigh_logger.net_weight as net_weight,weigh_logger.nt as nt,weigh_logger.site_code as site_code,weigh_logger.site_name as site_name,weigh_logger.tare_weight as tare_weight,weigh_logger.token as token,weigh_logger.update_at as update_at,weigh_logger.vehicle_no as vehicle_no from weigh_logger `
	list := []weighLoggerModel.WeighLogger{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from weigh_logger `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 过磅记录
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&weighLoggerModel.WeighLogger{Token:token})
	return db.RowsAffected,db.Error
}
