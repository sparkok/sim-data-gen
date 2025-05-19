package controllers

import (
	mineProductDao "sim_data_gen/daos/mineProduct"
	"sim_data_gen/models/mineChangeStatus"
	. "sim_data_gen/services"
	. "sim_data_gen/utils"
	"time"

	"github.com/beego/beego/v2/server/web"
)

// Operations about SetOfBoomGroups
type SimulateDataController struct {
	web.Controller
}

func (this *SimulateDataController) GenSimulateData(request map[string]interface{}) (string, *string, error) {
	//这是一个测试函数,为了方便
	//sim := simulator.GetDataSimulator()
	var startWeighTime time.Time
	var err error
	productName := extractStrParam(request, "productName", "")
	calculateFromStart := extractBoolParam(request, "calculateFromStart", true)
	minProbability := extractIntParam(request, "minProbability", mineChangeStatus.HIGH_PROBABILITY)
	predictTime := extractDateTimeParam(request, "supposedNow", ProductionsContentSvr.GetNow4FetchData())
	var maxTimeAsUtc *int = RefInt(int(predictTime.Unix()))
	//if _, existed := request["supposedNow"]; existed {
	//	maxTimeAsUtc = RefInt(int(predictTime.Unix()))
	//} else {
	//	maxTimeAsUtc = nil
	//}
	//这是根据本地时区显示的时间,假设docker的时区没有设置好,会显示错误的日期,因此这里使用北京时间

	dateStrOfFetchData := predictTime.In(SystemZone).Format("2006-01-02")
	dateStrOfSavePlan := time.Now().In(SystemZone).Format("2006-01-02")
	//认为中子仪数据会延迟半分钟收到
	if calculateFromStart {
		startWeighTime = Change2TheDay(time.Date(2000, 0, 0, 5, 0, 0, 0, SystemZone), dateStrOfFetchData)
	} else {
		startWeighTime = startWeighTime.Add(-25 * time.Minute)
	}

	if _, err = mineProductDao.GetObjByName(&productName); err != nil {
		return "", nil, err
	}

	GetProductionsContents().CalculateContentsOfProductions(dateStrOfFetchData, dateStrOfSavePlan, minProbability, true, false, predictTime, ByBindLogger, "", maxTimeAsUtc, productName, nil, false)
	return "", nil, nil
}

func extractStrParam(request map[string]interface{}, key string, defaultValue string) string {
	value, existed := request[key].(string)
	if !existed {
		return defaultValue
	}
	return value
}
func extractBoolParam(request map[string]interface{}, key string, defaultValue bool) bool {
	value, existed := request[key].(bool)
	if !existed {
		return defaultValue
	}
	return value
}
func extractIntParam(request map[string]interface{}, key string, defaultValue int) int {
	value, existed := request[key].(int)
	if !existed {
		return defaultValue
	}
	return value
}
func extractTimeParam(request map[string]interface{}, key string, defaultValue time.Time) time.Time {
	tmpStr, existed := request[key].(string)
	if !existed {
		return defaultValue
	}
	value, err := time.Parse("15:04:05", tmpStr)
	if err != nil {
		return defaultValue
	}
	return value
}
func extractDateTimeParam(request map[string]interface{}, key string, defaultValue time.Time) time.Time {
	tmpStr, existed := request[key].(string)
	if !existed {
		return defaultValue
	}
	value, err := time.ParseInLocation("2006-01-02 15:04:05", tmpStr, SystemZone)
	if err != nil {
		return defaultValue
	}
	return value
}

func init() {
	web.Router("/sim_data_gen/setofboomgroups/genSimulateData", &SimulateDataController{}, "get,post:GenSimulateData")
}
