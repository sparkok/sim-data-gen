package productAndBoomGroup

import "sim_data_gen/models/boomGroup"

/**
* 实体类 ProductAndBoomGroup
 */
type ProductAndBoomGroup struct {
	TokenOfBoomGroup   *string `gorm:"column:boom_group_id"`
	TokenOfMineProduct *string `gorm:"column:mine_product_id"`
	Name               *string
	Token              *string `gorm:"primaryKey;"`
}

/**
* 能获取全部属性的实体类 ProductAndBoomGroup
 */
type ProductAndBoomGroupFully struct {
	ProductAndBoomGroup
	BoomGroupDesp      *string
	TokenOfBoomGroup   *string
	MineProductDesp    *string
	TokenOfMineProduct *string
}

type ProductAndBoomGroupInDetailFully struct {
	boomGroup.BoomGroup
	ProductAndBoomGroup
	BoomGroupDesp      *string
	TokenOfBoomGroup   *string
	MineProductDesp    *string
	TokenOfMineProduct *string
}

func (this *ProductAndBoomGroupFully) Convert2Obj() *ProductAndBoomGroup {
	obj := new(ProductAndBoomGroup)

	obj.TokenOfBoomGroup = this.TokenOfBoomGroup
	obj.TokenOfMineProduct = this.TokenOfMineProduct
	obj.Name = this.Name
	obj.Token = this.Token
	return obj
}
