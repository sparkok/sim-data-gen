package boomGroupInfo

import (
	"sim_data_gen/utils"
	"time"
)

type BoomGroupInfoEx struct {
	BoomGroupInfo
	CreatedAt *string
}

func (this *BoomGroupInfoEx) Convert2Obj() *BoomGroupInfo {
	obj := new(BoomGroupInfo)
	if this.CreatedAt != nil {
		timeValue, err := time.ParseInLocation(time.RFC3339, *this.CreatedAt, utils.SystemZone)
		if err != nil {
			panic(err)
		}
		obj.CreatedAt = utils.RefTime(timeValue)
	}

	//time.ParseInLocation()
	//obj.CreatedAt = "this.CreatedAt"
	obj.TokenOfBoomGroup = this.TokenOfBoomGroup
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
	obj.Probability1 = this.Probability1
	obj.Probability10 = this.Probability10
	obj.Probability11 = this.Probability11
	obj.Probability12 = this.Probability12
	obj.Probability13 = this.Probability13
	obj.Probability14 = this.Probability14
	obj.Probability15 = this.Probability15
	obj.Probability16 = this.Probability16
	obj.Probability17 = this.Probability17
	obj.Probability18 = this.Probability18
	obj.Probability19 = this.Probability19
	obj.Probability2 = this.Probability2
	obj.Probability20 = this.Probability20
	obj.Probability3 = this.Probability3
	obj.Probability4 = this.Probability4
	obj.Probability5 = this.Probability5
	obj.Probability6 = this.Probability6
	obj.Probability7 = this.Probability7
	obj.Probability8 = this.Probability8
	obj.Probability9 = this.Probability9
	obj.Token = this.Token
	return obj
}

/**
* 实体类 BoomGroupInfo
 */
type BoomGroupInfo struct {
	TokenOfBoomGroup *string `gorm:"column:boom_group_id"`
	CreatedAt        *time.Time
	Material1        *float64 `gorm:"default:0;"`
	Material10       *float64 `gorm:"default:0;"`
	Material11       *float64 `gorm:"default:0;"`
	Material12       *float64 `gorm:"default:0;"`
	Material13       *float64 `gorm:"default:0;"`
	Material14       *float64 `gorm:"default:0;"`
	Material15       *float64 `gorm:"default:0;"`
	Material16       *float64 `gorm:"default:0;"`
	Material17       *float64 `gorm:"default:0;"`
	Material18       *float64 `gorm:"default:0;"`
	Material19       *float64 `gorm:"default:0;"`
	Material2        *float64 `gorm:"default:0;"`
	Material20       *float64 `gorm:"default:0;"`
	Material3        *float64 `gorm:"default:0;"`
	Material4        *float64 `gorm:"default:0;"`
	Material5        *float64 `gorm:"default:0;"`
	Material6        *float64 `gorm:"default:0;"`
	Material7        *float64 `gorm:"default:0;"`
	Material8        *float64 `gorm:"default:0;"`
	Material9        *float64 `gorm:"default:0;"`
	Name             *string
	Probability1     *float64 `gorm:"default:0;"`
	Probability10    *float64 `gorm:"default:0;"`
	Probability11    *float64 `gorm:"default:0;"`
	Probability12    *float64 `gorm:"default:0;"`
	Probability13    *float64 `gorm:"default:0;"`
	Probability14    *float64 `gorm:"default:0;"`
	Probability15    *float64 `gorm:"default:0;"`
	Probability16    *float64 `gorm:"default:0;"`
	Probability17    *float64 `gorm:"default:0;"`
	Probability18    *float64 `gorm:"default:0;"`
	Probability19    *float64 `gorm:"default:0;"`
	Probability2     *float64 `gorm:"default:0;"`
	Probability20    *float64 `gorm:"default:0;"`
	Probability3     *float64 `gorm:"default:0;"`
	Probability4     *float64 `gorm:"default:0;"`
	Probability5     *float64 `gorm:"default:0;"`
	Probability6     *float64 `gorm:"default:0;"`
	Probability7     *float64 `gorm:"default:0;"`
	Probability8     *float64 `gorm:"default:0;"`
	Probability9     *float64 `gorm:"default:0;"`
	Token            *string  `gorm:"primaryKey;"`
}

/**
* 能获取全部属性的实体类 BoomGroupInfo
 */
type BoomGroupInfoFully struct {
	BoomGroupInfo
	BoomGroupDesp    *string
	TokenOfBoomGroup *string
}

func (this *BoomGroupInfoFully) Convert2Obj() *BoomGroupInfo {
	obj := new(BoomGroupInfo)

	obj.TokenOfBoomGroup = this.TokenOfBoomGroup
	obj.CreatedAt = this.CreatedAt
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
	obj.Probability1 = this.Probability1
	obj.Probability10 = this.Probability10
	obj.Probability11 = this.Probability11
	obj.Probability12 = this.Probability12
	obj.Probability13 = this.Probability13
	obj.Probability14 = this.Probability14
	obj.Probability15 = this.Probability15
	obj.Probability16 = this.Probability16
	obj.Probability17 = this.Probability17
	obj.Probability18 = this.Probability18
	obj.Probability19 = this.Probability19
	obj.Probability2 = this.Probability2
	obj.Probability20 = this.Probability20
	obj.Probability3 = this.Probability3
	obj.Probability4 = this.Probability4
	obj.Probability5 = this.Probability5
	obj.Probability6 = this.Probability6
	obj.Probability7 = this.Probability7
	obj.Probability8 = this.Probability8
	obj.Probability9 = this.Probability9
	obj.Token = this.Token
	return obj
}
