package common

import "time"

type Count struct {
	Count int64 `orm:"column(Count)"`
}
type DataTime struct {
	Value *time.Time `orm:"column(DataTimeOfNeutron)"`
}
