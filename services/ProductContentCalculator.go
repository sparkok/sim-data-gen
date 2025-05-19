package services

import (
	mine_assignment "sim_data_gen/mine_assignment"
	boomGroupInfoModel "sim_data_gen/models/boomGroupInfo"
	productAndBoomGroupModel "sim_data_gen/models/productAndBoomGroup"
	"sim_data_gen/utils"
)

// ProductionContentCalculator 用于计算生产内容的计算器
type ProductionContentCalculator struct {
	// ContentPrevMap 存储之前的生产内容信息
	ContentPrevMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully
	// RealTimeContentMap 存储实时的生产内容信息
	RealTimeContentMap map[string]boomGroupInfoModel.BoomGroupInfo
}

// AfterCalculate 是一个回调函数类型，用于在计算完成后执行某些操作
type AfterCalculate func(result bool)

// NewProductionContentCalculator 创建一个新的 ProductionContentCalculator 实例
func NewProductionContentCalculator(prevBoomGroupInfos []productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, realTimeContents []boomGroupInfoModel.BoomGroupInfo) *ProductionContentCalculator {
	contentPrevMap := make(map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully)
	realTimeContentMap := make(map[string]boomGroupInfoModel.BoomGroupInfo)

	// 将之前的生产内容信息填充到 ContentPrevMap 中
	for _, prevBoomGroupInfo := range prevBoomGroupInfos {
		contentPrevMap[*prevBoomGroupInfo.TokenOfBoomGroup] = prevBoomGroupInfo
	}

	// 将实时的生产内容信息填充到 RealTimeContentMap 中
	for _, realTimeContent := range realTimeContents {
		realTimeContentMap[*realTimeContent.TokenOfBoomGroup] = realTimeContent
	}

	// 返回初始化后的 ProductionContentCalculator 实例
	return &ProductionContentCalculator{ContentPrevMap: contentPrevMap, RealTimeContentMap: realTimeContentMap}
}

// Calculate 执行生产内容的计算
func (this *ProductionContentCalculator) Calculate(url string, contentPercentLimit []mine_assignment.ContentPercentLimitReq, diggers []mine_assignment.DiggerReq, boomGroup2Map func(boomGroup productAndBoomGroupModel.ProductAndBoomGroupInDetailFully) map[string]interface{}, boomGroupInfo2Map func(boomGroupInfo boomGroupInfoModel.BoomGroupInfo) map[string]interface{}, afterCalculate func(hasResult bool, boomGroupsMap map[string]productAndBoomGroupModel.ProductAndBoomGroupInDetailFully, results []mine_assignment.ResultData)) {
	go func() {
		// 创建数据
		boomGroups := []map[string]interface{}{}

		// 遍历 ContentPrevMap 和 RealTimeContentMap 进行处理
		for key, value := range this.ContentPrevMap {
			if content, existed := this.RealTimeContentMap[key]; existed {
				boomGroupInfoAsMap := boomGroupInfo2Map(content)
				if boomGroupInfoAsMap != nil {
					if value.Used == nil || value.High == nil {
						utils.Logger.Error("boomGroup used or high is nil")
					}
					if *value.Used >= *value.High {
						continue
					}
					boomGroupInfoAsMap["max"] = *value.High - *value.Used
					boomGroups = append(boomGroups, boomGroupInfoAsMap)
				}
			} else {
				boomGroupAsMap := boomGroup2Map(value)
				if boomGroupAsMap != nil {
					boomGroups = append(boomGroups, boomGroupAsMap)
				}
			}
		}

		// 查询分配结果并执行回调函数
		results, ok := mine_assignment.QueryAssignMine4BoomGroupsMap(url, contentPercentLimit, diggers, boomGroups)
		afterCalculate(ok, this.ContentPrevMap, results)
	}()
}
