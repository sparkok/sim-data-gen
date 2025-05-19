package productAndBoomGroup
import (
	"bytes"
	"sim_data_gen/requests/common"
	productAndBoomGroupModel "sim_data_gen/models/productAndBoomGroup"
)


/**
* ProductAndBoomGroup请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {

	//`json:"tokenOfBoomGroup"`
	TokenOfBoomGroup *string 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string  
	//`json:"name"`
	Name *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.TokenOfBoomGroup != nil {
      condition.WriteString("and (product_and_boom_group.boom_group_id = '" + *this.TokenOfBoomGroup + "')")
    }
    if this.TokenOfMineProduct != nil {
      condition.WriteString("and (product_and_boom_group.mine_product_id = '" + *this.TokenOfMineProduct + "')")
    }
    if this.Name != nil {
      condition.WriteString("and (product_and_boom_group.name like '%" + *this.Name + "%')")
    }
    if this.Token != nil {
      condition.WriteString("and (product_and_boom_group.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"tokenOfBoomGroup"`
	TokenOfBoomGroup *string 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string 
	//`json:"name"`
	Name *string 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2ProductAndBoomGroup() productAndBoomGroupModel.ProductAndBoomGroup  {
	var productAndBoomGroup = productAndBoomGroupModel.ProductAndBoomGroup{}

	productAndBoomGroup.TokenOfBoomGroup = this.TokenOfBoomGroup
	productAndBoomGroup.TokenOfMineProduct = this.TokenOfMineProduct
	productAndBoomGroup.Name = this.Name
	productAndBoomGroup.Token = this.Token
	return productAndBoomGroup
}

type UpdateObj struct {

	//`json:"tokenOfBoomGroup"`
	TokenOfBoomGroup *string 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string 
	//`json:"name"`
	Name *string 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2ProductAndBoomGroup() productAndBoomGroupModel.ProductAndBoomGroup  {
	var productAndBoomGroup = productAndBoomGroupModel.ProductAndBoomGroup{}

	productAndBoomGroup.TokenOfBoomGroup = this.TokenOfBoomGroup
	productAndBoomGroup.TokenOfMineProduct = this.TokenOfMineProduct
	productAndBoomGroup.Name = this.Name
	productAndBoomGroup.Token = this.Token
	return productAndBoomGroup
}
