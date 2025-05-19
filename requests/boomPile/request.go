package boomPile
import (
	"bytes"
	"sim_data_gen/requests/common"
	boomPileModel "sim_data_gen/models/boomPile"
)


/**
* BoomPile请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"bench"`
	Bench *string  
	//`json:"boomDate"`
	BoomDate *string  
	//`json:"costToGo"`
	CostToGo *string  
	//`json:"geom"`
	Geom *string  
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
	 
	//`json:"mineType"`
	MineType *string  
	//`json:"name"`
	Name *string  
	//`json:"nt"`
	Nt *string  
	//`json:"quantity"`
	Quantity *string 
	 
	//`json:"status"`
	Status *string  
	//`json:"tag"`
	Tag *string  
	//`json:"token"`
	Token *string  
	//`json:"used"`
	Used *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.Bench != nil {
      condition.WriteString("and (boom_pile.bench like '%" + *this.Bench + "%')")
    }
    if this.BoomDate != nil {
      condition.WriteString("and (boom_pile.boom_date like '%" + *this.BoomDate + "%')")
    }
    if this.CostToGo != nil {
     condition.WriteString("and (boom_pile.cost_to_go = " + *this.CostToGo + ")")
    }
    if this.Geom != nil {
      condition.WriteString("and (boom_pile.geom like '%" + *this.Geom + "%')")
    }
    if this.Material1 != nil {
     condition.WriteString("and (boom_pile.material1 = " + *this.Material1 + ")")
    }
    if this.Material10 != nil {
     condition.WriteString("and (boom_pile.material10 = " + *this.Material10 + ")")
    }
    if this.Material11 != nil {
     condition.WriteString("and (boom_pile.material11 = " + *this.Material11 + ")")
    }
    if this.Material12 != nil {
     condition.WriteString("and (boom_pile.material12 = " + *this.Material12 + ")")
    }
    if this.Material13 != nil {
     condition.WriteString("and (boom_pile.material13 = " + *this.Material13 + ")")
    }
    if this.Material14 != nil {
     condition.WriteString("and (boom_pile.material14 = " + *this.Material14 + ")")
    }
    if this.Material15 != nil {
     condition.WriteString("and (boom_pile.material15 = " + *this.Material15 + ")")
    }
    if this.Material16 != nil {
     condition.WriteString("and (boom_pile.material16 = " + *this.Material16 + ")")
    }
    if this.Material17 != nil {
     condition.WriteString("and (boom_pile.material17 = " + *this.Material17 + ")")
    }
    if this.Material18 != nil {
     condition.WriteString("and (boom_pile.material18 = " + *this.Material18 + ")")
    }
    if this.Material19 != nil {
     condition.WriteString("and (boom_pile.material19 = " + *this.Material19 + ")")
    }
    if this.Material2 != nil {
     condition.WriteString("and (boom_pile.material2 = " + *this.Material2 + ")")
    }
    if this.Material20 != nil {
     condition.WriteString("and (boom_pile.material20 = " + *this.Material20 + ")")
    }
    if this.Material3 != nil {
     condition.WriteString("and (boom_pile.material3 = " + *this.Material3 + ")")
    }
    if this.Material4 != nil {
     condition.WriteString("and (boom_pile.material4 = " + *this.Material4 + ")")
    }
    if this.Material5 != nil {
     condition.WriteString("and (boom_pile.material5 = " + *this.Material5 + ")")
    }
    if this.Material6 != nil {
     condition.WriteString("and (boom_pile.material6 = " + *this.Material6 + ")")
    }
    if this.Material7 != nil {
     condition.WriteString("and (boom_pile.material7 = " + *this.Material7 + ")")
    }
    if this.Material8 != nil {
     condition.WriteString("and (boom_pile.material8 = " + *this.Material8 + ")")
    }
    if this.Material9 != nil {
     condition.WriteString("and (boom_pile.material9 = " + *this.Material9 + ")")
    }
    if this.MineType != nil {
      condition.WriteString("and (boom_pile.mine_type = '" + *this.MineType + "')")
    }
    if this.Name != nil {
      condition.WriteString("and (boom_pile.name like '%" + *this.Name + "%')")
    }
    if this.Nt != nil {
      condition.WriteString("and (boom_pile.nt like '%" + *this.Nt + "%')")
    }
    if this.Quantity != nil {
     condition.WriteString("and (boom_pile.quantity = " + *this.Quantity + ")")
    }
    if this.Status != nil {
      condition.WriteString("and (boom_pile.status = '" + *this.Status + "')")
    }
    if this.Tag != nil {
      condition.WriteString("and (boom_pile.tag like '%" + *this.Tag + "%')")
    }
    if this.Token != nil {
      condition.WriteString("and (boom_pile.token like '%" + *this.Token + "%')")
    }
    if this.Used != nil {
     condition.WriteString("and (boom_pile.used = " + *this.Used + ")")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"bench"`
	Bench *string 
	//`json:"boomDate"`
	BoomDate *string 
	//`json:"costToGo"`
	CostToGo *int 
	//`json:"geom"`
	Geom *string 
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
	
	//`json:"mineType"`
	MineType *string 
	//`json:"name"`
	Name *string 
	//`json:"nt"`
	Nt *string 
	//`json:"quantity"`
	Quantity *float64 
	
	//`json:"status"`
	Status *string 
	//`json:"tag"`
	Tag *string 
	//`json:"token"`
	Token *string 
	//`json:"used"`
	Used *float64 
}
func (this *CreateObj) Convert2BoomPile() boomPileModel.BoomPile  {
	var boomPile = boomPileModel.BoomPile{}

	boomPile.Bench = this.Bench
	boomPile.BoomDate = this.BoomDate
	boomPile.CostToGo = this.CostToGo
	boomPile.Geom = this.Geom
	boomPile.Material1 = this.Material1
	boomPile.Material10 = this.Material10
	boomPile.Material11 = this.Material11
	boomPile.Material12 = this.Material12
	boomPile.Material13 = this.Material13
	boomPile.Material14 = this.Material14
	boomPile.Material15 = this.Material15
	boomPile.Material16 = this.Material16
	boomPile.Material17 = this.Material17
	boomPile.Material18 = this.Material18
	boomPile.Material19 = this.Material19
	boomPile.Material2 = this.Material2
	boomPile.Material20 = this.Material20
	boomPile.Material3 = this.Material3
	boomPile.Material4 = this.Material4
	boomPile.Material5 = this.Material5
	boomPile.Material6 = this.Material6
	boomPile.Material7 = this.Material7
	boomPile.Material8 = this.Material8
	boomPile.Material9 = this.Material9
	
	boomPile.MineType = this.MineType
	boomPile.Name = this.Name
	boomPile.Nt = this.Nt
	boomPile.Quantity = this.Quantity
	
	boomPile.Status = this.Status
	boomPile.Tag = this.Tag
	boomPile.Token = this.Token
	boomPile.Used = this.Used
	return boomPile
}

type UpdateObj struct {

	//`json:"bench"`
	Bench *string 
	//`json:"boomDate"`
	BoomDate *string 
	//`json:"costToGo"`
	CostToGo *int 
	//`json:"geom"`
	Geom *string 
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
	
	//`json:"mineType"`
	MineType *string 
	//`json:"name"`
	Name *string 
	//`json:"nt"`
	Nt *string 
	//`json:"quantity"`
	Quantity *float64 
	
	//`json:"status"`
	Status *string 
	//`json:"tag"`
	Tag *string 
	//`json:"token"`
	Token *string 
	//`json:"used"`
	Used *float64 
}
func (this *UpdateObj) Convert2BoomPile() boomPileModel.BoomPile  {
	var boomPile = boomPileModel.BoomPile{}

	boomPile.Bench = this.Bench
	boomPile.BoomDate = this.BoomDate
	boomPile.CostToGo = this.CostToGo
	boomPile.Geom = this.Geom
	boomPile.Material1 = this.Material1
	boomPile.Material10 = this.Material10
	boomPile.Material11 = this.Material11
	boomPile.Material12 = this.Material12
	boomPile.Material13 = this.Material13
	boomPile.Material14 = this.Material14
	boomPile.Material15 = this.Material15
	boomPile.Material16 = this.Material16
	boomPile.Material17 = this.Material17
	boomPile.Material18 = this.Material18
	boomPile.Material19 = this.Material19
	boomPile.Material2 = this.Material2
	boomPile.Material20 = this.Material20
	boomPile.Material3 = this.Material3
	boomPile.Material4 = this.Material4
	boomPile.Material5 = this.Material5
	boomPile.Material6 = this.Material6
	boomPile.Material7 = this.Material7
	boomPile.Material8 = this.Material8
	boomPile.Material9 = this.Material9
	
	boomPile.MineType = this.MineType
	boomPile.Name = this.Name
	boomPile.Nt = this.Nt
	boomPile.Quantity = this.Quantity
	
	boomPile.Status = this.Status
	boomPile.Tag = this.Tag
	boomPile.Token = this.Token
	boomPile.Used = this.Used
	return boomPile
}
