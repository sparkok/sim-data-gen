package test

import (
	. "sim_data_gen/utils"
	"fmt"
	"testing"
)

func TestThingSearch(t *testing.T) {
	// 初始化搜索器
	searcher := NewThingSearch(8)

	// 添加测试点
	points := []struct {
		id   string
		x, y float64
	}{
		{"A", 1.0, 2.0},
		{"B", 3.0, 4.0},
		{"C", 1.5, 2.5},
		{"D", 5.0, 5.0},
	}
	for _, p := range points {
		searcher.AddPoint(p.id, p.x, p.y)
	}

	// 查询以 (1.2, 2.3) 为圆心，半径2.0的圆形区域
	if result, err := searcher.SearchCircle(1.2, 2.3, 2.0); err == nil {
		fmt.Println("找到", len(result), "个点:")
		for _, p := range result {
			fmt.Printf("- %s\n", p.ID)
		}
	} else {
		fmt.Println("查询失败:", err)
	}
}
