package mineProduct
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	mineProductModel "sim_data_gen/models/mineProduct"
	. "sim_data_gen/utils"
)

// 列出 产品
func ListObj(tx ... *gorm.DB)([]mineProductModel.MineProduct,error){
	list := []mineProductModel.MineProduct{}
	db := GetDb(tx ...).Table("mine_product").Find(&list)
	return list,db.Error
}

/**
* mineProduct数据库操作类 
*/
func CreateObj(mineProduct *mineProductModel.MineProduct,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(mineProduct)
	return db.RowsAffected,db.Error
}

//  更新产品
func UpdateObj(mineProduct *mineProductModel.MineProduct,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(mineProduct)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(mineProduct *mineProductModel.MineProduct,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(mineProduct)
	return db.RowsAffected,db.Error
}

//  更新产品
func SaveObj(mineProduct *mineProductModel.MineProduct,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(mineProduct)
	return db.RowsAffected,db.Error
}


// 根据ID获取产品
func GetObjById(token *string,tx ... *gorm.DB)(mineProductModel.MineProduct,error){
	mineProduct := mineProductModel.MineProduct{Token:token}
	result := mineProductModel.MineProduct{}
	db := GetDb(tx ...).Where(&mineProduct).Take(&result)
	return result,db.Error
}
//  用分页方式列出 产品
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]mineProductModel.MineProduct,error){
	sql := `select mine_product.content_limits as content_limits,mine_product.mat_indexes as mat_indexes,mine_product.name as name,mine_product.status as status,mine_product.token as token from mine_product `
	list := []mineProductModel.MineProduct{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from mine_product `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 产品
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&mineProductModel.MineProduct{Token:token})
	return db.RowsAffected,db.Error
}
