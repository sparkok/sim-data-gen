package entity

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"sort"
	"time"
)

type ContentResultList struct {
	CorrectResults    map[string][]ContentResult `json:"correctResults"`
	LorryCarries      []LorryCarry               `json:"lorryCarries"`
	Lorry2BoomGroupId map[string]string          `json:"Lorry2BoomGroupId"`
}

func (t *ContentResultList) AddResult(boomGroupId string, result ContentResult) {
	if t.CorrectResults == nil {
		t.CorrectResults = make(map[string][]ContentResult)
	}
	boomGroupWithResult := t.CorrectResults[boomGroupId]
	if boomGroupWithResult == nil {
		t.CorrectResults[boomGroupId] = []ContentResult{result}
		//应该改成根据配矿单元id建立map
	} else {
		t.CorrectResults[boomGroupId] = append(t.CorrectResults[boomGroupId], result)
	}
}
func (t *ContentResultList) Results2Txt() string {
	value, err := json.Marshal(t.CorrectResults)
	if err == nil {
		return ""
	}
	return string(value)
}

func (t *ContentResultList) Enumerate(f func(boomGroupId string, results []ContentResult)) {
	for boomGroupId, results := range t.CorrectResults {
		f(boomGroupId, results)
	}
}
func (t *ContentResultList) EnumerateLastLorryCarries(
	latest func(boomGroupId string, results []ContentResult, contentContentAndP [2]float64),
) {
	lastLorryCarries := map[string]LorryCarry{}
	for _, lorryCarry := range t.LorryCarries {
		if !lorryCarry.IncludeCorrectResult() {
			continue
		}
		lastLorryCarries[lorryCarry.BoomGroupId] = lorryCarry
	}
	for _, lastLorryCarry := range lastLorryCarries {
		latest(lastLorryCarry.BoomGroupId, lastLorryCarry.Results, lastLorryCarry.PrimaryMaterialContentAndP)
	}
}

func (t *ContentResultList) Size() int {
	return len(t.CorrectResults)
}

func (t *ContentResultList) GenerateReport(path string) {

}

func (t *ContentResultList) AddWholeLorryMine(boomGroupId, lorryName string, results []ContentResult, primaryMaterialContentAndP [2]float64) {
	t.LorryCarries = append(t.LorryCarries, LorryCarry{
		CarryId:                    results[0].CarryId,
		Results:                    results,
		BoomGroupId:                boomGroupId,
		LorryName:                  lorryName,
		PrimaryMaterialContentAndP: primaryMaterialContentAndP,
	})
}

type Item struct {
	Purity float64
}

const MIN_PURITY = 0.80

func FilterLowPurityResult(results []ContentResult) []ContentResult {
	rst := make([]ContentResult, 0)
	for _, result := range results {
		if result.Purity > MIN_PURITY {
			rst = append(rst, result)
		}
	}
	return rst
}

func (t *ContentResultList) GenerateExcel(contents map[string][2]float64, path string) (string, error) {
	// 删除已存在的Excel文件
	var err error
	if _, err = os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			return "", err
		}
	}

	// 创建一个新的Excel文件
	f := excelize.NewFile()
	// 创建一个工作表
	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", err
	}
	if index == 0 {
		f.SetActiveSheet(index)
	}

	var headers []string
	// 车牌,品位准确度,车次,氧化钙含量,重量,时间
	headers = []string{"序号", "配矿单元", "卡车", "准确度", "车次", "氧化钙", "地磅矿量", "地磅时间", "中子仪流量", "中子仪时间", "氧化钙预测值", "标准差"}

	// 设置表头
	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s%d", colName, 1) // 生成如"A1"
		f.SetCellValue(sheetName, cell, header)
	}

	carried2Contents := make(map[int][2]float64)
	for _, result := range t.LorryCarries {
		carried2Contents[result.CarryId] = result.PrimaryMaterialContentAndP
	}

	// 写入数据
	row := 2
	// 将数据提取到一个切片中以便排序
	var data [][]interface{}
	var contentAndP [2]float64
	var ok bool
	for boomGroupId, result := range t.CorrectResults {
		for _, item := range result {
			if contentAndP, ok = carried2Contents[item.CarryId]; !ok {
				contentAndP[0] = -1
				contentAndP[1] = -1
			}
			itemValue := []interface{}{
				item.Index,
				boomGroupId,
				*item.VehicleNo,
				item.Purity,
				item.CarryId,
				*item.NeutronData.Mat1,
				*item.BridgeData.NetWeight,
				*item.BridgeData.CheckTime,
				*item.NeutronData.Flux,
				UnixTime2StrTime(item.NeutronData.TestAt),
				contentAndP[0],
				contentAndP[1],
			}
			if item.Purity < MIN_PURITY {
				itemValue[10] = nil
				itemValue[11] = nil
			}
			data = append(data, itemValue)
		}
	}

	//// 按照G列（时间）排序
	sort.Slice(data, func(i, j int) bool {
		return data[i][0].(int) < data[j][0].(int)
	})

	// 写入排序后的数据
	for _, rowdata := range data {
		for col, value := range rowdata {
			colName, _ := excelize.ColumnNumberToName(col + 1)
			cell := fmt.Sprintf("%s%d", colName, row)
			if value != nil {
				f.SetCellValue(sheetName, cell, value)
			} else {
				f.SetCellValue(sheetName, cell, "-")
			}
		}
		row++
	}

	// 写入第二个工作表的数据
	sheetName2 := "最近准确值"
	index2, err := f.NewSheet(sheetName2)
	if err != nil {
		return "", err
	}
	if index2 == 0 {
		f.SetActiveSheet(index2)
	}

	headers = []string{"配矿单元", "氧化钙", "准确度"}

	// 设置表头
	for i, header := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s%d", colName, 1) // 生成如"A1"
		f.SetCellValue(sheetName2, cell, header)
	}

	row = 2
	for boomGroupId, content := range contents {
		// 写入LorryNumber
		f.SetCellValue(sheetName2, fmt.Sprintf("A%d", row), boomGroupId)
		f.SetCellValue(sheetName2, fmt.Sprintf("B%d", row), content[0])
		f.SetCellValue(sheetName2, fmt.Sprintf("C%d", row), content[1])
		row++
	}

	// 保存Excel文件
	if err = f.SaveAs(path); err != nil {
		return "", err
	}

	return path, nil
}

func UnixTime2StrTime(unixTime *int) string {
	if unixTime == nil {
		return ""
	} else {
		// 将 Unix 时间转换为字符串格式的时间
		return time.Unix(int64(*unixTime), 0).Format("2006-01-02 15:04:05")
	}
}
