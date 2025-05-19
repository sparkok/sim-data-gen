package entity

import (
	"sim_data_gen/models/weighLogger"
	yAnalyserModel "sim_data_gen/models/yAnalyser"
)

type ContentResult struct {
	VehicleNo         *string                  `json:"vehicleNo"`
	LboundOfTotalMass float64                  `json:"LboundOfTotalMass"`
	UboundOfTotalMass float64                  `json:"UboundOfTotalMass"`
	NeutronData       yAnalyserModel.YAnalyser `json:"NeutronData"`
	BridgeData        weighLogger.WeighLogger  `json:"BridgeData"`
	Purity            float64                  `json:"Purity"`
	Name              string                   `json:"Name"`
	CarryId           int                      `json:"CarryId"` // > 0 ,如果运输序号相同,这表示是一车矿石
	Index             int                      `json:"Index"`
}
