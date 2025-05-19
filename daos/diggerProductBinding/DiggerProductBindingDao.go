package diggerProductBinding
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	diggerProductBindingModel "sim_data_gen/models/diggerProductBinding"
	. "sim_data_gen/utils"
)

// 列出 挖机产品绑定
func ListObj(tx ... *gorm.DB)([]diggerProductBindingModel.DiggerProductBinding,error){
	list := []diggerProductBindingModel.DiggerProductBinding{}
	db := GetDb(tx ...).Table("digger_product_binding").Find(&list)
	return list,db.Error
}

/**
* diggerProductBinding数据库操作类 
*/
func CreateObj(diggerProductBinding *diggerProductBindingModel.DiggerProductBinding,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(diggerProductBinding)
	return db.RowsAffected,db.Error
}

//  更新挖机产品绑定
func UpdateObj(diggerProductBinding *diggerProductBindingModel.DiggerProductBinding,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(diggerProductBinding)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(diggerProductBinding *diggerProductBindingModel.DiggerProductBinding,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(diggerProductBinding)
	return db.RowsAffected,db.Error
}

//  更新挖机产品绑定
func SaveObj(diggerProductBinding *diggerProductBindingModel.DiggerProductBinding,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(diggerProductBinding)
	return db.RowsAffected,db.Error
}


// 根据ID获取挖机产品绑定
func GetObjByIdFully(token *string,tx ... *gorm.DB)(diggerProductBindingModel.DiggerProductBindingFully,error){
	sql := `select digger_product_binding.name as name,digger_product_binding.token as token,digger0.token as token_of_digger,digger0.name as digger_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from digger_product_binding 
 		left join digger digger0 on (digger0.token = digger_product_binding.digger_id)
 		left join mine_product mine_product0 on (mine_product0.token = digger_product_binding.mine_product_id) where  digger_product_binding.token = ?`
	diggerProductBinding := diggerProductBindingModel.DiggerProductBindingFully{}
	db := GetDb(tx ...).Raw(sql,token).Take(&diggerProductBinding)
	return diggerProductBinding,db.Error
}

// 根据ID获取挖机产品绑定
func GetObjById(token *string,tx ... *gorm.DB)(diggerProductBindingModel.DiggerProductBinding,error){
	diggerProductBinding := diggerProductBindingModel.DiggerProductBinding{Token:token}
	result := diggerProductBindingModel.DiggerProductBinding{}
	db := GetDb(tx ...).Where(&diggerProductBinding).Take(&result)
	return result,db.Error
}

//  用分页方式列出 挖机产品绑定
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]diggerProductBindingModel.DiggerProductBindingFully,error){
	sql := `select digger_product_binding.name as name,digger_product_binding.token as token,digger0.token as token_of_digger,digger0.name as digger_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from digger_product_binding 
 		left join digger digger0 on (digger0.token = digger_product_binding.digger_id)
 		left join mine_product mine_product0 on (mine_product0.token = digger_product_binding.mine_product_id)`
	list := []diggerProductBindingModel.DiggerProductBindingFully{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from digger_product_binding 
 		left join digger digger0 on (digger0.token = digger_product_binding.digger_id)
 		left join mine_product mine_product0 on (mine_product0.token = digger_product_binding.mine_product_id)`
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 挖机产品绑定
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&diggerProductBindingModel.DiggerProductBinding{Token:token})
	return db.RowsAffected,db.Error
}
