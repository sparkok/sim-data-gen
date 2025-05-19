package mineProduct
import (
	"bytes"
	"sim_data_gen/requests/common"
	mineProductModel "sim_data_gen/models/mineProduct"
)


/**
* MineProduct请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"contentLimits"`
	ContentLimits *string  
	//`json:"matIndexes"`
	MatIndexes *string  
	//`json:"name"`
	Name *string 
	 
	//`json:"status"`
	Status *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.ContentLimits != nil {
      condition.WriteString("and (mine_product.content_limits like '%" + *this.ContentLimits + "%')")
    }
    if this.MatIndexes != nil {
      condition.WriteString("and (mine_product.mat_indexes like '%" + *this.MatIndexes + "%')")
    }
    if this.Name != nil {
      condition.WriteString("and (mine_product.name like '%" + *this.Name + "%')")
    }
    if this.Status != nil {
      condition.WriteString("and (mine_product.status = '" + *this.Status + "')")
    }
    if this.Token != nil {
      condition.WriteString("and (mine_product.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"contentLimits"`
	ContentLimits *string 
	//`json:"matIndexes"`
	MatIndexes *string 
	//`json:"name"`
	Name *string 
	
	//`json:"status"`
	Status *string 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2MineProduct() mineProductModel.MineProduct  {
	var mineProduct = mineProductModel.MineProduct{}

	mineProduct.ContentLimits = this.ContentLimits
	mineProduct.MatIndexes = this.MatIndexes
	mineProduct.Name = this.Name
	
	mineProduct.Status = this.Status
	mineProduct.Token = this.Token
	return mineProduct
}

type UpdateObj struct {

	//`json:"contentLimits"`
	ContentLimits *string 
	//`json:"matIndexes"`
	MatIndexes *string 
	//`json:"name"`
	Name *string 
	
	//`json:"status"`
	Status *string 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2MineProduct() mineProductModel.MineProduct  {
	var mineProduct = mineProductModel.MineProduct{}

	mineProduct.ContentLimits = this.ContentLimits
	mineProduct.MatIndexes = this.MatIndexes
	mineProduct.Name = this.Name
	
	mineProduct.Status = this.Status
	mineProduct.Token = this.Token
	return mineProduct
}
