package lorryDiggerBindingLog

/**
* 实体类 LorryDiggerBindingLog
 */
type LorryDiggerBindingLog struct {
	TokenOfBoomGroup *string `gorm:"column:boom_group_id"`
	BoomGroupName    *string
	DateFlag         *string
	TokenOfDigger    *string `gorm:"column:digger_id"`
	DiggerName       *string
	EndUtc           *int
	TokenOfLorry     *string `gorm:"column:lorry_id"`
	LorryName        *string
	StartUtc         *int
	Token            *string `gorm:"primaryKey;"`
}

func (this LorryDiggerBindingLog) Clone() LorryDiggerBindingLog {
	cloneObj := LorryDiggerBindingLog{}
	cloneObj = this
	return cloneObj
}

/**
* 能获取全部属性的实体类 LorryDiggerBindingLog
 */
type LorryDiggerBindingLogFully struct {
	LorryDiggerBindingLog
	BoomGroupDesp    *string
	TokenOfBoomGroup *string
	DiggerDesp       *string
	TokenOfDigger    *string
	LorryDesp        *string
	TokenOfLorry     *string
}

func (this *LorryDiggerBindingLogFully) Convert2Obj() *LorryDiggerBindingLog {
	obj := new(LorryDiggerBindingLog)

	obj.TokenOfBoomGroup = this.TokenOfBoomGroup
	obj.BoomGroupName = this.BoomGroupName
	obj.DateFlag = this.DateFlag
	obj.TokenOfDigger = this.TokenOfDigger
	obj.DiggerName = this.DiggerName
	obj.StartUtc = this.StartUtc
	obj.TokenOfLorry = this.TokenOfLorry
	obj.LorryName = this.LorryName
	obj.Token = this.Token
	obj.EndUtc = this.EndUtc
	return obj
}
