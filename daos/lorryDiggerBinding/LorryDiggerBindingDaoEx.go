package lorryDiggerBinding

import (
	lorryDiggerBindingModel "sim_data_gen/models/lorryDiggerBinding"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

/**
* lorryDiggerBinding数据库操作类
 */
func SearchMatchedDiggerByToken(token *string, tx ...*gorm.DB) (lorryDiggerBindingModel.LorryDiggerBindingFully, error) {
	sql := `select lorry_digger_binding.enabled as enabled,lorry_digger_binding.level as level,lorry_digger_binding.name as name,lorry_digger_binding.token as token,digger0.token as token_of_digger,digger0.name as digger_desp,lorry0.token as token_of_lorry,lorry0.name as lorry_desp from lorry_digger_binding 
 		left join digger digger0 on (digger0.token = lorry_digger_binding.digger_id)
 		left join lorry lorry0 on (lorry0.token = lorry_digger_binding.lorry_id) 
		where  lorry_digger_binding.lorry_id = ? and lorry_digger_binding.enabled = 1
		limit 1`
	lorryDiggerBindingFully := lorryDiggerBindingModel.LorryDiggerBindingFully{}
	db := GetDb(tx...).Raw(sql, token).Take(&lorryDiggerBindingFully)
	return lorryDiggerBindingFully, db.Error
}
func SearchMatchedDiggerByLorryName(name string, tx ...*gorm.DB) (lorryDiggerBindingModel.LorryDiggerBindingFully, error) {
	sql := `select lorry_digger_binding.enabled as enabled,lorry_digger_binding.level as level,lorry_digger_binding.name as name,lorry_digger_binding.token as token,digger0.token as token_of_digger,digger0.name as digger_desp,lorry0.token as token_of_lorry,lorry0.name as lorry_desp from lorry_digger_binding 
 		left join digger digger0 on (digger0.token = lorry_digger_binding.digger_id)
 		left join lorry lorry0 on (lorry0.token = lorry_digger_binding.lorry_id) 
		where  lorry0.name = ? and lorry_digger_binding.enabled = 1
		limit 1`
	lorryDiggerBindingFully := lorryDiggerBindingModel.LorryDiggerBindingFully{}
	db := GetDb(tx...).Raw(sql, name).Take(&lorryDiggerBindingFully)
	return lorryDiggerBindingFully, db.Error
}
