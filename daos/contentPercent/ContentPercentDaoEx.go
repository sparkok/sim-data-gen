package contentPercent

import (
	contentPercentModel "sim_data_gen/models/contentPercent"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

//	用分页方式列出 品位
//
// 根据ID获取品位
func ListObjByIdFully(tokenOfMineProduct *string, tx ...*gorm.DB) (contentPercents []contentPercentModel.ContentPercentFully, err error) {
	sql := `select content_percent.max_value as max_value,content_percent.min_value as min_value,content_percent.name as name,content_percent.status as status,content_percent.token as token,content_percent.num as num,composition0.token as token_of_composition,composition0.name as composition_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from content_percent 
 		left join composition composition0 on (composition0.token = content_percent.composition_id)
 		left join mine_product mine_product0 on (mine_product0.token = content_percent.mine_product_id) where  content_percent.mine_product_id = ?`
	db := GetDb(tx...).Raw(sql, tokenOfMineProduct).Find(&contentPercents)
	return contentPercents, db.Error
}
