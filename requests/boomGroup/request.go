package boomGroup

import (
	"bytes"
	boomGroupModel "sim_data_gen/models/boomGroup"
	"sim_data_gen/requests/common"
	"sim_data_gen/utils"
)

/**
* BoomGroup请求类
 */
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {

	//`json:"distance"`
	Distance *string
	//`json:"geom"`
	Geom *string
	//`json:"high"`
	High *string
	//`json:"low"`
	Low *string
	//`json:"material1"`
	Material1 *string
	//`json:"material10"`
	Material10 *string
	//`json:"material11"`
	Material11 *string
	//`json:"material12"`
	Material12 *string
	//`json:"material13"`
	Material13 *string
	//`json:"material14"`
	Material14 *string
	//`json:"material15"`
	Material15 *string
	//`json:"material16"`
	Material16 *string
	//`json:"material17"`
	Material17 *string
	//`json:"material18"`
	Material18 *string
	//`json:"material19"`
	Material19 *string
	//`json:"material2"`
	Material2 *string
	//`json:"material20"`
	Material20 *string
	//`json:"material3"`
	Material3 *string
	//`json:"material4"`
	Material4 *string
	//`json:"material5"`
	Material5 *string
	//`json:"material6"`
	Material6 *string
	//`json:"material7"`
	Material7 *string
	//`json:"material8"`
	Material8 *string
	//`json:"material9"`
	Material9 *string
	//`json:"name"`
	Name *string
	//`json:"nt"`
	Nt *string
	//`json:"number"`
	Number *string
	//`json:"tokenOfPile"`
	TokenOfPile *string

	//`json:"status"`
	Status *string
	//`json:"tag"`
	Tag *string
	//`json:"token"`
	Token *string
	//`json:"used"`
	Used *string
	//`json:"x"`
	X *string
	//`json:"y"`
	Y *string
}

func (this *SearchInfo) GetConditions() string {
	var condition bytes.Buffer
	if this.Distance != nil {
		condition.WriteString("and (boom_group.distance = " + *this.Distance + ")")
	}
	if this.Geom != nil {
		condition.WriteString("and (boom_group.geom like '%" + *this.Geom + "%')")
	}
	if this.High != nil {
		condition.WriteString("and (boom_group.high = " + *this.High + ")")
	}
	if this.Low != nil {
		condition.WriteString("and (boom_group.low = " + *this.Low + ")")
	}
	if this.Material1 != nil {
		condition.WriteString("and (boom_group.material1 = " + *this.Material1 + ")")
	}
	if this.Material10 != nil {
		condition.WriteString("and (boom_group.material10 = " + *this.Material10 + ")")
	}
	if this.Material11 != nil {
		condition.WriteString("and (boom_group.material11 = " + *this.Material11 + ")")
	}
	if this.Material12 != nil {
		condition.WriteString("and (boom_group.material12 = " + *this.Material12 + ")")
	}
	if this.Material13 != nil {
		condition.WriteString("and (boom_group.material13 = " + *this.Material13 + ")")
	}
	if this.Material14 != nil {
		condition.WriteString("and (boom_group.material14 = " + *this.Material14 + ")")
	}
	if this.Material15 != nil {
		condition.WriteString("and (boom_group.material15 = " + *this.Material15 + ")")
	}
	if this.Material16 != nil {
		condition.WriteString("and (boom_group.material16 = " + *this.Material16 + ")")
	}
	if this.Material17 != nil {
		condition.WriteString("and (boom_group.material17 = " + *this.Material17 + ")")
	}
	if this.Material18 != nil {
		condition.WriteString("and (boom_group.material18 = " + *this.Material18 + ")")
	}
	if this.Material19 != nil {
		condition.WriteString("and (boom_group.material19 = " + *this.Material19 + ")")
	}
	if this.Material2 != nil {
		condition.WriteString("and (boom_group.material2 = " + *this.Material2 + ")")
	}
	if this.Material20 != nil {
		condition.WriteString("and (boom_group.material20 = " + *this.Material20 + ")")
	}
	if this.Material3 != nil {
		condition.WriteString("and (boom_group.material3 = " + *this.Material3 + ")")
	}
	if this.Material4 != nil {
		condition.WriteString("and (boom_group.material4 = " + *this.Material4 + ")")
	}
	if this.Material5 != nil {
		condition.WriteString("and (boom_group.material5 = " + *this.Material5 + ")")
	}
	if this.Material6 != nil {
		condition.WriteString("and (boom_group.material6 = " + *this.Material6 + ")")
	}
	if this.Material7 != nil {
		condition.WriteString("and (boom_group.material7 = " + *this.Material7 + ")")
	}
	if this.Material8 != nil {
		condition.WriteString("and (boom_group.material8 = " + *this.Material8 + ")")
	}
	if this.Material9 != nil {
		condition.WriteString("and (boom_group.material9 = " + *this.Material9 + ")")
	}
	if this.Name != nil {
		condition.WriteString("and (boom_group.name like '%" + *this.Name + "%')")
	}
	if this.Nt != nil {
		condition.WriteString("and (boom_group.nt like '%" + *this.Nt + "%')")
	}
	if this.Number != nil {
		condition.WriteString("and (boom_group.number like '%" + *this.Number + "%')")
	}
	if this.TokenOfPile != nil {
		condition.WriteString("and (boom_group.pile_id = '" + *this.TokenOfPile + "')")
	}
	if this.Status != nil {
		condition.WriteString("and (boom_group.status = '" + *this.Status + "')")
	}
	if this.Tag != nil {
		condition.WriteString("and (boom_group.tag like '%" + *this.Tag + "%')")
	}
	if this.Token != nil {
		condition.WriteString("and (boom_group.token like '%" + *this.Token + "%')")
	}
	if this.Used != nil {
		condition.WriteString("and (boom_group.used = " + *this.Used + ")")
	}
	if this.X != nil {
		condition.WriteString("and (boom_group.x = " + *this.X + ")")
	}
	if this.Y != nil {
		condition.WriteString("and (boom_group.y = " + *this.Y + ")")
	}

	if condition.Len() > 4 {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"distance"`
	Distance *float64
	//`json:"geom"`
	Geom *utils.Geometry `gorm:"type:geometry"`
	//`json:"high"`
	High *float64
	//`json:"low"`
	Low *float64
	//`json:"material1"`
	Material1 *float64
	//`json:"material10"`
	Material10 *float64
	//`json:"material11"`
	Material11 *float64
	//`json:"material12"`
	Material12 *float64
	//`json:"material13"`
	Material13 *float64
	//`json:"material14"`
	Material14 *float64
	//`json:"material15"`
	Material15 *float64
	//`json:"material16"`
	Material16 *float64
	//`json:"material17"`
	Material17 *float64
	//`json:"material18"`
	Material18 *float64
	//`json:"material19"`
	Material19 *float64
	//`json:"material2"`
	Material2 *float64
	//`json:"material20"`
	Material20 *float64
	//`json:"material3"`
	Material3 *float64
	//`json:"material4"`
	Material4 *float64
	//`json:"material5"`
	Material5 *float64
	//`json:"material6"`
	Material6 *float64
	//`json:"material7"`
	Material7 *float64
	//`json:"material8"`
	Material8 *float64
	//`json:"material9"`
	Material9 *float64
	//`json:"name"`
	Name *string
	//`json:"nt"`
	Nt *string
	//`json:"number"`
	Number *string
	//`json:"tokenOfPile"`
	TokenOfPile *string

	//`json:"status"`
	Status *string
	//`json:"tag"`
	Tag *string
	//`json:"token"`
	Token *string
	//`json:"used"`
	Used *float64
	//`json:"x"`
	X *float64
	//`json:"y"`
	Y *float64
}

func (this *CreateObj) Convert2BoomGroup() boomGroupModel.BoomGroup {
	var boomGroup = boomGroupModel.BoomGroup{}

	boomGroup.Distance = this.Distance
	boomGroup.Geom = this.Geom
	boomGroup.High = this.High
	boomGroup.Low = this.Low
	boomGroup.Material1 = this.Material1
	boomGroup.Material10 = this.Material10
	boomGroup.Material11 = this.Material11
	boomGroup.Material12 = this.Material12
	boomGroup.Material13 = this.Material13
	boomGroup.Material14 = this.Material14
	boomGroup.Material15 = this.Material15
	boomGroup.Material16 = this.Material16
	boomGroup.Material17 = this.Material17
	boomGroup.Material18 = this.Material18
	boomGroup.Material19 = this.Material19
	boomGroup.Material2 = this.Material2
	boomGroup.Material20 = this.Material20
	boomGroup.Material3 = this.Material3
	boomGroup.Material4 = this.Material4
	boomGroup.Material5 = this.Material5
	boomGroup.Material6 = this.Material6
	boomGroup.Material7 = this.Material7
	boomGroup.Material8 = this.Material8
	boomGroup.Material9 = this.Material9
	boomGroup.Name = this.Name
	boomGroup.Nt = this.Nt
	boomGroup.Number = this.Number
	boomGroup.TokenOfPile = this.TokenOfPile

	boomGroup.Status = this.Status
	boomGroup.Tag = this.Tag
	boomGroup.Token = this.Token
	boomGroup.Used = this.Used
	boomGroup.X = this.X
	boomGroup.Y = this.Y
	return boomGroup
}

type UpdateObj struct {

	//`json:"distance"`
	Distance *float64
	//`json:"geom"`
	Geom *utils.Geometry `gorm:"type:geometry"`
	//`json:"high"`
	High *float64
	//`json:"low"`
	Low *float64
	//`json:"material1"`
	Material1 *float64
	//`json:"material10"`
	Material10 *float64
	//`json:"material11"`
	Material11 *float64
	//`json:"material12"`
	Material12 *float64
	//`json:"material13"`
	Material13 *float64
	//`json:"material14"`
	Material14 *float64
	//`json:"material15"`
	Material15 *float64
	//`json:"material16"`
	Material16 *float64
	//`json:"material17"`
	Material17 *float64
	//`json:"material18"`
	Material18 *float64
	//`json:"material19"`
	Material19 *float64
	//`json:"material2"`
	Material2 *float64
	//`json:"material20"`
	Material20 *float64
	//`json:"material3"`
	Material3 *float64
	//`json:"material4"`
	Material4 *float64
	//`json:"material5"`
	Material5 *float64
	//`json:"material6"`
	Material6 *float64
	//`json:"material7"`
	Material7 *float64
	//`json:"material8"`
	Material8 *float64
	//`json:"material9"`
	Material9 *float64
	//`json:"name"`
	Name *string
	//`json:"nt"`
	Nt *string
	//`json:"number"`
	Number *string
	//`json:"tokenOfPile"`
	TokenOfPile *string

	//`json:"status"`
	Status *string
	//`json:"tag"`
	Tag *string
	//`json:"token"`
	Token *string
	//`json:"used"`
	Used *float64
	//`json:"x"`
	X *float64
	//`json:"y"`
	Y *float64
}

func (this *UpdateObj) Convert2BoomGroup() boomGroupModel.BoomGroup {
	var boomGroup = boomGroupModel.BoomGroup{}

	boomGroup.Distance = this.Distance
	boomGroup.Geom = this.Geom
	boomGroup.High = this.High
	boomGroup.Low = this.Low
	boomGroup.Material1 = this.Material1
	boomGroup.Material10 = this.Material10
	boomGroup.Material11 = this.Material11
	boomGroup.Material12 = this.Material12
	boomGroup.Material13 = this.Material13
	boomGroup.Material14 = this.Material14
	boomGroup.Material15 = this.Material15
	boomGroup.Material16 = this.Material16
	boomGroup.Material17 = this.Material17
	boomGroup.Material18 = this.Material18
	boomGroup.Material19 = this.Material19
	boomGroup.Material2 = this.Material2
	boomGroup.Material20 = this.Material20
	boomGroup.Material3 = this.Material3
	boomGroup.Material4 = this.Material4
	boomGroup.Material5 = this.Material5
	boomGroup.Material6 = this.Material6
	boomGroup.Material7 = this.Material7
	boomGroup.Material8 = this.Material8
	boomGroup.Material9 = this.Material9
	boomGroup.Name = this.Name
	boomGroup.Nt = this.Nt
	boomGroup.Number = this.Number
	boomGroup.TokenOfPile = this.TokenOfPile

	boomGroup.Status = this.Status
	boomGroup.Tag = this.Tag
	boomGroup.Token = this.Token
	boomGroup.Used = this.Used
	boomGroup.X = this.X
	boomGroup.Y = this.Y
	return boomGroup
}
