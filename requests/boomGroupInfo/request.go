package boomGroupInfo
import (
	"bytes"
	"sim_data_gen/requests/common"
	boomGroupInfoModel "sim_data_gen/models/boomGroupInfo"
)
import "time"


/**
* BoomGroupInfo请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {

	//`json:"tokenOfBoomGroup"`
	TokenOfBoomGroup *string  
	//`json:"createdAt"`
	CreatedAt *string  
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
	//`json:"probability1"`
	Probability1 *string  
	//`json:"probability10"`
	Probability10 *string  
	//`json:"probability11"`
	Probability11 *string  
	//`json:"probability12"`
	Probability12 *string  
	//`json:"probability13"`
	Probability13 *string  
	//`json:"probability14"`
	Probability14 *string  
	//`json:"probability15"`
	Probability15 *string  
	//`json:"probability16"`
	Probability16 *string  
	//`json:"probability17"`
	Probability17 *string  
	//`json:"probability18"`
	Probability18 *string  
	//`json:"probability19"`
	Probability19 *string  
	//`json:"probability2"`
	Probability2 *string  
	//`json:"probability20"`
	Probability20 *string  
	//`json:"probability3"`
	Probability3 *string  
	//`json:"probability4"`
	Probability4 *string  
	//`json:"probability5"`
	Probability5 *string  
	//`json:"probability6"`
	Probability6 *string  
	//`json:"probability7"`
	Probability7 *string  
	//`json:"probability8"`
	Probability8 *string  
	//`json:"probability9"`
	Probability9 *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.TokenOfBoomGroup != nil {
      condition.WriteString("and (boom_group_info.boom_group_id = '" + *this.TokenOfBoomGroup + "')")
    }
    if this.CreatedAt != nil {
      condition.WriteString("and (boom_group_info.created_at = '" + *this.CreatedAt + "')")
    }
    if this.Material1 != nil {
     condition.WriteString("and (boom_group_info.material1 = " + *this.Material1 + ")")
    }
    if this.Material10 != nil {
     condition.WriteString("and (boom_group_info.material10 = " + *this.Material10 + ")")
    }
    if this.Material11 != nil {
     condition.WriteString("and (boom_group_info.material11 = " + *this.Material11 + ")")
    }
    if this.Material12 != nil {
     condition.WriteString("and (boom_group_info.material12 = " + *this.Material12 + ")")
    }
    if this.Material13 != nil {
     condition.WriteString("and (boom_group_info.material13 = " + *this.Material13 + ")")
    }
    if this.Material14 != nil {
     condition.WriteString("and (boom_group_info.material14 = " + *this.Material14 + ")")
    }
    if this.Material15 != nil {
     condition.WriteString("and (boom_group_info.material15 = " + *this.Material15 + ")")
    }
    if this.Material16 != nil {
     condition.WriteString("and (boom_group_info.material16 = " + *this.Material16 + ")")
    }
    if this.Material17 != nil {
     condition.WriteString("and (boom_group_info.material17 = " + *this.Material17 + ")")
    }
    if this.Material18 != nil {
     condition.WriteString("and (boom_group_info.material18 = " + *this.Material18 + ")")
    }
    if this.Material19 != nil {
     condition.WriteString("and (boom_group_info.material19 = " + *this.Material19 + ")")
    }
    if this.Material2 != nil {
     condition.WriteString("and (boom_group_info.material2 = " + *this.Material2 + ")")
    }
    if this.Material20 != nil {
     condition.WriteString("and (boom_group_info.material20 = " + *this.Material20 + ")")
    }
    if this.Material3 != nil {
     condition.WriteString("and (boom_group_info.material3 = " + *this.Material3 + ")")
    }
    if this.Material4 != nil {
     condition.WriteString("and (boom_group_info.material4 = " + *this.Material4 + ")")
    }
    if this.Material5 != nil {
     condition.WriteString("and (boom_group_info.material5 = " + *this.Material5 + ")")
    }
    if this.Material6 != nil {
     condition.WriteString("and (boom_group_info.material6 = " + *this.Material6 + ")")
    }
    if this.Material7 != nil {
     condition.WriteString("and (boom_group_info.material7 = " + *this.Material7 + ")")
    }
    if this.Material8 != nil {
     condition.WriteString("and (boom_group_info.material8 = " + *this.Material8 + ")")
    }
    if this.Material9 != nil {
     condition.WriteString("and (boom_group_info.material9 = " + *this.Material9 + ")")
    }
    if this.Name != nil {
      condition.WriteString("and (boom_group_info.name like '%" + *this.Name + "%')")
    }
    if this.Probability1 != nil {
     condition.WriteString("and (boom_group_info.probability1 = " + *this.Probability1 + ")")
    }
    if this.Probability10 != nil {
     condition.WriteString("and (boom_group_info.probability10 = " + *this.Probability10 + ")")
    }
    if this.Probability11 != nil {
     condition.WriteString("and (boom_group_info.probability11 = " + *this.Probability11 + ")")
    }
    if this.Probability12 != nil {
     condition.WriteString("and (boom_group_info.probability12 = " + *this.Probability12 + ")")
    }
    if this.Probability13 != nil {
     condition.WriteString("and (boom_group_info.probability13 = " + *this.Probability13 + ")")
    }
    if this.Probability14 != nil {
     condition.WriteString("and (boom_group_info.probability14 = " + *this.Probability14 + ")")
    }
    if this.Probability15 != nil {
     condition.WriteString("and (boom_group_info.probability15 = " + *this.Probability15 + ")")
    }
    if this.Probability16 != nil {
     condition.WriteString("and (boom_group_info.probability16 = " + *this.Probability16 + ")")
    }
    if this.Probability17 != nil {
     condition.WriteString("and (boom_group_info.probability17 = " + *this.Probability17 + ")")
    }
    if this.Probability18 != nil {
     condition.WriteString("and (boom_group_info.probability18 = " + *this.Probability18 + ")")
    }
    if this.Probability19 != nil {
     condition.WriteString("and (boom_group_info.probability19 = " + *this.Probability19 + ")")
    }
    if this.Probability2 != nil {
     condition.WriteString("and (boom_group_info.probability2 = " + *this.Probability2 + ")")
    }
    if this.Probability20 != nil {
     condition.WriteString("and (boom_group_info.probability20 = " + *this.Probability20 + ")")
    }
    if this.Probability3 != nil {
     condition.WriteString("and (boom_group_info.probability3 = " + *this.Probability3 + ")")
    }
    if this.Probability4 != nil {
     condition.WriteString("and (boom_group_info.probability4 = " + *this.Probability4 + ")")
    }
    if this.Probability5 != nil {
     condition.WriteString("and (boom_group_info.probability5 = " + *this.Probability5 + ")")
    }
    if this.Probability6 != nil {
     condition.WriteString("and (boom_group_info.probability6 = " + *this.Probability6 + ")")
    }
    if this.Probability7 != nil {
     condition.WriteString("and (boom_group_info.probability7 = " + *this.Probability7 + ")")
    }
    if this.Probability8 != nil {
     condition.WriteString("and (boom_group_info.probability8 = " + *this.Probability8 + ")")
    }
    if this.Probability9 != nil {
     condition.WriteString("and (boom_group_info.probability9 = " + *this.Probability9 + ")")
    }
    if this.Token != nil {
      condition.WriteString("and (boom_group_info.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"tokenOfBoomGroup"`
	TokenOfBoomGroup *string 
	//`json:"createdAt"`
	CreatedAt *time.Time 
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
	//`json:"probability1"`
	Probability1 *float64 
	//`json:"probability10"`
	Probability10 *float64 
	//`json:"probability11"`
	Probability11 *float64 
	//`json:"probability12"`
	Probability12 *float64 
	//`json:"probability13"`
	Probability13 *float64 
	//`json:"probability14"`
	Probability14 *float64 
	//`json:"probability15"`
	Probability15 *float64 
	//`json:"probability16"`
	Probability16 *float64 
	//`json:"probability17"`
	Probability17 *float64 
	//`json:"probability18"`
	Probability18 *float64 
	//`json:"probability19"`
	Probability19 *float64 
	//`json:"probability2"`
	Probability2 *float64 
	//`json:"probability20"`
	Probability20 *float64 
	//`json:"probability3"`
	Probability3 *float64 
	//`json:"probability4"`
	Probability4 *float64 
	//`json:"probability5"`
	Probability5 *float64 
	//`json:"probability6"`
	Probability6 *float64 
	//`json:"probability7"`
	Probability7 *float64 
	//`json:"probability8"`
	Probability8 *float64 
	//`json:"probability9"`
	Probability9 *float64 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2BoomGroupInfo() boomGroupInfoModel.BoomGroupInfo  {
	var boomGroupInfo = boomGroupInfoModel.BoomGroupInfo{}

	boomGroupInfo.TokenOfBoomGroup = this.TokenOfBoomGroup
	boomGroupInfo.CreatedAt = this.CreatedAt
	boomGroupInfo.Material1 = this.Material1
	boomGroupInfo.Material10 = this.Material10
	boomGroupInfo.Material11 = this.Material11
	boomGroupInfo.Material12 = this.Material12
	boomGroupInfo.Material13 = this.Material13
	boomGroupInfo.Material14 = this.Material14
	boomGroupInfo.Material15 = this.Material15
	boomGroupInfo.Material16 = this.Material16
	boomGroupInfo.Material17 = this.Material17
	boomGroupInfo.Material18 = this.Material18
	boomGroupInfo.Material19 = this.Material19
	boomGroupInfo.Material2 = this.Material2
	boomGroupInfo.Material20 = this.Material20
	boomGroupInfo.Material3 = this.Material3
	boomGroupInfo.Material4 = this.Material4
	boomGroupInfo.Material5 = this.Material5
	boomGroupInfo.Material6 = this.Material6
	boomGroupInfo.Material7 = this.Material7
	boomGroupInfo.Material8 = this.Material8
	boomGroupInfo.Material9 = this.Material9
	boomGroupInfo.Name = this.Name
	boomGroupInfo.Probability1 = this.Probability1
	boomGroupInfo.Probability10 = this.Probability10
	boomGroupInfo.Probability11 = this.Probability11
	boomGroupInfo.Probability12 = this.Probability12
	boomGroupInfo.Probability13 = this.Probability13
	boomGroupInfo.Probability14 = this.Probability14
	boomGroupInfo.Probability15 = this.Probability15
	boomGroupInfo.Probability16 = this.Probability16
	boomGroupInfo.Probability17 = this.Probability17
	boomGroupInfo.Probability18 = this.Probability18
	boomGroupInfo.Probability19 = this.Probability19
	boomGroupInfo.Probability2 = this.Probability2
	boomGroupInfo.Probability20 = this.Probability20
	boomGroupInfo.Probability3 = this.Probability3
	boomGroupInfo.Probability4 = this.Probability4
	boomGroupInfo.Probability5 = this.Probability5
	boomGroupInfo.Probability6 = this.Probability6
	boomGroupInfo.Probability7 = this.Probability7
	boomGroupInfo.Probability8 = this.Probability8
	boomGroupInfo.Probability9 = this.Probability9
	boomGroupInfo.Token = this.Token
	return boomGroupInfo
}

type UpdateObj struct {

	//`json:"tokenOfBoomGroup"`
	TokenOfBoomGroup *string 
	//`json:"createdAt"`
	CreatedAt *time.Time 
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
	//`json:"probability1"`
	Probability1 *float64 
	//`json:"probability10"`
	Probability10 *float64 
	//`json:"probability11"`
	Probability11 *float64 
	//`json:"probability12"`
	Probability12 *float64 
	//`json:"probability13"`
	Probability13 *float64 
	//`json:"probability14"`
	Probability14 *float64 
	//`json:"probability15"`
	Probability15 *float64 
	//`json:"probability16"`
	Probability16 *float64 
	//`json:"probability17"`
	Probability17 *float64 
	//`json:"probability18"`
	Probability18 *float64 
	//`json:"probability19"`
	Probability19 *float64 
	//`json:"probability2"`
	Probability2 *float64 
	//`json:"probability20"`
	Probability20 *float64 
	//`json:"probability3"`
	Probability3 *float64 
	//`json:"probability4"`
	Probability4 *float64 
	//`json:"probability5"`
	Probability5 *float64 
	//`json:"probability6"`
	Probability6 *float64 
	//`json:"probability7"`
	Probability7 *float64 
	//`json:"probability8"`
	Probability8 *float64 
	//`json:"probability9"`
	Probability9 *float64 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2BoomGroupInfo() boomGroupInfoModel.BoomGroupInfo  {
	var boomGroupInfo = boomGroupInfoModel.BoomGroupInfo{}

	boomGroupInfo.TokenOfBoomGroup = this.TokenOfBoomGroup
	boomGroupInfo.CreatedAt = this.CreatedAt
	boomGroupInfo.Material1 = this.Material1
	boomGroupInfo.Material10 = this.Material10
	boomGroupInfo.Material11 = this.Material11
	boomGroupInfo.Material12 = this.Material12
	boomGroupInfo.Material13 = this.Material13
	boomGroupInfo.Material14 = this.Material14
	boomGroupInfo.Material15 = this.Material15
	boomGroupInfo.Material16 = this.Material16
	boomGroupInfo.Material17 = this.Material17
	boomGroupInfo.Material18 = this.Material18
	boomGroupInfo.Material19 = this.Material19
	boomGroupInfo.Material2 = this.Material2
	boomGroupInfo.Material20 = this.Material20
	boomGroupInfo.Material3 = this.Material3
	boomGroupInfo.Material4 = this.Material4
	boomGroupInfo.Material5 = this.Material5
	boomGroupInfo.Material6 = this.Material6
	boomGroupInfo.Material7 = this.Material7
	boomGroupInfo.Material8 = this.Material8
	boomGroupInfo.Material9 = this.Material9
	boomGroupInfo.Name = this.Name
	boomGroupInfo.Probability1 = this.Probability1
	boomGroupInfo.Probability10 = this.Probability10
	boomGroupInfo.Probability11 = this.Probability11
	boomGroupInfo.Probability12 = this.Probability12
	boomGroupInfo.Probability13 = this.Probability13
	boomGroupInfo.Probability14 = this.Probability14
	boomGroupInfo.Probability15 = this.Probability15
	boomGroupInfo.Probability16 = this.Probability16
	boomGroupInfo.Probability17 = this.Probability17
	boomGroupInfo.Probability18 = this.Probability18
	boomGroupInfo.Probability19 = this.Probability19
	boomGroupInfo.Probability2 = this.Probability2
	boomGroupInfo.Probability20 = this.Probability20
	boomGroupInfo.Probability3 = this.Probability3
	boomGroupInfo.Probability4 = this.Probability4
	boomGroupInfo.Probability5 = this.Probability5
	boomGroupInfo.Probability6 = this.Probability6
	boomGroupInfo.Probability7 = this.Probability7
	boomGroupInfo.Probability8 = this.Probability8
	boomGroupInfo.Probability9 = this.Probability9
	boomGroupInfo.Token = this.Token
	return boomGroupInfo
}
