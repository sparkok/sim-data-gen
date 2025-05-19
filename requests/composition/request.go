package composition
import (
	"bytes"
	"sim_data_gen/requests/common"
	compositionModel "sim_data_gen/models/composition"
)


/**
* Composition请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"name"`
	Name *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.Name != nil {
      condition.WriteString("and (composition.name like '%" + *this.Name + "%')")
    }
    if this.Token != nil {
      condition.WriteString("and (composition.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"name"`
	Name *string 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2Composition() compositionModel.Composition  {
	var composition = compositionModel.Composition{}

	composition.Name = this.Name
	composition.Token = this.Token
	return composition
}

type UpdateObj struct {

	//`json:"name"`
	Name *string 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2Composition() compositionModel.Composition  {
	var composition = compositionModel.Composition{}

	composition.Name = this.Name
	composition.Token = this.Token
	return composition
}
