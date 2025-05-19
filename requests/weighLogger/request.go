package weighLogger
import (
	"bytes"
	"sim_data_gen/requests/common"
	weighLoggerModel "sim_data_gen/models/weighLogger"
)
import "time"


/**
* WeighLogger请求类 
*/
type PageObj struct {
	common.PageObjBase
	SearchInfo SearchInfo
}
type SearchInfo struct {
 
	//`json:"busiNo"`
	BusiNo *string  
	//`json:"checkTime"`
	CheckTime *string 
	 
	//`json:"direction"`
	Direction *string  
	//`json:"grossWeight"`
	GrossWeight *string  
	//`json:"netWeight"`
	NetWeight *string  
	//`json:"nt"`
	Nt *string  
	//`json:"siteCode"`
	SiteCode *string  
	//`json:"siteName"`
	SiteName *string  
	//`json:"tareWeight"`
	TareWeight *string  
	//`json:"token"`
	Token *string  
	//`json:"updateAt"`
	UpdateAt *string  
	//`json:"vehicleNo"`
	VehicleNo *string 
}

func (this *SearchInfo) GetConditions() string{
    var condition bytes.Buffer
    if this.BusiNo != nil {
      condition.WriteString("and (weigh_logger.busi_no like '%" + *this.BusiNo + "%')")
    }
    if this.CheckTime != nil {
      condition.WriteString("and (weigh_logger.check_time like '%" + *this.CheckTime + "%')")
    }
    if this.Direction != nil {
      condition.WriteString("and (weigh_logger.direction = '" + *this.Direction + "')")
    }
    if this.GrossWeight != nil {
     condition.WriteString("and (weigh_logger.gross_weight = " + *this.GrossWeight + ")")
    }
    if this.NetWeight != nil {
     condition.WriteString("and (weigh_logger.net_weight = " + *this.NetWeight + ")")
    }
    if this.Nt != nil {
      condition.WriteString("and (weigh_logger.nt like '%" + *this.Nt + "%')")
    }
    if this.SiteCode != nil {
      condition.WriteString("and (weigh_logger.site_code like '%" + *this.SiteCode + "%')")
    }
    if this.SiteName != nil {
      condition.WriteString("and (weigh_logger.site_name like '%" + *this.SiteName + "%')")
    }
    if this.TareWeight != nil {
     condition.WriteString("and (weigh_logger.tare_weight = " + *this.TareWeight + ")")
    }
    if this.Token != nil {
      condition.WriteString("and (weigh_logger.token like '%" + *this.Token + "%')")
    }
    if this.UpdateAt != nil {
      condition.WriteString("and (weigh_logger.update_at = '" + *this.UpdateAt + "')")
    }
    if this.VehicleNo != nil {
      condition.WriteString("and (weigh_logger.vehicle_no like '%" + *this.VehicleNo + "%')")
    }

	if(condition.Len() > 4) {
		return " where " + condition.String()[4:]
	}
	return ""
}

type CreateObj struct {

	//`json:"busiNo"`
	BusiNo *string 
	//`json:"checkTime"`
	CheckTime *string 
	
	//`json:"direction"`
	Direction *string 
	//`json:"grossWeight"`
	GrossWeight *float64 
	//`json:"netWeight"`
	NetWeight *float64 
	//`json:"nt"`
	Nt *string 
	//`json:"siteCode"`
	SiteCode *string 
	//`json:"siteName"`
	SiteName *string 
	//`json:"tareWeight"`
	TareWeight *float64 
	//`json:"token"`
	Token *string 
	//`json:"updateAt"`
	UpdateAt *time.Time 
	//`json:"vehicleNo"`
	VehicleNo *string 
}
func (this *CreateObj) Convert2WeighLogger() weighLoggerModel.WeighLogger  {
	var weighLogger = weighLoggerModel.WeighLogger{}

	weighLogger.BusiNo = this.BusiNo
	weighLogger.CheckTime = this.CheckTime
	
	weighLogger.Direction = this.Direction
	weighLogger.GrossWeight = this.GrossWeight
	weighLogger.NetWeight = this.NetWeight
	weighLogger.Nt = this.Nt
	weighLogger.SiteCode = this.SiteCode
	weighLogger.SiteName = this.SiteName
	weighLogger.TareWeight = this.TareWeight
	weighLogger.Token = this.Token
	weighLogger.UpdateAt = this.UpdateAt
	weighLogger.VehicleNo = this.VehicleNo
	return weighLogger
}

type UpdateObj struct {

	//`json:"busiNo"`
	BusiNo *string 
	//`json:"checkTime"`
	CheckTime *string 
	
	//`json:"direction"`
	Direction *string 
	//`json:"grossWeight"`
	GrossWeight *float64 
	//`json:"netWeight"`
	NetWeight *float64 
	//`json:"nt"`
	Nt *string 
	//`json:"siteCode"`
	SiteCode *string 
	//`json:"siteName"`
	SiteName *string 
	//`json:"tareWeight"`
	TareWeight *float64 
	//`json:"token"`
	Token *string 
	//`json:"updateAt"`
	UpdateAt *time.Time 
	//`json:"vehicleNo"`
	VehicleNo *string 
}
func (this *UpdateObj) Convert2WeighLogger() weighLoggerModel.WeighLogger  {
	var weighLogger = weighLoggerModel.WeighLogger{}

	weighLogger.BusiNo = this.BusiNo
	weighLogger.CheckTime = this.CheckTime
	
	weighLogger.Direction = this.Direction
	weighLogger.GrossWeight = this.GrossWeight
	weighLogger.NetWeight = this.NetWeight
	weighLogger.Nt = this.Nt
	weighLogger.SiteCode = this.SiteCode
	weighLogger.SiteName = this.SiteName
	weighLogger.TareWeight = this.TareWeight
	weighLogger.Token = this.Token
	weighLogger.UpdateAt = this.UpdateAt
	weighLogger.VehicleNo = this.VehicleNo
	return weighLogger
}
