package contentPercent


/**
* 实体类 ContentPercent
*/
type ContentPercent struct {

	TokenOfComposition *string `gorm:"column:composition_id"`                                              
	MaxValue *float64 `gorm:"default:0;"`                                              
	MinValue *float64 `gorm:"default:0;"`
	TokenOfMineProduct *string `gorm:"column:mine_product_id"`                                              
	Name *string                                               
	Num *int `gorm:"default:0;"`                                              
	Status *int `gorm:"default:0;"`                                              
	Token *string `gorm:"primaryKey;"`
}
/**
* 能获取全部属性的实体类 ContentPercent
*/
type ContentPercentFully struct {
	ContentPercent                                                                                                                          
	CompositionDesp *string 
	TokenOfComposition *string                                                                                                                          
	MineProductDesp *string 
	TokenOfMineProduct *string
}
func (this* ContentPercentFully) Convert2Obj() *ContentPercent {
	obj := new(ContentPercent)

	obj.TokenOfComposition = this.TokenOfComposition
	obj.MaxValue = this.MaxValue
	obj.MinValue = this.MinValue
	obj.TokenOfMineProduct = this.TokenOfMineProduct
	obj.Name = this.Name
	obj.Num = this.Num
	obj.Status = this.Status
	obj.Token = this.Token
	return obj
}
