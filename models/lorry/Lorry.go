package lorry

/**
* 实体类 Lorry
 */
type Lorry struct {
	Attribs             *string
	Capacity            *float64
	Carried             *float64
	Name                *string
	Status              *int
	TokenOfTargetDigger *string `gorm:"column:target_digger_id"`
	TokenOfTargetGroup  *string `gorm:"column:target_group_id"`
	TokenOfTeam *string `gorm:"column:team_id"` 
	Token               *string `gorm:"primaryKey;"`
}
type LorryWithComm struct {
	Attribs             *string
	Capacity            *float64
	Carried             *float64
	Name                *string
	Status              *int
	Utc                 *int
	LastCommUtc         *int
	TokenOfTargetDigger *string `gorm:"column:target_digger_id"`
	TokenOfTargetGroup  *string `gorm:"column:target_group_id"`
	TokenOfTeam *string `gorm:"column:team_id"`
	Token               *string `gorm:"primaryKey;"`
}
type LorryWithSearchInfo struct {
	Attribs             *string
	Capacity            *float64
	Carried             *float64
	Name                *string
	Status              *int
	TokenOfTargetDigger *string `gorm:"column:target_digger_id"`
	TokenOfTargetGroup  *string `gorm:"column:target_group_id"`
	TokenOfTeam *string `gorm:"column:team_id"`
	Token               *string `gorm:"primaryKey;"`
	SearchDistance      *float64
	SearchUtc           *int
	DiggerX             *float64
	DiggerY             *float64
	GpsUtc              *int
	LorryX              *float64
	LorryY              *float64
}

/**
* 能获取全部属性的实体类 Lorry
 */
type LorryFully struct {
	Lorry
	TargetDiggerDesp    *string
	TokenOfTargetDigger *string
	TargetGroupDesp     *string
	TokenOfTargetGroup  *string
	TeamDesp *string 
	TokenOfTeam *string
}

func (this *LorryFully) Convert2Obj() *Lorry {
	obj := new(Lorry)

	obj.Attribs = this.Attribs
	obj.Capacity = this.Capacity
	obj.Carried = this.Carried
	obj.Name = this.Name
	obj.Status = this.Status
	obj.TokenOfTargetDigger = this.TokenOfTargetDigger
	obj.TokenOfTargetGroup = this.TokenOfTargetGroup
	obj.TokenOfTeam = this.TokenOfTeam
	obj.Token = this.Token
	return obj
}
func Status2Txt(status *int) string {
	if status == nil {
		return "未知"
	}
	switch *status {
	case 0:
		return "休息"
	case 1:
		return "空车"
	case 2:
		return "载货"
	case 3:
		return "已接单"
	default:
		return "其他"
	}

}
