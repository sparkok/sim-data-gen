package weighLogger

import "time"

/**
* 实体类 WeighLogger
 */
type WeighLogger struct {
	BusiNo       *string    `json:"busiNo"`                        //业务号 进矿区和出矿区是一个业务号
	CheckTime    *string    `json:"checkTime"`                     //车辆过磅的时间
	Direction    *string    `json:"direction"`                     //1 - 车辆是重车离开矿区,0 - 车辆是重车进入矿区
	GrossWeight  *float64   `gorm:"default:0;" json:"grossWeight"` //毛重
	NetWeight    *float64   `gorm:"default:0;" json:"netWeight"`   //净重
	Nt           *string    `json:"nt"`                            //备注
	SiteCode     *string    `json:"siteCode"`                      //和产品编号对应 例如:CP1
	SiteName     *string    `json:"siteName"`                      //和产品编号对应 例如:CP1
	TareWeight   *float64   `gorm:"default:0;" json:"tareWeight"`  //皮重
	Token        *string    `gorm:"primaryKey;" json:"token"`      //主键
	UpdateAt     *time.Time `gorm:"default:CURRENT_TIMESTAMP;" json:"updateAt"`
	VehicleNo    *string    `json:"vehicleNo"`     //这里填写卡车的名词而不是id
	MaterialCode *string    `json:"material_code"` //三明矿: 0102010003 或 0102010023
	MaterialName *string    `json:"material_name"` //三明矿: 石灰石碎石 或 废石
}
