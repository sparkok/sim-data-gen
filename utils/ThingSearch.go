package utils

import (
	"errors"
	"github.com/dhconnelly/rtreego"
	"math"
	"sync"
)

// --------------------------
//      核心实现
// --------------------------

// ThingPoint 表示一个二维点
type ThingPoint struct {
	ID  string
	Loc rtreego.Point
}
type SearchResult struct {
	ID        string
	Distance  float64
	Elevation *float64
}

// 实现 rtreego.Spatial 接口
func (p *ThingPoint) Bounds() rtreego.Rect {
	rect, _ := rtreego.NewRect(
		rtreego.Point{p.Loc[0], p.Loc[1]},
		[]float64{1e-9, 1e-9}, // 极小矩形模拟点
	)
	return rect
}

// ThingSearch 主结构
type ThingSearch struct {
	tree     *rtreego.Rtree         // R树实例
	pointMap map[string]*ThingPoint // ID到点的映射
	lock     sync.RWMutex           // 读写锁
}

// 自动创建R树（简化参数）
func NewThingSearch(aboutSize int) *ThingSearch {
	min, max := autoCalcParams(aboutSize)
	return &ThingSearch{
		tree:     rtreego.NewTree(2, min, max), // 固定参数，适合中小数据量
		pointMap: make(map[string]*ThingPoint),
	}
}

// --------------------------
//      增删改操作
// --------------------------

// 添加点（线程安全）
func (t *ThingSearch) AddPoint(id string, x, y float64) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, exists := t.pointMap[id]; !exists {
		p := &ThingPoint{ID: id, Loc: rtreego.Point{x, y}}
		t.tree.Insert(p)
		t.pointMap[id] = p
	}
}

// 删除点
func (t *ThingSearch) RemovePoint(id string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	p, exists := t.pointMap[id]
	if !exists {
		return errors.New("point not found")
	}

	if t.tree.Delete(p) {
		delete(t.pointMap, id)
		return nil
	}
	return errors.New("delete failed")
}

// 更新点坐标
func (t *ThingSearch) UpdatePoint(id string, newX, newY float64) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	p, exists := t.pointMap[id]
	if !exists {
		return errors.New("point not found")
	}

	// 先删旧点
	if !t.tree.Delete(p) {
		return errors.New("delete old point failed")
	}

	// 插入新点
	newP := &ThingPoint{ID: id, Loc: rtreego.Point{newX, newY}}
	t.tree.Insert(newP)
	t.pointMap[id] = newP
	return nil
}

// --------------------------
//      圆形范围查询
// --------------------------

// 查询圆形区域内的所有点
func (t *ThingSearch) SearchCircle(cx, cy, radius float64) ([]*SearchResult, error) {
	if radius <= 0 {
		return nil, errors.New("radius must be positive")
	}

	t.lock.RLock()
	defer t.lock.RUnlock()

	// 生成包围圆的最小矩形
	searchRect := t.createBoundingRect(cx, cy, radius)

	// 先进行矩形范围查询
	candidates := t.tree.SearchIntersect(searchRect)

	// 精确筛选圆形内的点
	results := make([]*SearchResult, 0)
	center := rtreego.Point{cx, cy}
	distance := 0.0
	for _, item := range candidates {
		p := item.(*ThingPoint)
		distance = t.distance(p.Loc, center, radius)
		if distance < radius {
			result := &SearchResult{ID: p.ID, Distance: distance}
			results = append(results, result)
		}
	}
	return results, nil
}

// 创建包围圆的外接矩形
func (t *ThingSearch) createBoundingRect(cx, cy, r float64) rtreego.Rect {
	lower := rtreego.Point{cx - r, cy - r}
	upper := rtreego.Point{cx + r, cy + r}
	rect, _ := rtreego.NewRectFromPoints(lower, upper)
	return rect
}

// 智能参数计算（核心算法）
func autoCalcParams(n int) (min, max int) {
	switch {
	case n <= 0:
		return 2, 4 // 默认值
	case n <= 8:
		return 2, 4 // 小数据优化
	default:
		base := math.Log2(float64(n))
		min = int(math.Ceil(base))
		max = 2 * min

		// 限制最大值防止过度分裂
		if max > 50 {
			max = 50
			min = 25
		}
		return
	}
}

// 判断点是否在圆内
func (t *ThingSearch) distance(p, center rtreego.Point, radius float64) float64 {
	dx := p[0] - center[0]
	dy := p[1] - center[1]
	return math.Sqrt(dx*dx + dy*dy)
}
