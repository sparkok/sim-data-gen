package common

import "fmt"

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    //`json:"page" form:"page"`         // 页码
	PageSize int    //`json:"pageSize" form:"pageSize"` // 每页大小
}
type PageObjBase struct {
	PageInfo PageInfo
	OrderInfo OrderInfo
}

func (this *PageObjBase)GetOrder() string {
	if(len(this.OrderInfo.OrderKey) == 0) {
		return ""
	}
	if (this.OrderInfo.Desc) {
		return fmt.Sprintf(" order by %v desc", this.OrderInfo.OrderKey)
	} else {
		return fmt.Sprintf("order by %v asc", this.OrderInfo.OrderKey)
	}
}
type OrderInfo struct {
	OrderKey string `json:"OrderKey,omitempty"` //排序使用的key
	Desc bool `json:"Desc,omitempty"` //是否是降序
}

// GetById Find by id structure
type GetById struct {
	Id *string //`json:"id" form:"id"` // 主键ID
}

type IdsReq struct {
	Ids []*string `json:"ids" form:"ids"`
}

type Empty struct{}
//找一个一般人不用的数当空值
const EmptyNumber = 0x31415929