package snap

import (
	. "sim_data_gen/entity"
	"encoding/json"
	"os"
)

type SnapTask struct {
	TaskId     string   `json:"taskId"`
	TaskName   string   `json:"taskName"`
	TaskType   string   `json:"taskType"`
	TaskStatus string   `json:"taskStatus"`
	TaskDesc   string   `json:"taskDesc"`
	SnapData   []string `json:"SnapData"`
	//Results    ContentResultList `json:"results"`
	DataLoaded bool
}

//func (t *SnapTask) SetResults(results ContentResultList) {
//	t.Results = results
//}

func (t *SnapTask) GenereteReport(contentResultList *ContentResultList, reportFilePath string) {
	contentResultList.GenerateReport(reportFilePath)
}

func (t *SnapTask) GenerateExcel(contentResultList *ContentResultList, contents map[string][2]float64, excelFilePath string) (string, error) {
	return contentResultList.GenerateExcel(contents, excelFilePath)
}

func (t *SnapTask) Save(filePath string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
