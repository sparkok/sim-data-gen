package mineChangeStatus
import (
	"bytes"
	"sim_data_gen/requests/common"
	mineChangeStatusModel "sim_data_gen/models/mineChangeStatus"
)
import "time"


/**
* MineChangeStatus请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"changing"`
	Changing *string  
	//`json:"createdAt"`
	CreatedAt *string  
	//`json:"dataTimeOfBridge"`
	DataTimeOfBridge *string  
	//`json:"dataTimeOfNeutron"`
	DataTimeOfNeutron *string  
	//`json:"dataType"`
	DataType *string  
	//`json:"dateFlag"`
	DateFlag *string  
	//`json:"name"`
	Name *string  
	//`json:"productName"`
	ProductName *string  
	//`json:"token"`
	Token *string  
	//`json:"totalMassOfBridge"`
	TotalMassOfBridge *string  
	//`json:"totalMassOfNeutron"`
	TotalMassOfNeutron *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.Changing != nil {
      condition.WriteString("and (mine_change_status.changing like '%" + *this.Changing + "%')")
    }
    if this.CreatedAt != nil {
      condition.WriteString("and (mine_change_status.created_at = '" + *this.CreatedAt + "')")
    }
    if this.DataTimeOfBridge != nil {
      condition.WriteString("and (mine_change_status.data_time_of_bridge = '" + *this.DataTimeOfBridge + "')")
    }
    if this.DataTimeOfNeutron != nil {
      condition.WriteString("and (mine_change_status.data_time_of_neutron = '" + *this.DataTimeOfNeutron + "')")
    }
    if this.DataType != nil {
     condition.WriteString("and (mine_change_status.data_type = " + *this.DataType + ")")
    }
    if this.DateFlag != nil {
      condition.WriteString("and (mine_change_status.date_flag like '%" + *this.DateFlag + "%')")
    }
    if this.Name != nil {
      condition.WriteString("and (mine_change_status.name like '%" + *this.Name + "%')")
    }
    if this.ProductName != nil {
      condition.WriteString("and (mine_change_status.product_name like '%" + *this.ProductName + "%')")
    }
    if this.Token != nil {
      condition.WriteString("and (mine_change_status.token like '%" + *this.Token + "%')")
    }
    if this.TotalMassOfBridge != nil {
     condition.WriteString("and (mine_change_status.total_mass_of_bridge = " + *this.TotalMassOfBridge + ")")
    }
    if this.TotalMassOfNeutron != nil {
     condition.WriteString("and (mine_change_status.total_mass_of_neutron = " + *this.TotalMassOfNeutron + ")")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"changing"`
	Changing *string 
	//`json:"createdAt"`
	CreatedAt *time.Time 
	//`json:"dataTimeOfBridge"`
	DataTimeOfBridge *time.Time 
	//`json:"dataTimeOfNeutron"`
	DataTimeOfNeutron *time.Time 
	//`json:"dataType"`
	DataType *int 
	//`json:"dateFlag"`
	DateFlag *string 
	//`json:"name"`
	Name *string 
	//`json:"productName"`
	ProductName *string 
	//`json:"token"`
	Token *string 
	//`json:"totalMassOfBridge"`
	TotalMassOfBridge *float64 
	//`json:"totalMassOfNeutron"`
	TotalMassOfNeutron *float64 
}
func (this *CreateObj) Convert2MineChangeStatus() mineChangeStatusModel.MineChangeStatus  {
	var mineChangeStatus = mineChangeStatusModel.MineChangeStatus{}

	mineChangeStatus.Changing = this.Changing
	mineChangeStatus.CreatedAt = this.CreatedAt
	mineChangeStatus.DataTimeOfBridge = this.DataTimeOfBridge
	mineChangeStatus.DataTimeOfNeutron = this.DataTimeOfNeutron
	mineChangeStatus.DataType = this.DataType
	mineChangeStatus.DateFlag = this.DateFlag
	mineChangeStatus.Name = this.Name
	mineChangeStatus.ProductName = this.ProductName
	mineChangeStatus.Token = this.Token
	mineChangeStatus.TotalMassOfBridge = this.TotalMassOfBridge
	mineChangeStatus.TotalMassOfNeutron = this.TotalMassOfNeutron
	return mineChangeStatus
}

type UpdateObj struct {

	//`json:"changing"`
	Changing *string 
	//`json:"createdAt"`
	CreatedAt *time.Time 
	//`json:"dataTimeOfBridge"`
	DataTimeOfBridge *time.Time 
	//`json:"dataTimeOfNeutron"`
	DataTimeOfNeutron *time.Time 
	//`json:"dataType"`
	DataType *int 
	//`json:"dateFlag"`
	DateFlag *string 
	//`json:"name"`
	Name *string 
	//`json:"productName"`
	ProductName *string 
	//`json:"token"`
	Token *string 
	//`json:"totalMassOfBridge"`
	TotalMassOfBridge *float64 
	//`json:"totalMassOfNeutron"`
	TotalMassOfNeutron *float64 
}
func (this *UpdateObj) Convert2MineChangeStatus() mineChangeStatusModel.MineChangeStatus  {
	var mineChangeStatus = mineChangeStatusModel.MineChangeStatus{}

	mineChangeStatus.Changing = this.Changing
	mineChangeStatus.CreatedAt = this.CreatedAt
	mineChangeStatus.DataTimeOfBridge = this.DataTimeOfBridge
	mineChangeStatus.DataTimeOfNeutron = this.DataTimeOfNeutron
	mineChangeStatus.DataType = this.DataType
	mineChangeStatus.DateFlag = this.DateFlag
	mineChangeStatus.Name = this.Name
	mineChangeStatus.ProductName = this.ProductName
	mineChangeStatus.Token = this.Token
	mineChangeStatus.TotalMassOfBridge = this.TotalMassOfBridge
	mineChangeStatus.TotalMassOfNeutron = this.TotalMassOfNeutron
	return mineChangeStatus
}
