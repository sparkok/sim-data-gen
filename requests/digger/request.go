package digger
import (
	"bytes"
	"sim_data_gen/requests/common"
	diggerModel "sim_data_gen/models/digger"
)


/**
* Digger请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"attribs"`
	Attribs *string  
	//`json:"name"`
	Name *string  
	//`json:"produce"`
	Produce *string  
	//`json:"speed"`
	Speed *string 
	 
	//`json:"status"`
	Status *string  
	//`json:"token"`
	Token *string  
	//`json:"utc"`
	Utc *string  
	//`json:"x"`
	X *string  
	//`json:"y"`
	Y *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.Attribs != nil {
      condition.WriteString("and (digger.attribs like '%" + *this.Attribs + "%')")
    }
    if this.Name != nil {
      condition.WriteString("and (digger.name like '%" + *this.Name + "%')")
    }
    if this.Produce != nil {
     condition.WriteString("and (digger.produce = " + *this.Produce + ")")
    }
    if this.Speed != nil {
     condition.WriteString("and (digger.speed = " + *this.Speed + ")")
    }
    if this.Status != nil {
     condition.WriteString("and (digger.status = " + *this.Status + ")")
    }
    if this.Token != nil {
      condition.WriteString("and (digger.token like '%" + *this.Token + "%')")
    }
    if this.Utc != nil {
     condition.WriteString("and (digger.utc = " + *this.Utc + ")")
    }
    if this.X != nil {
     condition.WriteString("and (digger.x = " + *this.X + ")")
    }
    if this.Y != nil {
     condition.WriteString("and (digger.y = " + *this.Y + ")")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"attribs"`
	Attribs *string 
	//`json:"name"`
	Name *string 
	//`json:"produce"`
	Produce *float64 
	//`json:"speed"`
	Speed *float64 
	
	//`json:"status"`
	Status *int 
	//`json:"token"`
	Token *string 
	//`json:"utc"`
	Utc *int 
	//`json:"x"`
	X *float64 
	//`json:"y"`
	Y *float64 
}
func (this *CreateObj) Convert2Digger() diggerModel.Digger  {
	var digger = diggerModel.Digger{}

	digger.Attribs = this.Attribs
	digger.Name = this.Name
	digger.Produce = this.Produce
	digger.Speed = this.Speed
	
	digger.Status = this.Status
	digger.Token = this.Token
	digger.Utc = this.Utc
	digger.X = this.X
	digger.Y = this.Y
	return digger
}

type UpdateObj struct {

	//`json:"attribs"`
	Attribs *string 
	//`json:"name"`
	Name *string 
	//`json:"produce"`
	Produce *float64 
	//`json:"speed"`
	Speed *float64 
	
	//`json:"status"`
	Status *int 
	//`json:"token"`
	Token *string 
	//`json:"utc"`
	Utc *int 
	//`json:"x"`
	X *float64 
	//`json:"y"`
	Y *float64 
}
func (this *UpdateObj) Convert2Digger() diggerModel.Digger  {
	var digger = diggerModel.Digger{}

	digger.Attribs = this.Attribs
	digger.Name = this.Name
	digger.Produce = this.Produce
	digger.Speed = this.Speed
	
	digger.Status = this.Status
	digger.Token = this.Token
	digger.Utc = this.Utc
	digger.X = this.X
	digger.Y = this.Y
	return digger
}
