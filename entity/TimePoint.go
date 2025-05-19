package entity

import (
	"github.com/google/btree"
)

type TimePoint struct {
	Key   int64
	Value interface{}
}

func (tp TimePoint) Less(other btree.Item) bool {
	return tp.Key < other.(TimePoint).Key
}

type TimeSeries struct {
	tree *btree.BTree
}

func NewTimeSeries() *TimeSeries {
	return &TimeSeries{tree: btree.New(32)}
}

// 插入数据
func (ts *TimeSeries) Insert(tp TimePoint) {
	ts.tree.ReplaceOrInsert(tp)
}

// 查找最接近的时间点
func (ts *TimeSeries) FindClosest(value int64) (closest TimePoint, found bool) {
	//从大到小遍历,找到小于等于的最大值
	found = false
	ts.tree.DescendLessOrEqual(TimePoint{Key: value}, func(item btree.Item) bool {
		found = true
		closest = item.(TimePoint)
		return false // 仅取第一个匹配项
	})
	return closest, found
}
