package location

import (
	"sim_data_gen/models/common"
	locationModel "sim_data_gen/models/location"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 列出 位置
func ListObj(tx ...*gorm.DB) ([]locationModel.Location, error) {
	list := []locationModel.Location{}
	db := GetDb(tx...).Table("location").Find(&list)
	return list, db.Error
}

/**
* location数据库操作类
 */
func CreateObj(location *locationModel.Location, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(location)
	return db.RowsAffected, db.Error
}

//  更新位置
func UpdateObj(location *locationModel.Location, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Updates(location)
	return db.RowsAffected, db.Error
}

//  更新非null字段
func UpdateObjItem(location *locationModel.Location, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).UpdateColumns(location)
	return db.RowsAffected, db.Error
}

//  更新位置
func SaveObj(location *locationModel.Location, tx ...*gorm.DB) (int64, error) {
	//能保存为null的字段
	db := GetDb(tx...).Save(location)
	return db.RowsAffected, db.Error
}

// 根据ID获取位置
func GetObjById(token *string, tx ...*gorm.DB) (locationModel.Location, error) {
	location := locationModel.Location{Token: token}
	result := locationModel.Location{}
	db := GetDb(tx...).Where(&location).Take(&result)
	return result, db.Error
}

//  用分页方式列出 位置
func PageObj(conditions string, order string, curPage int, pageSize int, tx ...*gorm.DB) ([]locationModel.Location, error) {
	sql := `select location.alarm as alarm,location.heading as heading,location.name as name,location.speed as speed,location.status as status,location.token as token,location.utc as utc,location.x as x,location.y as y from location `
	list := []locationModel.Location{}
	limitAndOffset := MakeLimitOffset(curPage, pageSize)
	db := GetDb(tx...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list, db.Error
}

func Count4Page(conditions string, tx ...*gorm.DB) (int64, error) {
	sql := `select count(*) as Count from location `
	count := common.Count{}
	db := GetDb(tx...).Raw(sql + conditions).Take(&count)
	return count.Count, db.Error
}

// 根据id删除 位置
func DeleteObj(token *string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&locationModel.Location{Token: token})
	return db.RowsAffected, db.Error
}
