package diggerProductBinding

import (
	diggerProductBindingModel "sim_data_gen/models/diggerProductBinding"
	. "sim_data_gen/utils"
	"fmt"
	"gorm.io/gorm"
)

// 列出 挖机产品绑定
func ListDiggerOfMineProductExt(tokenOfMineProduct string, includeNoBoundDigger bool, tx ...*gorm.DB) ([]diggerProductBindingModel.DiggerProductBindingFully1, error) {
	sql := `select digger_product_binding.name as name,
       			digger_product_binding.token as token,
       			digger.token as token_of_digger,
       			digger.name as digger_desp,
       			mine_product.token as token_of_mine_product,
       			mine_product.name as mine_product_desp,
       			digger.speed as digger_speed 
			from digger 
 			left join digger_product_binding on (digger.token = digger_product_binding.digger_id)
 			left join mine_product on (mine_product.token = digger_product_binding.mine_product_id)`
	if includeNoBoundDigger {
		sql += fmt.Sprintf(` where ( digger.status = 1 ) and ( digger_product_binding.mine_product_id = '%s' or digger_product_binding.mine_product_id is null)`, tokenOfMineProduct)
	} else {
		sql += fmt.Sprintf(` where  ( digger.status = 1 ) and ( digger_product_binding.mine_product_id = '%s' )`, tokenOfMineProduct)
	}
	sql += ` order by digger_desp`
	list := []diggerProductBindingModel.DiggerProductBindingFully1{}
	db := GetDb(tx...).Raw(sql).Find(&list)
	return list, db.Error
}
