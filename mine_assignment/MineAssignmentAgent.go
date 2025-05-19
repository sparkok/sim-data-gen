package mine_assignment

import (
	diggerModel "sim_data_gen/models/digger"
	diggerProductBindingModel "sim_data_gen/models/diggerProductBinding"
	"sim_data_gen/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BoomGroupReq struct {
	Name     string  `json:"name"`
	Mat1     float64 `json:"mat1"`
	Mat2     int     `json:"mat2"`
	Mat3     int     `json:"mat3"`
	Min      int     `json:"min"`
	Max      int     `json:"max"`
	Distance int     `json:"distance"`
}

type BoomGroupSettingReq struct {
	MaxCountOfOptionMineGroup int `json:"maxCountOfOptionMineGroup"`
	MinCountOfOptionMineGroup int `json:"minCountOfOptionMineGroup"`
}

type ContentPercentLimitReq struct {
	Max  float64 `json:"max"`
	Min  float64 `json:"min"`
	Name string  `json:"name"`
}

type DiggerReq struct {
	Name  string  `json:"name"`
	Power float64 `json:"power"`
	Token string  `json:"token"`
}

type GoalReq struct {
	CanSubDiggers       int     `json:"canSubDiggers"`
	IdealDigWeight      float64 `json:"idealDigWeight"`
	IdealRunHours       float64 `json:"idealRunHours"`
	MaxBiasInProportion float64 `json:"maxBiasInProportion"`
}

type OutputReq struct {
	Max                float64 `json:"max"`
	Min                float64 `json:"min"`
	MineAssignmentGoal int     `json:"mineAssignmentGoal"`
}

type RuleReq struct {
	Expression string `json:"expression"`
	Name       string `json:"name"`
}

type RequestDataReq struct {
	BoomGroups           []map[string]interface{} `json:"boomGroups"`
	BoomGroupSetting     BoomGroupSettingReq      `json:"boomGroupSetting"`
	ContentPercentLimits []ContentPercentLimitReq `json:"contentPercentLimits"`
	Diggers              []DiggerReq              `json:"diggers"`
	Goal                 GoalReq                  `json:"goal"`
	Md5                  string                   `json:"md5"`
	Output               OutputReq                `json:"output"`
	Prohibitions         []interface{}            `json:"prohibitions"`
	Rules                []RuleReq                `json:"rules"`
}

func QueryTest2(url string) ([]ResultData, bool) {
	boomGroups := []map[string]interface{}{
		map[string]interface{}{
			"Name":     "1#",
			"cao":      47,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "2#",
			"cao":      44,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "3#",
			"cao":      51,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "4#",
			"cao":      49,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "5#",
			"cao":      45,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "6#",
			"cao":      46,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "7#",
			"cao":      47,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "8#",
			"cao":      53,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "9#",
			"cao":      52,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
		map[string]interface{}{
			"Name":     "10#",
			"cao":      46.5,
			"mgo":      0,
			"sm":       0,
			"min":      0,
			"max":      5000,
			"distance": 100,
		},
	}
	contentPercentLimits := []ContentPercentLimitReq{
		{Max: 69.0, Min: 46.0, Name: "cao"},
		{Max: 5.0, Min: 0.0, Name: "sm"},
		{Max: 1.0, Min: 0.0, Name: "mgo"},
	}

	diggers := []DiggerReq{
		{Name: "Digger1", Power: 300.0, Token: "001"},
		{Name: "Digger4", Power: 300.0, Token: "004"},
	}
	return QueryForExample(url, contentPercentLimits, diggers, boomGroups)
}
func QueryForExample(url string, contentPercentLimits []ContentPercentLimitReq, diggers []DiggerReq, boomGroups []map[string]interface{}) ([]ResultData, bool) {
	boomGroupSetting := BoomGroupSettingReq{
		MaxCountOfOptionMineGroup: 0,
		MinCountOfOptionMineGroup: 0,
	}

	//contentPercentLimits := []ContentPercentLimitReq{
	//	{Max: 69.0, Min: 46.0, Name: "cao"},
	//	{Max: 5.0, Min: 0.0, Name: "sm"},
	//	{Max: 1.0, Min: 0.0, Name: "mgo"},
	//}

	goal := GoalReq{
		CanSubDiggers:       0,
		IdealDigWeight:      10000.0,
		IdealRunHours:       8.0,
		MaxBiasInProportion: 0.2,
	}

	//md5 := "4ffd3d42686b5f540aef2d6b170e4e7e"

	output := OutputReq{
		Max:                500000.0,
		Min:                0.0,
		MineAssignmentGoal: 6,
	}

	rules := []RuleReq{
		{Expression: "\tlocal a = VecSum(boomGroupOccupies);\n\tlocal b = 2;\n\treturn (a + b)\n", Name: "rule1"},
	}

	data := RequestDataReq{
		BoomGroups:           boomGroups,
		BoomGroupSetting:     boomGroupSetting,
		ContentPercentLimits: contentPercentLimits,
		Diggers:              diggers,
		Goal:                 goal,
		//Md5:                  md5,
		Output:       output,
		Prohibitions: []interface{}{},
		Rules:        rules,
	}

	// 将数据编码为 JSON
	if jsonData, err := json.MarshalIndent(data, "", "  "); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, false
	} else {
		fmt.Println(string(jsonData))
		return Post(url, jsonData)
	}
}
func QueryAssignMine(url string, params map[string]interface{}) ([]ResultData, bool) {
	if requestAsText, err := json.Marshal(params); err != nil {
		return nil, false
	} else {
		return Post(url, requestAsText)
	}
}
func QueryAssignMine4BoomGroupsMap(url string, contentPercentLimit []ContentPercentLimitReq, diggers []DiggerReq, boomGroups []map[string]interface{}) ([]ResultData, bool) {
	return QueryForExample(url, contentPercentLimit, diggers, boomGroups)
}

func Post(url string, requestAsText []byte) ([]ResultData, bool) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(requestAsText)))
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, false
	}

	req.Header.Set("Content-Type", "application/json")

	rep, err := client.Do(req)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, false
	}
	data, err := io.ReadAll(rep.Body)
	return parseResult(string(data), true)
}

func QueryExample1(url string) ([]ResultData, bool) {
	/*
	  "boomGroups": [
	    {
	      "name": "1#",
	      "cao": 47,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "2#",
	      "cao": 44,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "3#",
	      "cao": 51,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "4#",
	      "cao": 49,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "5#",
	      "cao": 45,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "6#",
	      "cao": 46,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "7#",
	      "cao": 47,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "8#",
	      "cao": 53,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "9#",
	      "cao": 52,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    },
	    {
	      "name": "10#",
	      "cao": 46.5,
	      "mgo": 0,
	      "sm": 0,
	      "min": 0,
	      "max": 5000,
	      "distance": 100
	    }
	  ],
	  "boomGroupSetting": {
	    "maxCountOfOptionMineGroup": 0,
	    "minCountOfOptionMineGroup": 0
	  },
	  "contentPercentLimits": [
	    {
	      "max": 69.0,
	      "min": 46.0,
	      "name": "cao"
	    },
	    {
	      "max": 5.0,
	      "min": 0.0,
	      "name": "sm"
	    },
	    {
	      "max": 1.0,
	      "min": 0.0,
	      "name": "mgo"
	    }
	  ],
	  "diggers": [
	    {
	      "name": "Digger1",
	      "power": 300.0,
	      "token": "001"
	    },
	    {
	      "name": "Digger4",
	      "power": 300.0,
	      "token": "004"
	    }
	  ],
	  "goal": {
	    "canSubDiggers": 0,
	    "idealDigWeight": 10000.0,
	    "idealRunHours": 8.0,
	    "maxBiasInProportion": 0.2
	  },
	  "md5": "4ffd3d42686b5f540aef2d6b170e4e7e",
	  "output": {
	    "max": 500000.0,
	    "min": 0.0,
	    "mineAssignmentGoal": 6
	  },
	  "prohibitions": [],
	  "rules": [
	    {
	      "expression": "\tlocal a = VecSum(boomGroupOccupies);\n\tlocal b = 2;\n\treturn (a + b)\n",
	      "name": "rule1"
	    }
	  ]
	*/
	query := map[string]interface{}{
		"boomGroups": []map[string]interface{}{
			{
				"name":     "1#",
				"cao":      47,
				"mgo":      0,
				"sm":       0,
				"min":      0,
				"max":      5000,
				"distance": 100,
			},
			{
				"name":     "2#",
				"cao":      44,
				"mgo":      0,
				"sm":       0,
				"min":      0,
				"max":      5000,
				"distance": 100,
			},
			{
				"name":     "3#",
				"cao":      51,
				"mgo":      0,
				"sm":       0,
				"min":      0,
				"max":      5000,
				"distance": 100,
			},
			// 此处省略了中间重复格式的部分，实际完整代码需补充完整
			{
				"name":     "10#",
				"cao":      46.5,
				"mgo":      0,
				"sm":       0,
				"min":      0,
				"max":      5000,
				"distance": 100,
			},
		},
		"boomGroupSetting": map[string]interface{}{
			"maxCountOfOptionMineGroup": 0,
			"minCountOfOptionMineGroup": 0,
		},
		"contentPercentLimits": []map[string]interface{}{
			{
				"max":  69.0,
				"min":  46.0,
				"name": "cao",
			},
			{
				"max":  5.0,
				"min":  0.0,
				"name": "sm",
			},
			{
				"max":  1.0,
				"min":  0.0,
				"name": "mgo",
			},
		},
		"diggers": []map[string]interface{}{
			{
				"name":  "Digger1",
				"power": 300.0,
				"token": "001",
			},
			{
				"name":  "Digger4",
				"power": 300.0,
				"token": "004",
			},
		},
		"goal": map[string]interface{}{
			"canSubDiggers":       0,
			"idealDigWeight":      10000.0,
			"idealRunHours":       8.0,
			"maxBiasInProportion": 0.2,
		},
		"md5": "4ffd3d42686b5f540aef2d6b170e4e7e",
		"output": map[string]interface{}{
			"max":                500000.0,
			"min":                0.0,
			"mineAssignmentGoal": 6,
		},
		"prohibitions": []interface{}{},
		"rules": []map[string]interface{}{
			{
				"expression": "\tlocal a = VecSum(boomGroupOccupies);\n\tlocal b = 2;\n\treturn (a + b)\n",
				"name":       "rule1",
			},
		},
	}
	return QueryAssignMine(url, query)
}
func parseResult(data string, debug bool) ([]ResultData, bool) {
	var response Response
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		fmt.Println("解析JSON出错:", err)
		return nil, false
	}
	if response.Code != 0 {
		return nil, false
	}
	if debug {
		// 输出解析后的结构体内容，可按需进行后续操作
		for i, result := range response.Value.Result {
			fmt.Printf("解析后的结果 %d:\n", i+1)
			fmt.Println("BoomGroupOccupations:", result.BoomGroupOccupations)
			fmt.Println("ContentPercents:", result.ContentPercents)
			fmt.Println("Goal:", result.Goal)
			fmt.Println("MD5:", result.MD5)
			fmt.Println("Result:", result.Result)
			fmt.Println("SolveResult:", result.SolveResult)
			fmt.Println("TakenSeconds:", result.TakenSeconds)
			fmt.Println("TotalEffect:", result.TotalEffect)
			fmt.Println("TotalMass:", result.TotalMass)
			fmt.Println("BoomGroupItems:", result.Value)
			fmt.Println()
		}
	}
	return response.Value.Result, true
}

func NewMineAssignmentReqFromDiggerObjs(diggerObjs []diggerModel.Digger) []DiggerReq {
	diggerReq := make([]DiggerReq, len(diggerObjs))
	for i, diggerObj := range diggerObjs {
		diggerReq[i] = DiggerReq{
			Name:  *diggerObj.Name,
			Power: *diggerObj.Speed,
			Token: *diggerObj.Token,
		}
	}
	return diggerReq
}
func NewMineAssignmentReqFromDiggerProductBindings(bindings []diggerProductBindingModel.DiggerProductBindingFully1) []DiggerReq {
	diggerReq := make([]DiggerReq, len(bindings))
	for i, diggerObj := range bindings {
		diggerReq[i] = DiggerReq{
			Name:  *diggerObj.DiggerDesp,
			Power: *diggerObj.DiggerSpeed,
			Token: *diggerObj.TokenOfDigger,
		}
	}
	return diggerReq

}
