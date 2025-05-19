package composition
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	compositionModel "sim_data_gen/models/composition"
	. "sim_data_gen/utils"
)

// 列出 化学成分
func ListObj(tx ... *gorm.DB)([]compositionModel.Composition,error){
	list := []compositionModel.Composition{}
	db := GetDb(tx ...).Table("composition").Find(&list)
	return list,db.Error
}

/**
* composition数据库操作类 
*/
func CreateObj(composition *compositionModel.Composition,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(composition)
	return db.RowsAffected,db.Error
}

//  更新化学成分
func UpdateObj(composition *compositionModel.Composition,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(composition)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(composition *compositionModel.Composition,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(composition)
	return db.RowsAffected,db.Error
}

//  更新化学成分
func SaveObj(composition *compositionModel.Composition,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(composition)
	return db.RowsAffected,db.Error
}


// 根据ID获取化学成分
func GetObjById(token *string,tx ... *gorm.DB)(compositionModel.Composition,error){
	composition := compositionModel.Composition{Token:token}
	result := compositionModel.Composition{}
	db := GetDb(tx ...).Where(&composition).Take(&result)
	return result,db.Error
}
//  用分页方式列出 化学成分
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]compositionModel.Composition,error){
	sql := `select composition.name as name,composition.token as token from composition `
	list := []compositionModel.Composition{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from composition `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 化学成分
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&compositionModel.Composition{Token:token})
	return db.RowsAffected,db.Error
}
