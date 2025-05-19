package boomPile
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	boomPileModel "sim_data_gen/models/boomPile"
	. "sim_data_gen/utils"
)

// 列出 爆堆
func ListObj(tx ... *gorm.DB)([]boomPileModel.BoomPile,error){
	list := []boomPileModel.BoomPile{}
	db := GetDb(tx ...).Table("boom_pile").Find(&list)
	return list,db.Error
}

/**
* boomPile数据库操作类 
*/
func CreateObj(boomPile *boomPileModel.BoomPile,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(boomPile)
	return db.RowsAffected,db.Error
}

//  更新爆堆
func UpdateObj(boomPile *boomPileModel.BoomPile,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(boomPile)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(boomPile *boomPileModel.BoomPile,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(boomPile)
	return db.RowsAffected,db.Error
}

//  更新爆堆
func SaveObj(boomPile *boomPileModel.BoomPile,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(boomPile)
	return db.RowsAffected,db.Error
}


// 根据ID获取爆堆
func GetObjById(token *string,tx ... *gorm.DB)(boomPileModel.BoomPile,error){
	boomPile := boomPileModel.BoomPile{Token:token}
	result := boomPileModel.BoomPile{}
	db := GetDb(tx ...).Where(&boomPile).Take(&result)
	return result,db.Error
}
//  用分页方式列出 爆堆
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]boomPileModel.BoomPile,error){
	sql := `select boom_pile.bench as bench,boom_pile.boom_date as boom_date,boom_pile.cost_to_go as cost_to_go,boom_pile.geom as geom,boom_pile.material1 as material1,boom_pile.material10 as material10,boom_pile.material11 as material11,boom_pile.material12 as material12,boom_pile.material13 as material13,boom_pile.material14 as material14,boom_pile.material15 as material15,boom_pile.material16 as material16,boom_pile.material17 as material17,boom_pile.material18 as material18,boom_pile.material19 as material19,boom_pile.material2 as material2,boom_pile.material20 as material20,boom_pile.material3 as material3,boom_pile.material4 as material4,boom_pile.material5 as material5,boom_pile.material6 as material6,boom_pile.material7 as material7,boom_pile.material8 as material8,boom_pile.material9 as material9,boom_pile.mine_type as mine_type,boom_pile.name as name,boom_pile.nt as nt,boom_pile.quantity as quantity,boom_pile.status as status,boom_pile.tag as tag,boom_pile.token as token,boom_pile.used as used from boom_pile `
	list := []boomPileModel.BoomPile{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from boom_pile `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 爆堆
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&boomPileModel.BoomPile{Token:token})
	return db.RowsAffected,db.Error
}
