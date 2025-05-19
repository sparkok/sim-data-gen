package message

type Msg struct {
	Params []interface{} `json:"params"`
	Id     string        `json:"id"`
	Obj    string        `json:"obj"`
	Method string        `json:"method"`
}
type MsgResp struct {
	Code   int          `json:"code"`
	Id     string       `json:"id"`
	Ret    *interface{} `json:"ret"`
	Method string       `json:"method"`
}
