package lorry

import (
	"sim_data_gen/models/common"
	lorryModel "sim_data_gen/models/lorry"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 列出 卡车
func ListObj(tx ...*gorm.DB) ([]lorryModel.Lorry, error) {
	list := []lorryModel.Lorry{}
	db := GetDb(tx...).Table("lorry").Find(&list)
	return list, db.Error
}

/**
* lorry数据库操作类
 */
func CreateObj(lorry *lorryModel.Lorry, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(lorry)
	return db.RowsAffected, db.Error
}

// 更新卡车
func UpdateObj(lorry *lorryModel.Lorry, tx ...*gorm.DB) (int64, error) {
	//注意null字段不会更新,只能设置为""
	db := GetDb(tx...).Updates(lorry)
	return db.RowsAffected, db.Error
}

// 更新卡车
func SaveObj(lorry *lorryModel.Lorry, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Save(lorry)
	return db.RowsAffected, db.Error
}

// 根据ID获取卡车
func GetObjByIdFully(token *string, tx ...*gorm.DB) (lorryModel.LorryFully, error) {
	sql := `select lorry.attribs as attribs,lorry.capacity as capacity,lorry.carried as carried,lorry.name as name,lorry.status as status,lorry.token as token,digger0.token as token_of_target_digger,digger0.name as target_digger_desp,boom_group0.token as token_of_target_group,boom_group0.name as target_group_desp,team0.token as token_of_team,team0.name as team_desp from lorry 
 		left join digger digger0 on (digger0.token = lorry.target_digger_id)
 		left join boom_group boom_group0 on (boom_group0.token = lorry.target_group_id)
 		left join team team0 on (team0.token = lorry.team_id) where  lorry.token = ?`
	lorry := lorryModel.LorryFully{}
	db := GetDb(tx...).Raw(sql, token).Take(&lorry)
	return lorry, db.Error
}

// 根据ID获取卡车
func GetObjById(token *string, tx ...*gorm.DB) (lorryModel.Lorry, error) {
	lorry := lorryModel.Lorry{Token: token}
	result := lorryModel.Lorry{}
	db := GetDb(tx...).Where(&lorry).Take(&result)
	return result, db.Error
}

// 用分页方式列出 卡车
func PageObj(conditions string, order string, curPage int, pageSize int, tx ...*gorm.DB) ([]lorryModel.LorryFully, error) {
	sql := `select lorry.attribs as attribs,lorry.capacity as capacity,lorry.carried as carried,lorry.name as name,lorry.status as status,lorry.token as token,digger0.token as token_of_target_digger,digger0.name as target_digger_desp,boom_group0.token as token_of_target_group,boom_group0.name as target_group_desp,team0.token as token_of_team,team0.name as team_desp from lorry 
 		left join digger digger0 on (digger0.token = lorry.target_digger_id)
 		left join boom_group boom_group0 on (boom_group0.token = lorry.target_group_id)
 		left join team team0 on (team0.token = lorry.team_id)`
	list := []lorryModel.LorryFully{}
	limitAndOffset := MakeLimitOffset(curPage, pageSize)
	db := GetDb(tx...).Raw(sql + conditions + order + " " + limitAndOffset).Find(&list)
	return list, db.Error
}

func Count4Page(conditions string, tx ...*gorm.DB) (int64, error) {
	sql := `select count(*) as Count from lorry 
 		left join digger digger0 on (digger0.token = lorry.target_digger_id)
 		left join boom_group boom_group0 on (boom_group0.token = lorry.target_group_id)
 		left join team team0 on (team0.token = lorry.team_id)`
	count := common.Count{}
	db := GetDb(tx...).Raw(sql + conditions).Take(&count)
	return count.Count, db.Error
}

// 根据id删除 卡车
func DeleteObj(token *string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&lorryModel.Lorry{Token: token})
	return db.RowsAffected, db.Error
}
