package diggerProductBinding


/**
* 实体类 DiggerProductBinding
*/
type DiggerProductBinding struct {

	TokenOfDigger *string `gorm:"column:digger_id"`
	TokenOfMineProduct *string `gorm:"column:mine_product_id"`                                              
	Name *string                                               
	Token *string `gorm:"primaryKey;"`
}
/**
* 能获取全部属性的实体类 DiggerProductBinding
*/
type DiggerProductBindingFully struct {
	DiggerProductBinding                                                                                                                          
	DiggerDesp *string 
	TokenOfDigger *string                                                                                                                          
	MineProductDesp *string 
	TokenOfMineProduct *string
}
func (this* DiggerProductBindingFully) Convert2Obj() *DiggerProductBinding {
	obj := new(DiggerProductBinding)

	obj.TokenOfDigger = this.TokenOfDigger
	obj.TokenOfMineProduct = this.TokenOfMineProduct
	obj.Name = this.Name
	obj.Token = this.Token
	return obj
}
