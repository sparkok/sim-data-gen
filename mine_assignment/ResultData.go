package mine_assignment

// Response 包含整个JSON响应的结构体
type Response struct {
	Code  int          `json:"code"`
	Value ElementValue `json:"value"`
}
type ElementValue struct {
	Message string       `json:"msg"`
	Result  []ResultData `json:"result"`
}

// BoomGroupOccupation 对应 "boomGroupOccupations" 数组中的元素结构体
type ElementBoomGroupOccupation struct {
	Group      string  `json:"group"`
	Occupation float64 `json:"occupation"`
}

// ContentPercent 对应 "contentPercents" 数组中的元素结构体
type ElementContentPercent struct {
	MaterialTotal float64  `json:"materialTotal"`
	Name          string   `json:"name"`
	Total         *float64 `json:"total"`
}

// ResultData 整体对应的结构体，包含了所有JSON中的字段
type ResultData struct {
	BoomGroupOccupations []ElementBoomGroupOccupation `json:"boomGroupOccupations"`
	ContentPercents      []ElementContentPercent      `json:"contentPercents"`
	Goal                 float64                      `json:"goal"`
	MD5                  string                       `json:"md5"`
	Result               bool                         `json:"result"`
	SolveResult          float64                      `json:"solveResult"`
	TakenSeconds         int64                        `json:"takenSeconds"`
	TotalEffect          float64                      `json:"totalEffect"`
	TotalMass            float64                      `json:"totalMass"`
	Value                float64                      `json:"value"`
}
