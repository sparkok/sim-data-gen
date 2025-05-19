package contentPercent
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	contentPercentModel "sim_data_gen/models/contentPercent"
	. "sim_data_gen/utils"
)

// 列出 品位
func ListObj(tx ... *gorm.DB)([]contentPercentModel.ContentPercent,error){
	list := []contentPercentModel.ContentPercent{}
	db := GetDb(tx ...).Table("content_percent").Find(&list)
	return list,db.Error
}

/**
* contentPercent数据库操作类 
*/
func CreateObj(contentPercent *contentPercentModel.ContentPercent,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(contentPercent)
	return db.RowsAffected,db.Error
}

//  更新品位
func UpdateObj(contentPercent *contentPercentModel.ContentPercent,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(contentPercent)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(contentPercent *contentPercentModel.ContentPercent,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(contentPercent)
	return db.RowsAffected,db.Error
}

//  更新品位
func SaveObj(contentPercent *contentPercentModel.ContentPercent,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(contentPercent)
	return db.RowsAffected,db.Error
}


// 根据ID获取品位
func GetObjByIdFully(token *string,tx ... *gorm.DB)(contentPercentModel.ContentPercentFully,error){
	sql := `select content_percent.max_value as max_value,content_percent.min_value as min_value,content_percent.name as name,content_percent.num as num,content_percent.status as status,content_percent.token as token,composition0.token as token_of_composition,composition0.name as composition_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from content_percent 
 		left join composition composition0 on (composition0.token = content_percent.composition_id)
 		left join mine_product mine_product0 on (mine_product0.token = content_percent.mine_product_id) where  content_percent.token = ?`
	contentPercent := contentPercentModel.ContentPercentFully{}
	db := GetDb(tx ...).Raw(sql,token).Take(&contentPercent)
	return contentPercent,db.Error
}

// 根据ID获取品位
func GetObjById(token *string,tx ... *gorm.DB)(contentPercentModel.ContentPercent,error){
	contentPercent := contentPercentModel.ContentPercent{Token:token}
	result := contentPercentModel.ContentPercent{}
	db := GetDb(tx ...).Where(&contentPercent).Take(&result)
	return result,db.Error
}

//  用分页方式列出 品位
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]contentPercentModel.ContentPercentFully,error){
	sql := `select content_percent.max_value as max_value,content_percent.min_value as min_value,content_percent.name as name,content_percent.num as num,content_percent.status as status,content_percent.token as token,composition0.token as token_of_composition,composition0.name as composition_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from content_percent 
 		left join composition composition0 on (composition0.token = content_percent.composition_id)
 		left join mine_product mine_product0 on (mine_product0.token = content_percent.mine_product_id)`
	list := []contentPercentModel.ContentPercentFully{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from content_percent 
 		left join composition composition0 on (composition0.token = content_percent.composition_id)
 		left join mine_product mine_product0 on (mine_product0.token = content_percent.mine_product_id)`
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 品位
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&contentPercentModel.ContentPercent{Token:token})
	return db.RowsAffected,db.Error
}
