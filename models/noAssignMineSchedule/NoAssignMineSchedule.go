package noAssignMineSchedule

import "sim_data_gen/utils"

/**
* 实体类 NoAssignMineSchedule
 */
type NoAssignMineSchedule struct {
	TokenOfBoomGroup *string `gorm:"column:boom_group_id"`
	TokenOfBoomPile  *string `gorm:"column:boom_pile_id"`
	TokenOfDigger    *string `gorm:"column:digger_id"`
	Enabled          *string
	Name             *string
	Priority         *int
	Token            *string `gorm:"primaryKey;"`
}

func (this NoAssignMineSchedule) Clone() *NoAssignMineSchedule {
	dst := new(NoAssignMineSchedule)
	dst.Name = utils.CopyRefString(this.Name)
	dst.TokenOfBoomGroup = utils.CopyRefString(this.TokenOfBoomGroup)
	dst.TokenOfBoomPile = utils.CopyRefString(this.TokenOfBoomPile)
	dst.TokenOfDigger = utils.CopyRefString(this.TokenOfDigger)
	dst.Token = utils.CopyRefString(this.Token)
	dst.Enabled = utils.CopyRefString(this.Enabled)
	dst.Priority = utils.CopyRefInt(this.Priority)
	return dst
}

/**
* 能获取全部属性的实体类 NoAssignMineSchedule
 */
type NoAssignMineScheduleFully struct {
	NoAssignMineSchedule
	BoomGroupDesp    *string
	TokenOfBoomGroup *string
	BoomPileDesp     *string
	TokenOfBoomPile  *string
	DiggerDesp       *string
	TokenOfDigger    *string
}

func (this *NoAssignMineScheduleFully) Convert2Obj() *NoAssignMineSchedule {
	obj := new(NoAssignMineSchedule)

	obj.TokenOfBoomGroup = this.TokenOfBoomGroup
	obj.TokenOfBoomPile = this.TokenOfBoomPile
	obj.TokenOfDigger = this.TokenOfDigger
	obj.Enabled = this.Enabled
	obj.Name = this.Name
	obj.Priority = this.Priority
	obj.Token = this.Token
	return obj
}
