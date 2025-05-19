package lorryNearbyTargetSpan
import (
	"bytes"
	"sim_data_gen/requests/common"
	lorryNearbyTargetSpanModel "sim_data_gen/models/lorryNearbyTargetSpan"
)


/**
* LorryNearbyTargetSpan请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"beginUtc"`
	BeginUtc *string  
	//`json:"dateFlag"`
	DateFlag *string  
	//`json:"endUtc"`
	EndUtc *string  
	//`json:"lorryId"`
	LorryId *string  
	//`json:"name"`
	Name *string  
	//`json:"nearbyObj"`
	NearbyObj *string 
	 
	//`json:"objType"`
	ObjType *string  
	//`json:"productName"`
	ProductName *string  
	//`json:"token"`
	Token *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.BeginUtc != nil {
     condition.WriteString("and (lorry_nearby_target_span.begin_utc = " + *this.BeginUtc + ")")
    }
    if this.DateFlag != nil {
      condition.WriteString("and (lorry_nearby_target_span.date_flag like '%" + *this.DateFlag + "%')")
    }
    if this.EndUtc != nil {
     condition.WriteString("and (lorry_nearby_target_span.end_utc = " + *this.EndUtc + ")")
    }
    if this.LorryId != nil {
      condition.WriteString("and (lorry_nearby_target_span.lorry_id like '%" + *this.LorryId + "%')")
    }
    if this.Name != nil {
      condition.WriteString("and (lorry_nearby_target_span.name like '%" + *this.Name + "%')")
    }
    if this.NearbyObj != nil {
      condition.WriteString("and (lorry_nearby_target_span.nearby_obj like '%" + *this.NearbyObj + "%')")
    }
    if this.ObjType != nil {
     condition.WriteString("and (lorry_nearby_target_span.obj_type = " + *this.ObjType + ")")
    }
    if this.ProductName != nil {
      condition.WriteString("and (lorry_nearby_target_span.product_name like '%" + *this.ProductName + "%')")
    }
    if this.Token != nil {
      condition.WriteString("and (lorry_nearby_target_span.token like '%" + *this.Token + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"beginUtc"`
	BeginUtc *int 
	//`json:"dateFlag"`
	DateFlag *string 
	//`json:"endUtc"`
	EndUtc *int 
	//`json:"lorryId"`
	LorryId *string 
	//`json:"name"`
	Name *string 
	//`json:"nearbyObj"`
	NearbyObj *string 
	
	//`json:"objType"`
	ObjType *int 
	//`json:"productName"`
	ProductName *string 
	//`json:"token"`
	Token *string 
}
func (this *CreateObj) Convert2LorryNearbyTargetSpan() lorryNearbyTargetSpanModel.LorryNearbyTargetSpan  {
	var lorryNearbyTargetSpan = lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{}

	lorryNearbyTargetSpan.BeginUtc = this.BeginUtc
	lorryNearbyTargetSpan.DateFlag = this.DateFlag
	lorryNearbyTargetSpan.EndUtc = this.EndUtc
	lorryNearbyTargetSpan.LorryId = this.LorryId
	lorryNearbyTargetSpan.Name = this.Name
	lorryNearbyTargetSpan.NearbyObj = this.NearbyObj
	
	lorryNearbyTargetSpan.ObjType = this.ObjType
	lorryNearbyTargetSpan.ProductName = this.ProductName
	lorryNearbyTargetSpan.Token = this.Token
	return lorryNearbyTargetSpan
}

type UpdateObj struct {

	//`json:"beginUtc"`
	BeginUtc *int 
	//`json:"dateFlag"`
	DateFlag *string 
	//`json:"endUtc"`
	EndUtc *int 
	//`json:"lorryId"`
	LorryId *string 
	//`json:"name"`
	Name *string 
	//`json:"nearbyObj"`
	NearbyObj *string 
	
	//`json:"objType"`
	ObjType *int 
	//`json:"productName"`
	ProductName *string 
	//`json:"token"`
	Token *string 
}
func (this *UpdateObj) Convert2LorryNearbyTargetSpan() lorryNearbyTargetSpanModel.LorryNearbyTargetSpan  {
	var lorryNearbyTargetSpan = lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{}

	lorryNearbyTargetSpan.BeginUtc = this.BeginUtc
	lorryNearbyTargetSpan.DateFlag = this.DateFlag
	lorryNearbyTargetSpan.EndUtc = this.EndUtc
	lorryNearbyTargetSpan.LorryId = this.LorryId
	lorryNearbyTargetSpan.Name = this.Name
	lorryNearbyTargetSpan.NearbyObj = this.NearbyObj
	
	lorryNearbyTargetSpan.ObjType = this.ObjType
	lorryNearbyTargetSpan.ProductName = this.ProductName
	lorryNearbyTargetSpan.Token = this.Token
	return lorryNearbyTargetSpan
}
