package boomGroup

import "sim_data_gen/utils"

/**
* 实体类 BoomGroup
 */
type BoomGroup struct {
	Distance    *float64        `gorm:"default:0;"`
	Geom        *utils.Geometry `gorm:"type:geometry"`
	High        *float64        `gorm:"default:0;"`
	Low         *float64        `gorm:"default:0;"`
	Material1   *float64        `gorm:"default:0;"`
	Material10  *float64        `gorm:"default:0;"`
	Material11  *float64        `gorm:"default:0;"`
	Material12  *float64        `gorm:"default:0;"`
	Material13  *float64        `gorm:"default:0;"`
	Material14  *float64        `gorm:"default:0;"`
	Material15  *float64        `gorm:"default:0;"`
	Material16  *float64        `gorm:"default:0;"`
	Material17  *float64        `gorm:"default:0;"`
	Material18  *float64        `gorm:"default:0;"`
	Material19  *float64        `gorm:"default:0;"`
	Material2   *float64        `gorm:"default:0;"`
	Material20  *float64        `gorm:"default:0;"`
	Material3   *float64        `gorm:"default:0;"`
	Material4   *float64        `gorm:"default:0;"`
	Material5   *float64        `gorm:"default:0;"`
	Material6   *float64        `gorm:"default:0;"`
	Material7   *float64        `gorm:"default:0;"`
	Material8   *float64        `gorm:"default:0;"`
	Material9   *float64        `gorm:"default:0;"`
	Name        *string
	Nt          *string
	Number      *string
	TokenOfPile *string `gorm:"column:pile_id"`
	Status      *string
	Tag         *string
	Token       *string  `gorm:"primaryKey;"`
	Used        *float64 `gorm:"default:0;"`
	X           *float64 `gorm:"default:0;"`
	Y           *float64 `gorm:"default:0;"`
}

/**
* 能获取全部属性的实体类 BoomGroup
 */
type BoomGroupFully struct {
	BoomGroup
	PileDesp    *string
	TokenOfPile *string
}

func (this *BoomGroupFully) Convert2Obj() *BoomGroup {
	obj := new(BoomGroup)

	obj.Distance = this.Distance
	obj.Geom = this.Geom
	obj.High = this.High
	obj.Low = this.Low
	obj.Material1 = this.Material1
	obj.Material10 = this.Material10
	obj.Material11 = this.Material11
	obj.Material12 = this.Material12
	obj.Material13 = this.Material13
	obj.Material14 = this.Material14
	obj.Material15 = this.Material15
	obj.Material16 = this.Material16
	obj.Material17 = this.Material17
	obj.Material18 = this.Material18
	obj.Material19 = this.Material19
	obj.Material2 = this.Material2
	obj.Material20 = this.Material20
	obj.Material3 = this.Material3
	obj.Material4 = this.Material4
	obj.Material5 = this.Material5
	obj.Material6 = this.Material6
	obj.Material7 = this.Material7
	obj.Material8 = this.Material8
	obj.Material9 = this.Material9
	obj.Name = this.Name
	obj.Nt = this.Nt
	obj.Number = this.Number
	obj.TokenOfPile = this.TokenOfPile
	obj.Status = this.Status
	obj.Tag = this.Tag
	obj.Token = this.Token
	obj.Used = this.Used
	obj.X = this.X
	obj.Y = this.Y
	return obj
}

// BoomGroup method helpers (same as before)
func (bg *BoomGroup) Name_() string {
	if bg != nil && bg.Name != nil {
		return *bg.Name
	}
	return "UnknownBoomGroup"
}
func (bg *BoomGroup) X_() float64 {
	if bg != nil && bg.X != nil {
		return *bg.X
	}
	return 0.0
}
func (bg *BoomGroup) Y_() float64 {
	if bg != nil && bg.Y != nil {
		return *bg.Y
	}
	return 0.0
}
func (bg *BoomGroup) High_() float64 {
	if bg != nil && bg.High != nil {
		return *bg.High
	}
	return 0.0
}
func (bg *BoomGroup) Low_() float64 {
	if bg != nil && bg.Low != nil {
		return *bg.Low
	}
	return 0.0
}
