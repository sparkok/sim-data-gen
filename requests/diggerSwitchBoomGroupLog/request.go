package diggerSwitchBoomGroupLog
import (
	"bytes"
	"sim_data_gen/requests/common"
	diggerSwitchBoomGroupLogModel "sim_data_gen/models/diggerSwitchBoomGroupLog"
)


/**
* DiggerSwitchBoomGroupLog请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"applyUtc"`
	ApplyUtc *string  
	//`json:"boomGroupId"`
	BoomGroupId *string  
	//`json:"dateFlag"`
	DateFlag *string  
	//`json:"diggerId"`
	DiggerId *string  
	//`json:"name"`
	Name *string 
	 
	//`json:"status"`
	Status *string  
	//`json:"submitUtc"`
	SubmitUtc *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.ApplyUtc != nil {
     condition.WriteString("and (digger_switch_boom_group_log.apply_utc = " + *this.ApplyUtc + ")")
    }
    if this.BoomGroupId != nil {
      condition.WriteString("and (digger_switch_boom_group_log.boom_group_id like '%" + *this.BoomGroupId + "%')")
    }
    if this.DateFlag != nil {
      condition.WriteString("and (digger_switch_boom_group_log.date_flag like '%" + *this.DateFlag + "%')")
    }
    if this.DiggerId != nil {
      condition.WriteString("and (digger_switch_boom_group_log.digger_id like '%" + *this.DiggerId + "%')")
    }
    if this.Name != nil {
      condition.WriteString("and (digger_switch_boom_group_log.name like '%" + *this.Name + "%')")
    }
    if this.Status != nil {
     condition.WriteString("and (digger_switch_boom_group_log.status = " + *this.Status + ")")
    }
    if this.SubmitUtc != nil {
     condition.WriteString("and (digger_switch_boom_group_log.submit_utc = " + *this.SubmitUtc + ")")
    }
    if this.Token != nil {
      condition.WriteString("and (digger_switch_boom_group_log.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"applyUtc"`
	ApplyUtc *int 
	//`json:"boomGroupId"`
	BoomGroupId *string 
	//`json:"dateFlag"`
	DateFlag *string 
	//`json:"diggerId"`
	DiggerId *string 
	//`json:"name"`
	Name *string 
	
	//`json:"status"`
	Status *int 
	//`json:"submitUtc"`
	SubmitUtc *int 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2DiggerSwitchBoomGroupLog() diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog  {
	var diggerSwitchBoomGroupLog = diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{}

	diggerSwitchBoomGroupLog.ApplyUtc = this.ApplyUtc
	diggerSwitchBoomGroupLog.BoomGroupId = this.BoomGroupId
	diggerSwitchBoomGroupLog.DateFlag = this.DateFlag
	diggerSwitchBoomGroupLog.DiggerId = this.DiggerId
	diggerSwitchBoomGroupLog.Name = this.Name
	
	diggerSwitchBoomGroupLog.Status = this.Status
	diggerSwitchBoomGroupLog.SubmitUtc = this.SubmitUtc
	diggerSwitchBoomGroupLog.Token = this.Token
	return diggerSwitchBoomGroupLog
}

type UpdateObj struct {

	//`json:"applyUtc"`
	ApplyUtc *int 
	//`json:"boomGroupId"`
	BoomGroupId *string 
	//`json:"dateFlag"`
	DateFlag *string 
	//`json:"diggerId"`
	DiggerId *string 
	//`json:"name"`
	Name *string 
	
	//`json:"status"`
	Status *int 
	//`json:"submitUtc"`
	SubmitUtc *int 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2DiggerSwitchBoomGroupLog() diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog  {
	var diggerSwitchBoomGroupLog = diggerSwitchBoomGroupLogModel.DiggerSwitchBoomGroupLog{}

	diggerSwitchBoomGroupLog.ApplyUtc = this.ApplyUtc
	diggerSwitchBoomGroupLog.BoomGroupId = this.BoomGroupId
	diggerSwitchBoomGroupLog.DateFlag = this.DateFlag
	diggerSwitchBoomGroupLog.DiggerId = this.DiggerId
	diggerSwitchBoomGroupLog.Name = this.Name
	
	diggerSwitchBoomGroupLog.Status = this.Status
	diggerSwitchBoomGroupLog.SubmitUtc = this.SubmitUtc
	diggerSwitchBoomGroupLog.Token = this.Token
	return diggerSwitchBoomGroupLog
}
