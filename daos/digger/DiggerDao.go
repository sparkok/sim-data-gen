package digger
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	diggerModel "sim_data_gen/models/digger"
	. "sim_data_gen/utils"
)

// 列出 挖机
func ListObj(tx ... *gorm.DB)([]diggerModel.Digger,error){
	list := []diggerModel.Digger{}
	db := GetDb(tx ...).Table("digger").Find(&list)
	return list,db.Error
}

/**
* digger数据库操作类 
*/
func CreateObj(digger *diggerModel.Digger,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(digger)
	return db.RowsAffected,db.Error
}

//  更新挖机
func UpdateObj(digger *diggerModel.Digger,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(digger)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(digger *diggerModel.Digger,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(digger)
	return db.RowsAffected,db.Error
}

//  更新挖机
func SaveObj(digger *diggerModel.Digger,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(digger)
	return db.RowsAffected,db.Error
}


// 根据ID获取挖机
func GetObjById(token *string,tx ... *gorm.DB)(diggerModel.Digger,error){
	digger := diggerModel.Digger{Token:token}
	result := diggerModel.Digger{}
	db := GetDb(tx ...).Where(&digger).Take(&result)
	return result,db.Error
}
//  用分页方式列出 挖机
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]diggerModel.Digger,error){
	sql := `select digger.attribs as attribs,digger.name as name,digger.produce as produce,digger.speed as speed,digger.status as status,digger.token as token,digger.utc as utc,digger.x as x,digger.y as y from digger `
	list := []diggerModel.Digger{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from digger `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 挖机
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&diggerModel.Digger{Token:token})
	return db.RowsAffected,db.Error
}
