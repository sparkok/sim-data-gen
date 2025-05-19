package lorryDiggerBinding

/**
* 实体类 LorryDiggerBinding
 */
type LorryDiggerBinding struct {
	TokenOfDigger *string `gorm:"column:digger_id"`
	Enabled       *int
	Level         *int
	TokenOfLorry  *string `gorm:"column:lorry_id"`
	Name          *string
	Token         *string `gorm:"primaryKey;"`
}

/**
* 能获取全部属性的实体类 LorryDiggerBinding
 */
type LorryDiggerBindingFully struct {
	LorryDiggerBinding
	DiggerDesp    *string
	TokenOfDigger *string
	LorryDesp     *string
	TokenOfLorry  *string
}

func (this *LorryDiggerBindingFully) Convert2Obj() *LorryDiggerBinding {
	obj := new(LorryDiggerBinding)

	obj.TokenOfDigger = this.TokenOfDigger
	obj.Enabled = this.Enabled
	obj.Level = this.Level
	obj.TokenOfLorry = this.TokenOfLorry
	obj.Name = this.Name
	obj.Token = this.Token
	return obj
}
