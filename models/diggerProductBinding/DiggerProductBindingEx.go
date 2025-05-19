package diggerProductBinding

/**
* 能获取全部属性的实体类 DiggerProductBinding
 */
type DiggerProductBindingFully1 struct {
	DiggerProductBinding
	DiggerDesp         *string
	TokenOfDigger      *string
	MineProductDesp    *string
	TokenOfMineProduct *string
	DiggerSpeed        *float64 `gorm:"column:digger_speed"`
}
