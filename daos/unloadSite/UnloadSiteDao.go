package unloadSite

import (
	"sim_data_gen/models/common"
	unloadSiteModel "sim_data_gen/models/unloadSite"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 列出 卸货点
func ListObj(tx ...*gorm.DB) ([]unloadSiteModel.UnloadSite, error) {
	list := []unloadSiteModel.UnloadSite{}
	db := GetDb(tx...).Table("unload_site").Find(&list)
	return list, db.Error
}

/**
* unloadSite数据库操作类
 */
func CreateObj(unloadSite *unloadSiteModel.UnloadSite, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(unloadSite)
	return db.RowsAffected, db.Error
}

//  更新卸货点
func UpdateObj(unloadSite *unloadSiteModel.UnloadSite, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Updates(unloadSite)
	return db.RowsAffected, db.Error
}

//  更新非null字段
func UpdateObjItem(unloadSite *unloadSiteModel.UnloadSite, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).UpdateColumns(unloadSite)
	return db.RowsAffected, db.Error
}

//  更新卸货点
func SaveObj(unloadSite *unloadSiteModel.UnloadSite, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Save(unloadSite)
	return db.RowsAffected, db.Error
}

// 根据ID获取卸货点
func GetObjById(token *string, tx ...*gorm.DB) (unloadSiteModel.UnloadSite, error) {
	unloadSite := unloadSiteModel.UnloadSite{Token: token}
	result := unloadSiteModel.UnloadSite{}
	db := GetDb(tx...).Where(&unloadSite).Take(&result)
	return result, db.Error
}

//  用分页方式列出 卸货点
func PageObj(conditions string, order string, curPage int, pageSize int, tx ...*gorm.DB) ([]unloadSiteModel.UnloadSite, error) {
	sql := `select unload_site.capacity as capacity,unload_site.geom as geom,unload_site.name as name,unload_site.nt as nt,unload_site.token as token,unload_site.x as x,unload_site.y as y from unload_site `
	list := []unloadSiteModel.UnloadSite{}
	limitAndOffset := MakeLimitOffset(curPage, pageSize)
	db := GetDb(tx...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list, db.Error
}

func Count4Page(conditions string, tx ...*gorm.DB) (int64, error) {
	sql := `select count(*) as Count from unload_site `
	count := common.Count{}
	db := GetDb(tx...).Raw(sql + conditions).Take(&count)
	return count.Count, db.Error
}

// 根据id删除 卸货点
func DeleteObj(token *string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&unloadSiteModel.UnloadSite{Token: token})
	return db.RowsAffected, db.Error
}
