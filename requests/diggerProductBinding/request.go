package diggerProductBinding
import (
	"bytes"
	"sim_data_gen/requests/common"
	diggerProductBindingModel "sim_data_gen/models/diggerProductBinding"
)


/**
* DiggerProductBinding请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {

	//`json:"tokenOfDigger"`
	TokenOfDigger *string 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string  
	//`json:"name"`
	Name *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.TokenOfDigger != nil {
      condition.WriteString("and (digger_product_binding.digger_id = '" + *this.TokenOfDigger + "')")
    }
    if this.TokenOfMineProduct != nil {
      condition.WriteString("and (digger_product_binding.mine_product_id = '" + *this.TokenOfMineProduct + "')")
    }
    if this.Name != nil {
      condition.WriteString("and (digger_product_binding.name like '%" + *this.Name + "%')")
    }
    if this.Token != nil {
      condition.WriteString("and (digger_product_binding.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"tokenOfDigger"`
	TokenOfDigger *string 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string 
	//`json:"name"`
	Name *string 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2DiggerProductBinding() diggerProductBindingModel.DiggerProductBinding  {
	var diggerProductBinding = diggerProductBindingModel.DiggerProductBinding{}

	diggerProductBinding.TokenOfDigger = this.TokenOfDigger
	diggerProductBinding.TokenOfMineProduct = this.TokenOfMineProduct
	diggerProductBinding.Name = this.Name
	diggerProductBinding.Token = this.Token
	return diggerProductBinding
}

type UpdateObj struct {

	//`json:"tokenOfDigger"`
	TokenOfDigger *string 
	//`json:"tokenOfMineProduct"`
	TokenOfMineProduct *string 
	//`json:"name"`
	Name *string 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2DiggerProductBinding() diggerProductBindingModel.DiggerProductBinding  {
	var diggerProductBinding = diggerProductBindingModel.DiggerProductBinding{}

	diggerProductBinding.TokenOfDigger = this.TokenOfDigger
	diggerProductBinding.TokenOfMineProduct = this.TokenOfMineProduct
	diggerProductBinding.Name = this.Name
	diggerProductBinding.Token = this.Token
	return diggerProductBinding
}
