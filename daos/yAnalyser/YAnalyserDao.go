package yAnalyser
import (
	"sim_data_gen/models/common"
	"gorm.io/gorm"
	yAnalyserModel "sim_data_gen/models/yAnalyser"
	. "sim_data_gen/utils"
)

// 列出 分析仪相关数据
func ListObj(tx ... *gorm.DB)([]yAnalyserModel.YAnalyser,error){
	list := []yAnalyserModel.YAnalyser{}
	db := GetDb(tx ...).Table("y_analyser").Find(&list)
	return list,db.Error
}

/**
* yAnalyser数据库操作类 
*/
func CreateObj(yAnalyser *yAnalyserModel.YAnalyser,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Create(yAnalyser)
	return db.RowsAffected,db.Error
}

//  更新分析仪相关数据
func UpdateObj(yAnalyser *yAnalyserModel.YAnalyser,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Updates(yAnalyser)
	return db.RowsAffected,db.Error
}

//  更新非null字段
func UpdateObjItem(yAnalyser *yAnalyserModel.YAnalyser,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).UpdateColumns(yAnalyser)
	return db.RowsAffected,db.Error
}

//  更新分析仪相关数据
func SaveObj(yAnalyser *yAnalyserModel.YAnalyser,tx ... *gorm.DB)(int64,error){
	db := GetDb(tx ...).Save(yAnalyser)
	return db.RowsAffected,db.Error
}


// 根据ID获取分析仪相关数据
func GetObjById(token *string,tx ... *gorm.DB)(yAnalyserModel.YAnalyser,error){
	yAnalyser := yAnalyserModel.YAnalyser{Token:token}
	result := yAnalyserModel.YAnalyser{}
	db := GetDb(tx ...).Where(&yAnalyser).Take(&result)
	return result,db.Error
}
//  用分页方式列出 分析仪相关数据
func PageObj(conditions string,order string,curPage int,pageSize int,tx ... *gorm.DB)([]yAnalyserModel.YAnalyser,error){
	sql := `select y_analyser.analyser_num as analyser_num,y_analyser.created_at as created_at,y_analyser.crushing_plant as crushing_plant,y_analyser.flux as flux,y_analyser.load as load,y_analyser.mat1 as mat1,y_analyser.mat10 as mat10,y_analyser.mat11 as mat11,y_analyser.mat12 as mat12,y_analyser.mat13 as mat13,y_analyser.mat14 as mat14,y_analyser.mat15 as mat15,y_analyser.mat16 as mat16,y_analyser.mat17 as mat17,y_analyser.mat18 as mat18,y_analyser.mat19 as mat19,y_analyser.mat2 as mat2,y_analyser.mat20 as mat20,y_analyser.mat3 as mat3,y_analyser.mat4 as mat4,y_analyser.mat5 as mat5,y_analyser.mat6 as mat6,y_analyser.mat7 as mat7,y_analyser.mat8 as mat8,y_analyser.mat9 as mat9,y_analyser.speed as speed,y_analyser.status as status,y_analyser.test_at as test_at,y_analyser.token as token from y_analyser `
	list := []yAnalyserModel.YAnalyser{}
	limitAndOffset := MakeLimitOffset(curPage,pageSize)
	db := GetDb(tx ...).Raw(sql + conditions + " " + order + " " + limitAndOffset).Find(&list)
	return list,db.Error
}

func Count4Page(conditions string,tx ... *gorm.DB)(int64,error){
	sql := `select count(*) as Count from y_analyser `
	count := common.Count{}
	db := GetDb(tx ...).Raw(sql + conditions).Take(&count)
	return count.Count,db.Error
}

// 根据id删除 分析仪相关数据
func DeleteObj(token *string,tx ... *gorm.DB) (int64,error){
	db := GetDb(tx ...).Delete(&yAnalyserModel.YAnalyser{Token:token})
	return db.RowsAffected,db.Error
}
