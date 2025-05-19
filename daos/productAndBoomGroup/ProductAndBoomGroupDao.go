package productAndBoomGroup

import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	productAndBoomGroupModel "sim_data_gen/models/productAndBoomGroup"
	. "sim_data_gen/utils"
)

// 列出 产品和配矿单元
func ListObj(tx ...*gorm.DB) ([]productAndBoomGroupModel.ProductAndBoomGroup, error) {
	list := []productAndBoomGroupModel.ProductAndBoomGroup{}
	db := GetDb(tx...).Table("product_and_boom_group").Find(&list)
	return list, db.Error
}

/**
* productAndBoomGroup数据库操作类
 */
func CreateObj(productAndBoomGroup *productAndBoomGroupModel.ProductAndBoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(productAndBoomGroup)
	return db.RowsAffected, db.Error
}

// 更新产品和配矿单元
func UpdateObj(productAndBoomGroup *productAndBoomGroupModel.ProductAndBoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Updates(productAndBoomGroup)
	return db.RowsAffected, db.Error
}

// 更新非null字段
func UpdateObjItem(productAndBoomGroup *productAndBoomGroupModel.ProductAndBoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).UpdateColumns(productAndBoomGroup)
	return db.RowsAffected, db.Error
}

// 更新产品和配矿单元
func SaveObj(productAndBoomGroup *productAndBoomGroupModel.ProductAndBoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Save(productAndBoomGroup)
	return db.RowsAffected, db.Error
}

// 根据ID获取产品和配矿单元
func GetObjByIdFully(token *string, tx ...*gorm.DB) (productAndBoomGroupModel.ProductAndBoomGroupFully, error) {
	sql := `select product_and_boom_group.name as name,product_and_boom_group.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from product_and_boom_group 
 		left join boom_group boom_group0 on (boom_group0.token = product_and_boom_group.boom_group_id)
 		left join mine_product mine_product0 on (mine_product0.token = product_and_boom_group.mine_product_id) where  product_and_boom_group.token = ?`
	productAndBoomGroup := productAndBoomGroupModel.ProductAndBoomGroupFully{}
	db := GetDb(tx...).Raw(sql, token).Take(&productAndBoomGroup)
	return productAndBoomGroup, db.Error
}

// 根据ID获取产品和配矿单元
func GetObjById(token *string, tx ...*gorm.DB) (productAndBoomGroupModel.ProductAndBoomGroup, error) {
	productAndBoomGroup := productAndBoomGroupModel.ProductAndBoomGroup{Token: token}
	result := productAndBoomGroupModel.ProductAndBoomGroup{}
	db := GetDb(tx...).Where(&productAndBoomGroup).Take(&result)
	return result, db.Error
}

// 用分页方式列出 产品和配矿单元
func PageObj(conditions string, order string, curPage int, pageSize int, tx ...*gorm.DB) ([]productAndBoomGroupModel.ProductAndBoomGroupFully, error) {
	sql := `select product_and_boom_group.name as name,product_and_boom_group.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from product_and_boom_group 
 		left join boom_group boom_group0 on (boom_group0.token = product_and_boom_group.boom_group_id)
 		left join mine_product mine_product0 on (mine_product0.token = product_and_boom_group.mine_product_id)`
	list := []productAndBoomGroupModel.ProductAndBoomGroupFully{}
	limitAndOffset := MakeLimitOffset(curPage, pageSize)
	db := GetDb(tx...).Raw(sql + conditions + order + " " + limitAndOffset).Find(&list)
	return list, db.Error
}

func Count4Page(conditions string, tx ...*gorm.DB) (int64, error) {
	sql := `select count(*) as Count from product_and_boom_group 
 		left join boom_group boom_group0 on (boom_group0.token = product_and_boom_group.boom_group_id)
 		left join mine_product mine_product0 on (mine_product0.token = product_and_boom_group.mine_product_id)`
	count := common.Count{}
	db := GetDb(tx...).Raw(sql + conditions).Take(&count)
	return count.Count, db.Error
}

// 根据id删除 产品和配矿单元
func DeleteObj(token *string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&productAndBoomGroupModel.ProductAndBoomGroup{Token: token})
	return db.RowsAffected, db.Error
}
