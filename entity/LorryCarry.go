package entity

import (
	"sim_data_gen/meterial"
	. "sim_data_gen/utils"
	"gonum.org/v1/gonum/stat"
)

type LorryCarry struct {
	CarryId                    int             `json:"carryId"`
	Results                    []ContentResult `json:"contentResult"`
	BoomGroupId                string          `json:"boomGroupId"`
	LorryName                  string          `json:"lorryName"`
	CarryTime                  string          `json:"carryTime"`
	PrimaryMaterialContentAndP [2]float64      `json:"primaryMaterialContentAndP"`
}

func (t *LorryCarry) IncludeCorrectResult() bool {
	for _, result := range t.Results {
		if AlmostOne(result.Purity) {
			return true
		}
	}
	return false
}
func (t *LorryCarry) CalculateContent(correctResultArray []ContentResult, materialAttrib meterial.MaterialAttrib, GetMaterialContentInResult func(meterial.MaterialAttrib, ContentResult) float64) [2]float64 {
	data := make([]float64, len(correctResultArray))
	weight := make([]float64, len(correctResultArray))
	for i, dataItem := range correctResultArray {
		val := GetMaterialContentInResult(materialAttrib, dataItem)
		if val < 0 {
			if materialAttrib.Required != nil && *materialAttrib.Required {
				Logger.Error("Material Content is error!")
				return [2]float64{0, 0}
			} else {
				val = -1
			}
		}
		weight[i] = dataItem.UboundOfTotalMass - dataItem.LboundOfTotalMass
		data[i] = val
	}
	return [2]float64{stat.Mean(data, weight), stat.StdDev(data, weight)}
}
