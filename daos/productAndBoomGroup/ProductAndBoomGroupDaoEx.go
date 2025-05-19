package productAndBoomGroup

import (
	productAndBoomGroupModel "sim_data_gen/models/productAndBoomGroup"
	. "sim_data_gen/utils"
	"gorm.io/gorm"
)

func ListByFilter(conditions string, order string, tx ...*gorm.DB) ([]productAndBoomGroupModel.ProductAndBoomGroupFully, error) {
	sql := `select product_and_boom_group.name as name,product_and_boom_group.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from product_and_boom_group 
 		left join boom_group boom_group0 on (boom_group0.token = product_and_boom_group.boom_group_id)
 		left join mine_product mine_product0 on (mine_product0.token = product_and_boom_group.mine_product_id)`
	list := []productAndBoomGroupModel.ProductAndBoomGroupFully{}
	db := GetDb(tx...).Raw(sql + conditions + order).Find(&list)
	return list, db.Error
}

// 用分页方式列出 产品和配矿单元
func ListEnabledObjByProductName(productName string, tx ...*gorm.DB) ([]productAndBoomGroupModel.ProductAndBoomGroupFully, error) {
	var order = ""
	//只查询和产品相关的数据
	var conditions = " where boom_group.status = \"1\" and mine_product.name = \"" + productName + "\""
	sql := `select product_and_boom_group.name as name,product_and_boom_group.token as token,boom_group0.token as token_of_boom_group,boom_group0.name as boom_group_desp,mine_product0.token as token_of_mine_product,mine_product0.name as mine_product_desp from product_and_boom_group 
 		left join boom_group on (boom_group.token = product_and_boom_group.boom_group_id)
 		left join mine_product on (mine_product.token = product_and_boom_group.mine_product_id)`
	//left join boom_group_info on (boom_group.token = boom_group_info.token.boom_group_id)
	list := []productAndBoomGroupModel.ProductAndBoomGroupFully{}
	db := GetDb(tx...).Raw(sql + conditions + order).Find(&list)
	return list, db.Error
}
func ListEnabledObjByProductToken(productToken string, includeUnboundBoomGroup bool, tx ...*gorm.DB) ([]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, error) {
	var order = ""
	//查询产品的配矿单元的数据
	//注意:status = '0' 表示使用中, '1' 表示已停用
	var conditions string
	if includeUnboundBoomGroup {
		conditions = " where boom_group.status = '0' and ( mine_product.token = '" + productToken + "' or mine_product.token is null )"
	} else {
		conditions = " where boom_group.status = '0' and ( mine_product.token = '" + productToken + "' )"
	}
	sql := `select 	product_and_boom_group.name as name,
       				product_and_boom_group.token as token,
       				boom_group.token as token_of_boom_group,
       				boom_group.name as boom_group_desp,
       				mine_product.token as token_of_mine_product,
       				mine_product.name as mine_product_desp,
					boom_group.distance as distance,
					boom_group.geom as geom,
					boom_group.high as high,
					boom_group.low as low,
					boom_group.material1 as material1,
					boom_group.material10 as material10,
					boom_group.material11 as material11,
					boom_group.material12 as material12,
					boom_group.material13 as material13,
					boom_group.material14 as material14,
					boom_group.material15 as material15,
					boom_group.material16 as material16,
					boom_group.material17 as material17,
					boom_group.material18 as material18,
					boom_group.material19 as material19,
					boom_group.material2 as material2,
					boom_group.material20 as material20,
					boom_group.material3 as material3,
					boom_group.material4 as material4,
					boom_group.material5 as material5,
					boom_group.material6 as material6,
					boom_group.material7 as material7,
					boom_group.material8 as material8,
					boom_group.material9 as material9,
					boom_group.name as name,
					boom_group.nt as nt,
					boom_group.number as number,
					boom_group.status as status,
					boom_group.tag as tag,
					boom_group.token as token,
					boom_group.used as used,
					boom_group.x as x,
					boom_group.y as y    		
			from boom_group
			left join product_and_boom_group on (boom_group.token = product_and_boom_group.boom_group_id)
			left join mine_product on (mine_product.token = product_and_boom_group.mine_product_id)
			`
	//left join boom_group_info on (boom_group.token = boom_group_info.token.boom_group_id)
	list := []productAndBoomGroupModel.ProductAndBoomGroupInDetailFully{}
	db := GetDb(tx...).Raw(sql + conditions + order).Find(&list)
	return list, db.Error
}
