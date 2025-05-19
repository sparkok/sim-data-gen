package lorryNearbyTargetSpan
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	lorryNearbyTargetSpanModel "sim_data_gen/models/lorryNearbyTargetSpan"
	. "sim_data_gen/utils"
)

// 列出 卡车附近的目标
func ListObj(tx ... *gorm.DB)([]lorryNearbyTargetSpanModel.LorryNearbyTargetSpan,error){
	list := []lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{}
	db := GetDb(tx ...).Table("lorry_nearby_target_span").Find(&list)
	return list,db.Error
}

/**
* lorryNearbyTargetSpan数据库操作类 
*/
func CreateObj(lorryNearbyTargetSpan *lorryNearbyTargetSpanModel.LorryNearbyTargetSpan,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(lorryNearbyTargetSpan)
	return db.RowsAffected,db.Error
}

//  更新卡车附近的目标
func UpdateObj(lorryNearbyTargetSpan *lorryNearbyTargetSpanModel.LorryNearbyTargetSpan,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(lorryNearbyTargetSpan)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(lorryNearbyTargetSpan *lorryNearbyTargetSpanModel.LorryNearbyTargetSpan,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(lorryNearbyTargetSpan)
	return db.RowsAffected,db.Error
}

//  更新卡车附近的目标
func SaveObj(lorryNearbyTargetSpan *lorryNearbyTargetSpanModel.LorryNearbyTargetSpan,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(lorryNearbyTargetSpan)
	return db.RowsAffected,db.Error
}


// 根据ID获取卡车附近的目标
func GetObjById(token *string,tx ... *gorm.DB)(lorryNearbyTargetSpanModel.LorryNearbyTargetSpan,error){
	lorryNearbyTargetSpan := lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{Token:token}
	result := lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{}
	db := GetDb(tx ...).Where(&lorryNearbyTargetSpan).Take(&result)
	return result,db.Error
}
//  用分页方式列出 卡车附近的目标
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]lorryNearbyTargetSpanModel.LorryNearbyTargetSpan,error){
	sql := `select lorry_nearby_target_span.begin_utc as begin_utc,lorry_nearby_target_span.date_flag as date_flag,lorry_nearby_target_span.end_utc as end_utc,lorry_nearby_target_span.lorry_id as lorry_id,lorry_nearby_target_span.name as name,lorry_nearby_target_span.nearby_obj as nearby_obj,lorry_nearby_target_span.obj_type as obj_type,lorry_nearby_target_span.product_name as product_name,lorry_nearby_target_span.token as token from lorry_nearby_target_span `
	list := []lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from lorry_nearby_target_span `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 卡车附近的目标
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{Token:token})
	return db.RowsAffected,db.Error
}
