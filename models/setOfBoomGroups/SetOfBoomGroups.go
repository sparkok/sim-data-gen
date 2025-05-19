package setOfBoomGroups

import "time"

/**
* 实体类 SetOfBoomGroups
 */
type SetOfBoomGroups struct {
	BoomGroupIds       *string
	CreatedAt          *time.Time
	DateFlag           *string
	Diggers            *string
	MatContents        *string
	TokenOfMineProduct *string `gorm:"column:mine_product_id"`
	Name               *string
	Nt                 *string
	Status             *int    `gorm:"default:0;"`
	Token              *string `gorm:"primaryKey;"`
	UpdateAt           *time.Time
}

/**
* 能获取全部属性的实体类 SetOfBoomGroups
 */
type SetOfBoomGroupsFully struct {
	SetOfBoomGroups
	MineProductDesp    *string
	TokenOfMineProduct *string
}

func (this *SetOfBoomGroupsFully) Convert2Obj() *SetOfBoomGroups {
	obj := new(SetOfBoomGroups)

	obj.BoomGroupIds = this.BoomGroupIds
	obj.CreatedAt = this.CreatedAt
	obj.DateFlag = this.DateFlag
	obj.Diggers = this.Diggers
	obj.MatContents = this.MatContents
	obj.TokenOfMineProduct = this.TokenOfMineProduct
	obj.Name = this.Name
	obj.Nt = this.Nt
	obj.Status = this.Status
	obj.Token = this.Token
	obj.UpdateAt = this.UpdateAt
	return obj
}
