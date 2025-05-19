package services

import (
	. "sim_data_gen/entity"
	. "sim_data_gen/meterial"
	boomGroupInfoModel "sim_data_gen/models/boomGroupInfo"
	productAndBoomGroupModel "sim_data_gen/models/productAndBoomGroup"
	"sim_data_gen/utils"
	"gonum.org/v1/gonum/stat"
)

func GetMaterialContentInBoomGroupInfo(materialAttrib MaterialAttrib, boomGroupInfo boomGroupInfoModel.BoomGroupInfo) float64 {
	switch materialAttrib.Index {
	case 1:
		return utils.FloatPtrValue(boomGroupInfo.Material1, -1)
	case 2:
		return utils.FloatPtrValue(boomGroupInfo.Material2, -1)
	case 3:
		return utils.FloatPtrValue(boomGroupInfo.Material3, -1)
	case 4:
		return utils.FloatPtrValue(boomGroupInfo.Material4, -1)
	case 5:
		return utils.FloatPtrValue(boomGroupInfo.Material5, -1)
	case 6:
		return utils.FloatPtrValue(boomGroupInfo.Material6, -1)
	case 7:
		return utils.FloatPtrValue(boomGroupInfo.Material7, -1)
	case 8:
		return utils.FloatPtrValue(boomGroupInfo.Material8, -1)
	case 9:
		return utils.FloatPtrValue(boomGroupInfo.Material9, -1)
	case 10:
		return utils.FloatPtrValue(boomGroupInfo.Material10, -1)
	case 11:
		return utils.FloatPtrValue(boomGroupInfo.Material11, -1)
	case 12:
		return utils.FloatPtrValue(boomGroupInfo.Material12, -1)
	case 13:
		return utils.FloatPtrValue(boomGroupInfo.Material13, -1)
	case 14:
		return utils.FloatPtrValue(boomGroupInfo.Material14, -1)
	case 15:
		return utils.FloatPtrValue(boomGroupInfo.Material15, -1)
	case 16:
		return utils.FloatPtrValue(boomGroupInfo.Material16, -1)
	case 17:
		return utils.FloatPtrValue(boomGroupInfo.Material17, -1)
	case 18:
		return utils.FloatPtrValue(boomGroupInfo.Material18, -1)
	case 19:
		return utils.FloatPtrValue(boomGroupInfo.Material19, -1)
	case 20:
		return utils.FloatPtrValue(boomGroupInfo.Material20, -1)
	}
	return -1
}
func GetMaterialContentInBoomGroup(materialAttrib MaterialAttrib, boomGroup productAndBoomGroupModel.ProductAndBoomGroupInDetailFully) float64 {
	switch materialAttrib.Index {
	case 1:
		return utils.FloatPtrValue(boomGroup.Material1, -1)
	case 2:
		return utils.FloatPtrValue(boomGroup.Material2, -1)
	case 3:
		return utils.FloatPtrValue(boomGroup.Material3, -1)
	case 4:
		return utils.FloatPtrValue(boomGroup.Material4, -1)
	case 5:
		return utils.FloatPtrValue(boomGroup.Material5, -1)
	case 6:
		return utils.FloatPtrValue(boomGroup.Material6, -1)
	case 7:
		return utils.FloatPtrValue(boomGroup.Material7, -1)
	case 8:
		return utils.FloatPtrValue(boomGroup.Material8, -1)
	case 9:
		return utils.FloatPtrValue(boomGroup.Material9, -1)
	case 10:
		return utils.FloatPtrValue(boomGroup.Material10, -1)
	case 11:
		return utils.FloatPtrValue(boomGroup.Material11, -1)
	case 12:
		return utils.FloatPtrValue(boomGroup.Material12, -1)
	case 13:
		return utils.FloatPtrValue(boomGroup.Material13, -1)
	case 14:
		return utils.FloatPtrValue(boomGroup.Material14, -1)
	case 15:
		return utils.FloatPtrValue(boomGroup.Material15, -1)
	case 16:
		return utils.FloatPtrValue(boomGroup.Material16, -1)
	case 17:
		return utils.FloatPtrValue(boomGroup.Material17, -1)
	case 18:
		return utils.FloatPtrValue(boomGroup.Material18, -1)
	case 19:
		return utils.FloatPtrValue(boomGroup.Material19, -1)
	case 20:
		return utils.FloatPtrValue(boomGroup.Material20, -1)
	}
	return -1
}

func CalculateContent(correctResultArray []ContentResult, materialAttrib MaterialAttrib) [2]float64 {
	data := make([]float64, len(correctResultArray))
	weight := make([]float64, len(correctResultArray))
	for i, dataItem := range correctResultArray {
		val := GetMaterialContentInResult(materialAttrib, dataItem)
		if val < 0 {
			if materialAttrib.Required != nil && *materialAttrib.Required {
				utils.Logger.Error("Material Content is error!")
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

// EstimateAndProbability 根据一组数据计算给定值的估计值和概率
// data 是原始数据数组，x 是要计算概率的目标值
//
//	func EstimateAndProbability(data []float64, x float64) (float64, float64) {
//		mean := stat.Mean(data, nil)
//		std := stat.StdDev(data, nil)
//		zScore := (x - mean) / std
//
//		// 使用distuv.Normal创建标准正态分布对象，再调用CDF方法
//		normalDist := distuv.Normal{Mu: 0, Sigma: 1}
//		probability := normalDist.CDF(zScore)
//
//		return mean, probability
//	}
func FillMaterialContentInResult(boomGroupInfo *boomGroupInfoModel.BoomGroupInfo, materialAttrib MaterialAttrib, value [2]float64) float64 {
	// ProbabilityX - 这里直接放入准确率
	switch materialAttrib.Index {
	case 1:
		boomGroupInfo.Material1 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability1 = utils.RefFloat64(value[1])
	case 2:
		boomGroupInfo.Material2 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability2 = utils.RefFloat64(value[1])
	case 3:
		boomGroupInfo.Material3 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability3 = utils.RefFloat64(value[1])
	case 4:
		boomGroupInfo.Material4 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability4 = utils.RefFloat64(value[1])
	case 5:
		boomGroupInfo.Material5 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability5 = utils.RefFloat64(value[1])
	case 6:
		boomGroupInfo.Material6 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability6 = utils.RefFloat64(value[1])
	case 7:
		boomGroupInfo.Material7 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability7 = utils.RefFloat64(value[1])
	case 8:
		boomGroupInfo.Material8 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability8 = utils.RefFloat64(value[1])
	case 9:
		boomGroupInfo.Material9 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability9 = utils.RefFloat64(value[1])
	case 10:
		boomGroupInfo.Material10 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability10 = utils.RefFloat64(value[1])
	case 11:
		boomGroupInfo.Material11 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability11 = utils.RefFloat64(value[1])
	case 12:
		boomGroupInfo.Material12 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability12 = utils.RefFloat64(value[1])
	case 13:
		boomGroupInfo.Material13 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability13 = utils.RefFloat64(value[1])
	case 14:
		boomGroupInfo.Material14 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability14 = utils.RefFloat64(value[1])
	case 15:
		boomGroupInfo.Material15 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability15 = utils.RefFloat64(value[1])
	case 16:
		boomGroupInfo.Material16 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability16 = utils.RefFloat64(value[1])
	case 17:
		boomGroupInfo.Material17 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability17 = utils.RefFloat64(value[1])
	case 18:
		boomGroupInfo.Material18 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability18 = utils.RefFloat64(value[1])
	case 19:
		boomGroupInfo.Material19 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability19 = utils.RefFloat64(value[1])
	case 20:
		boomGroupInfo.Material20 = utils.RefFloat64(value[0])
		boomGroupInfo.Probability20 = utils.RefFloat64(value[1])
	}
	return -1
}
func GetMaterialContentInResult(materialAttrib MaterialAttrib, result ContentResult) float64 {
	switch materialAttrib.Index {
	case 1:
		return utils.FloatPtrValue(result.NeutronData.Mat1, -1)
	case 2:
		return utils.FloatPtrValue(result.NeutronData.Mat2, -1)
	case 3:
		return utils.FloatPtrValue(result.NeutronData.Mat3, -1)
	case 4:
		return utils.FloatPtrValue(result.NeutronData.Mat4, -1)
	case 5:
		return utils.FloatPtrValue(result.NeutronData.Mat5, -1)
	case 6:
		return utils.FloatPtrValue(result.NeutronData.Mat6, -1)
	case 7:
		return utils.FloatPtrValue(result.NeutronData.Mat7, -1)
	case 8:
		return utils.FloatPtrValue(result.NeutronData.Mat8, -1)
	case 9:
		return utils.FloatPtrValue(result.NeutronData.Mat9, -1)
	case 10:
		return utils.FloatPtrValue(result.NeutronData.Mat10, -1)
	case 11:
		return utils.FloatPtrValue(result.NeutronData.Mat11, -1)
	case 12:
		return utils.FloatPtrValue(result.NeutronData.Mat12, -1)
	case 13:
		return utils.FloatPtrValue(result.NeutronData.Mat13, -1)
	case 14:
		return utils.FloatPtrValue(result.NeutronData.Mat14, -1)
	case 15:
		return utils.FloatPtrValue(result.NeutronData.Mat15, -1)
	case 16:
		return utils.FloatPtrValue(result.NeutronData.Mat16, -1)
	case 17:
		return utils.FloatPtrValue(result.NeutronData.Mat17, -1)
	case 18:
		return utils.FloatPtrValue(result.NeutronData.Mat18, -1)
	case 19:
		return utils.FloatPtrValue(result.NeutronData.Mat19, -1)
	case 20:
		return utils.FloatPtrValue(result.NeutronData.Mat20, -1)
	}
	return -1
}
