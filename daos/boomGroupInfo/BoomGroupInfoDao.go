package boomGroupInfo

import (
	boomGroupInfoModel "sim_data_gen/models/boomGroupInfo"
	"sim_data_gen/models/common"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 列出 采矿扩展
func ListObj(tx ...*gorm.DB) ([]boomGroupInfoModel.BoomGroupInfo, error) {
	list := []boomGroupInfoModel.BoomGroupInfo{}
	db := GetDb(tx...).Table("boom_group_info").Find(&list)
	return list, db.Error
}

/**
* boomGroupInfo数据库操作类
 */
func CreateObj(boomGroupInfo *boomGroupInfoModel.BoomGroupInfo, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(boomGroupInfo)
	return db.RowsAffected, db.Error
}

// 更新采矿扩展
func UpdateObj(boomGroupInfo *boomGroupInfoModel.BoomGroupInfo, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Updates(boomGroupInfo)
	return db.RowsAffected, db.Error
}

// 更新非null字段
func UpdateObjItem(boomGroupInfo *boomGroupInfoModel.BoomGroupInfo, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).UpdateColumns(boomGroupInfo)
	return db.RowsAffected, db.Error
}

// 更新采矿扩展
func SaveObj(boomGroupInfo *boomGroupInfoModel.BoomGroupInfo, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Save(boomGroupInfo)
	return db.RowsAffected, db.Error
}

// 根据ID获取采矿扩展
func GetObjByIdFully(token *string, tx ...*gorm.DB) (boomGroupInfoModel.BoomGroupInfoFully, error) {
	sql := `select boom_group_info.created_at as created_at,boom_group_info.material1 as material1,boom_group_info.material10 as material10,boom_group_info.material11 as material11,boom_group_info.material12 as material12,boom_group_info.material13 as material13,boom_group_info.material14 as material14,boom_group_info.material15 as material15,boom_group_info.material16 as material16,boom_group_info.material17 as material17,boom_group_info.material18 as material18,boom_group_info.material19 as material19,boom_group_info.material2 as material2,boom_group_info.material20 as material20,boom_group_info.material3 as material3,boom_group_info.material4 as material4,boom_group_info.material5 as material5,boom_group_info.material6 as material6,boom_group_info.material7 as material7,boom_group_info.material8 as material8,boom_group_info.material9 as material9,boom_group_info.name as name,boom_group_info.probability1 as probability1,boom_group_info.probability10 as probability10,boom_group_info.probability11 as probability11,boom_group_info.probability12 as probability12,boom_group_info.probability13 as probability13,boom_group_info.probability14 as probability14,boom_group_info.probability15 as probability15,boom_group_info.probability16 as probability16,boom_group_info.probability17 as probability17,boom_group_info.probability18 as probability18,boom_group_info.probability19 as probability19,boom_group_info.probability2 as probability2,boom_group_info.probability20 as probability20,boom_group_info.probability3 as probability3,boom_group_info.probability4 as probability4,boom_group_info.probability5 as probability5,boom_group_info.probability6 as probability6,boom_group_info.probability7 as probability7,boom_group_info.probability8 as probability8,boom_group_info.probability9 as probability9,boom_group_info.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp from boom_group_info 
 		left join boom_group boom_group0 on (boom_group0.token = boom_group_info.boom_group_id) where boom_group_info.token = ?`
	boomGroupInfo := boomGroupInfoModel.BoomGroupInfoFully{}
	db := GetDb(tx...).Raw(sql, token).Take(&boomGroupInfo)
	return boomGroupInfo, db.Error
}

// 根据ID获取采矿扩展
func GetObjById(token *string, tx ...*gorm.DB) (boomGroupInfoModel.BoomGroupInfo, error) {
	boomGroupInfo := boomGroupInfoModel.BoomGroupInfo{Token: token}
	result := boomGroupInfoModel.BoomGroupInfo{}
	db := GetDb(tx...).Where(&boomGroupInfo).Take(&result)
	return result, db.Error
}

// 用分页方式列出 采矿扩展
func PageObj(conditions string, order string, curPage int, pageSize int, tx ...*gorm.DB) ([]boomGroupInfoModel.BoomGroupInfoFully, error) {
	sql := `select boom_group_info.created_at as created_at,boom_group_info.material1 as material1,boom_group_info.material10 as material10,boom_group_info.material11 as material11,boom_group_info.material12 as material12,boom_group_info.material13 as material13,boom_group_info.material14 as material14,boom_group_info.material15 as material15,boom_group_info.material16 as material16,boom_group_info.material17 as material17,boom_group_info.material18 as material18,boom_group_info.material19 as material19,boom_group_info.material2 as material2,boom_group_info.material20 as material20,boom_group_info.material3 as material3,boom_group_info.material4 as material4,boom_group_info.material5 as material5,boom_group_info.material6 as material6,boom_group_info.material7 as material7,boom_group_info.material8 as material8,boom_group_info.material9 as material9,boom_group_info.name as name,boom_group_info.probability1 as probability1,boom_group_info.probability10 as probability10,boom_group_info.probability11 as probability11,boom_group_info.probability12 as probability12,boom_group_info.probability13 as probability13,boom_group_info.probability14 as probability14,boom_group_info.probability15 as probability15,boom_group_info.probability16 as probability16,boom_group_info.probability17 as probability17,boom_group_info.probability18 as probability18,boom_group_info.probability19 as probability19,boom_group_info.probability2 as probability2,boom_group_info.probability20 as probability20,boom_group_info.probability3 as probability3,boom_group_info.probability4 as probability4,boom_group_info.probability5 as probability5,boom_group_info.probability6 as probability6,boom_group_info.probability7 as probability7,boom_group_info.probability8 as probability8,boom_group_info.probability9 as probability9,boom_group_info.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp from boom_group_info 
 		left join boom_group boom_group0 on (boom_group0.token = boom_group_info.boom_group_id)`
	list := []boomGroupInfoModel.BoomGroupInfoFully{}
	limitAndOffset := MakeLimitOffset(curPage, pageSize)
	db := GetDb(tx...).Raw(sql + conditions + order + " " + limitAndOffset).Find(&list)
	return list, db.Error
}

func Count4Page(conditions string, tx ...*gorm.DB) (int64, error) {
	sql := `select count(*) as Count from boom_group_info 
 		left join boom_group boom_group0 on (boom_group0.token = boom_group_info.boom_group_id)`
	count := common.Count{}
	db := GetDb(tx...).Raw(sql + conditions).Take(&count)
	return count.Count, db.Error
}

// 根据id删除 采矿扩展
func DeleteObj(token *string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&boomGroupInfoModel.BoomGroupInfo{Token: token})
	return db.RowsAffected, db.Error
}
