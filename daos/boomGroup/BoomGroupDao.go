package boomGroup

import (
	boomGroupModel "sim_data_gen/models/boomGroup"
	"sim_data_gen/models/common"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

// 列出 采矿单元
func ListObj(tx ...*gorm.DB) ([]boomGroupModel.BoomGroup, error) {
	list := []boomGroupModel.BoomGroup{}
	db := GetDb(tx...).Table("boom_group").Find(&list)
	return list, db.Error
}

/**
* boomGroup数据库操作类
 */
func CreateObj(boomGroup *boomGroupModel.BoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Create(boomGroup)
	return db.RowsAffected, db.Error
}

// 更新采矿单元
func UpdateObj(boomGroup *boomGroupModel.BoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Updates(boomGroup)
	return db.RowsAffected, db.Error
}

// 更新非null字段
func UpdateObjItem(boomGroup *boomGroupModel.BoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).UpdateColumns(boomGroup)
	return db.RowsAffected, db.Error
}

// 更新采矿单元
func SaveObj(boomGroup *boomGroupModel.BoomGroup, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Save(boomGroup)
	return db.RowsAffected, db.Error
}

// 根据ID获取采矿单元
func GetObjByIdFully(token *string, tx ...*gorm.DB) (boomGroupModel.BoomGroupFully, error) {
	sql := `select boom_group.distance as distance,boom_group.geom as geom,boom_group.high as high,boom_group.low as low,boom_group.material1 as material1,boom_group.material10 as material10,boom_group.material11 as material11,boom_group.material12 as material12,boom_group.material13 as material13,boom_group.material14 as material14,boom_group.material15 as material15,boom_group.material16 as material16,boom_group.material17 as material17,boom_group.material18 as material18,boom_group.material19 as material19,boom_group.material2 as material2,boom_group.material20 as material20,boom_group.material3 as material3,boom_group.material4 as material4,boom_group.material5 as material5,boom_group.material6 as material6,boom_group.material7 as material7,boom_group.material8 as material8,boom_group.material9 as material9,boom_group.name as name,boom_group.nt as nt,boom_group.number as number,boom_group.status as status,boom_group.tag as tag,boom_group.token as token,boom_group.used as used,boom_group.x as x,boom_group.y as y,boom_pile0.token as token_of_pile,boom_pile0.name as pile_desp from boom_group 
 		left join boom_pile boom_pile0 on (boom_pile0.token = boom_group.pile_id) where  boom_group.token = ?`
	boomGroup := boomGroupModel.BoomGroupFully{}
	db := GetDb(tx...).Raw(sql, token).Take(&boomGroup)
	return boomGroup, db.Error
}

// 根据ID获取采矿单元
func GetObjById(token *string, tx ...*gorm.DB) (boomGroupModel.BoomGroup, error) {
	boomGroup := boomGroupModel.BoomGroup{Token: token}
	result := boomGroupModel.BoomGroup{}
	db := GetDb(tx...).Where(&boomGroup).Take(&result)
	return result, db.Error
}

// 用分页方式列出 采矿单元
func PageObj(conditions string, order string, curPage int, pageSize int, tx ...*gorm.DB) ([]boomGroupModel.BoomGroupFully, error) {
	sql := `select boom_group.distance as distance,boom_group.geom as geom,boom_group.high as high,boom_group.low as low,boom_group.material1 as material1,boom_group.material10 as material10,boom_group.material11 as material11,boom_group.material12 as material12,boom_group.material13 as material13,boom_group.material14 as material14,boom_group.material15 as material15,boom_group.material16 as material16,boom_group.material17 as material17,boom_group.material18 as material18,boom_group.material19 as material19,boom_group.material2 as material2,boom_group.material20 as material20,boom_group.material3 as material3,boom_group.material4 as material4,boom_group.material5 as material5,boom_group.material6 as material6,boom_group.material7 as material7,boom_group.material8 as material8,boom_group.material9 as material9,boom_group.name as name,boom_group.nt as nt,boom_group.number as number,boom_group.status as status,boom_group.tag as tag,boom_group.token as token,boom_group.used as used,boom_group.x as x,boom_group.y as y,boom_pile0.token as token_of_pile,boom_pile0.name as pile_desp from boom_group 
 		left join boom_pile boom_pile0 on (boom_pile0.token = boom_group.pile_id)`
	list := []boomGroupModel.BoomGroupFully{}
	limitAndOffset := MakeLimitOffset(curPage, pageSize)
	db := GetDb(tx...).Raw(sql + conditions + order + " " + limitAndOffset).Find(&list)
	return list, db.Error
}

func Count4Page(conditions string, tx ...*gorm.DB) (int64, error) {
	sql := `select count(*) as Count from boom_group 
 		left join boom_pile boom_pile0 on (boom_pile0.token = boom_group.pile_id)`
	count := common.Count{}
	db := GetDb(tx...).Raw(sql + conditions).Take(&count)
	return count.Count, db.Error
}

// 根据id删除 采矿单元
func DeleteObj(token *string, tx ...*gorm.DB) (int64, error) {
	db := GetDb(tx...).Delete(&boomGroupModel.BoomGroup{Token: token})
	return db.RowsAffected, db.Error
}
