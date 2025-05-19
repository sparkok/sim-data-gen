package contentPercent
import (
	"bytes"
	"sim_data_gen/requests/common"
	contentPercentModel "sim_data_gen/models/contentPercent"
)


/**
* ContentPercent请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {

	//`json:"tokenOfComposition"`
	TokenOfComposition *string  
	//`json:"maxValue"`
	MaxValue *string  
	//`json:"minValue"`
	MinValue *string 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string  
	//`json:"name"`
	Name *string  
	//`json:"num"`
	Num *string 
	 
	//`json:"status"`
	Status *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.TokenOfComposition != nil {
      condition.WriteString("and (content_percent.composition_id = '" + *this.TokenOfComposition + "')")
    }
    if this.MaxValue != nil {
     condition.WriteString("and (content_percent.max_value = " + *this.MaxValue + ")")
    }
    if this.MinValue != nil {
     condition.WriteString("and (content_percent.min_value = " + *this.MinValue + ")")
    }
    if this.TokenOfMineProduct != nil {
      condition.WriteString("and (content_percent.mine_product_id = '" + *this.TokenOfMineProduct + "')")
    }
    if this.Name != nil {
      condition.WriteString("and (content_percent.name like '%" + *this.Name + "%')")
    }
    if this.Num != nil {
     condition.WriteString("and (content_percent.num = " + *this.Num + ")")
    }
    if this.Status != nil {
     condition.WriteString("and (content_percent.status = " + *this.Status + ")")
    }
    if this.Token != nil {
      condition.WriteString("and (content_percent.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"tokenOfComposition"`
	TokenOfComposition *string 
	//`json:"maxValue"`
	MaxValue *float64 
	//`json:"minValue"`
	MinValue *float64 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string 
	//`json:"name"`
	Name *string 
	//`json:"num"`
	Num *int 
	
	//`json:"status"`
	Status *int 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2ContentPercent() contentPercentModel.ContentPercent  {
	var contentPercent = contentPercentModel.ContentPercent{}

	contentPercent.TokenOfComposition = this.TokenOfComposition
	contentPercent.MaxValue = this.MaxValue
	contentPercent.MinValue = this.MinValue
	contentPercent.TokenOfMineProduct = this.TokenOfMineProduct
	contentPercent.Name = this.Name
	contentPercent.Num = this.Num
	
	contentPercent.Status = this.Status
	contentPercent.Token = this.Token
	return contentPercent
}

type UpdateObj struct {

	//`json:"tokenOfComposition"`
	TokenOfComposition *string 
	//`json:"maxValue"`
	MaxValue *float64 
	//`json:"minValue"`
	MinValue *float64 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string 
	//`json:"name"`
	Name *string 
	//`json:"num"`
	Num *int 
	
	//`json:"status"`
	Status *int 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2ContentPercent() contentPercentModel.ContentPercent  {
	var contentPercent = contentPercentModel.ContentPercent{}

	contentPercent.TokenOfComposition = this.TokenOfComposition
	contentPercent.MaxValue = this.MaxValue
	contentPercent.MinValue = this.MinValue
	contentPercent.TokenOfMineProduct = this.TokenOfMineProduct
	contentPercent.Name = this.Name
	contentPercent.Num = this.Num
	
	contentPercent.Status = this.Status
	contentPercent.Token = this.Token
	return contentPercent
}
