package services

import (
	"sim_data_gen/utils"
	"errors"
	"fmt"
	"github.com/dhconnelly/rtreego"
	"github.com/spatial-go/geoos/space"
	"github.com/twpayne/go-geom"
	wktEncoding "github.com/twpayne/go-geom/encoding/wkt"
	"go.uber.org/zap"
	"math"
	"strings"
)

// Polygon 表示一个多边形，包含唯一标识和顶点列表
type Polygon struct {
	ID       string
	Name     string
	Vertices [][]float64
}

// 将 *geom.Polygon 转换为 space.Polygon
func ConvertToSpacePolygon(polygon *geom.Polygon) space.Polygon {
	// 获取多边形的坐标
	coords := polygon.Coords()

	// 创建一个 space.Polygon
	var spacePolygon space.Polygon

	// 遍历每个环（外环和内环）
	for _, ring := range coords {
		var points [][]float64
		// 遍历环中的每个点
		for _, coord := range ring {
			points = append(points, []float64{coord[0], coord[1]})
		}
		spacePolygon = append(spacePolygon, points)
	}

	return spacePolygon
}

// NewPolygonFromWKB 从WKB格式创建多边形

func NewPolygonFromWKT(id string, name string, wkt string) (*Polygon, error) {
	if geomObj, err := wktEncoding.Unmarshal(wkt); err != nil {
		utils.Logger.Error("failed to Unmarshal geom", zap.Error(err))
		//rawJSON, _ = fc.MarshalJSON()
		return nil, errors.New("爆堆的地理信息对象不正确")
	} else {

		polygonPtr := geomObj.(*geom.Polygon)
		outerShell, _ := BufferedOuterShellOfPolygon(polygonPtr, 15)
		coords := ExtractCoordsFromLineStringForFloats(outerShell, 3857)
		var vertices = make([][]float64, len(coords))
		for i, _ := range coords {
			//变成百万分之一度
			vertices[i] = []float64{coords[i][0] * 1000000, coords[i][1] * 1000000}
		}
		//配矿单元只有一个多边形,也就是最外层的
		return &Polygon{ID: id, Name: name, Vertices: vertices}, nil
	}
}
func BufferedOuterShellOfPolygon(poly *geom.Polygon, bufferDist float64) ([][]float64, error) {
	newPolygon := ConvertToSpacePolygon(poly)
	ret := newPolygon.Buffer(bufferDist, 8).(space.Polygon)
	return ret[0], nil
	//// 将go-geom的Polygon转换为geoos的Polygon
	//srcCoords := poly.Coords()[0]
	//pointCount := len(srcCoords)
	//outerShell := make(space.LineString, pointCount)
	//
	//// 遍历所有环（虽然用户已确保无洞，但需保持结构）
	//for i := 0; i < pointCount; i++ {
	//	outerShell[i] = space.Point{srcCoords[i].X(), srcCoords[i].Y()}
	//}
	//bufferedShell := outerShell.Buffer(bufferDist, 8)
	////每90度8个边
	//buffer, ok := bufferedShell.(space.Polygon)
	//if !ok {
	//	return nil, fmt.Errorf("缓冲区生成结果非Polygon类型")
	//}
	//return buffer[0], nil

}

//func BufferedOuterShellOfPolygon(poly *geom.Polygon, bufferDist float64) ([][]float64, error) {
//
//	// 将go-geom的Polygon转换为geoos的Polygon
//	srcCoords := poly.Coords()[0]
//	pointCount := len(srcCoords)
//	outerShell := make(space.LineString, pointCount)
//
//	// 遍历所有环（虽然用户已确保无洞，但需保持结构）
//	for i := 0; i < pointCount; i++ {
//		outerShell[i] = space.Point{srcCoords[i].X(), srcCoords[i].Y()}
//	}
//	bufferedShell := outerShell.Buffer(bufferDist, 8)
//	//每90度8个边
//	buffer, ok := bufferedShell.(space.Polygon)
//	if !ok {
//		return nil, fmt.Errorf("缓冲区生成结果非Polygon类型")
//	}
//	return buffer[0], nil
//
//}

// Bounds 计算多边形外接矩形
func (p *Polygon) Bounds() rtreego.Rect {
	if len(p.Vertices) == 0 {
		rect, _ := rtreego.NewRect(rtreego.Point{0, 0}, []float64{0, 0})
		return rect
	}

	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)

	for _, pt := range p.Vertices {
		if len(pt) < 2 {
			continue // 跳过无效点
		}
		x := pt[0]
		y := pt[1]
		minX = math.Min(minX, x)
		maxX = math.Max(maxX, x)
		minY = math.Min(minY, y)
		maxY = math.Max(maxY, y)
	}

	rect, _ := rtreego.NewRect(
		rtreego.Point{minX, minY},
		[]float64{maxX - minX, maxY - minY},
	)
	return rect
}
func NewLineString(coords [][]float64) (space.LineString, error) {
	var ls space.LineString
	for _, c := range coords {
		if len(c) != 2 {
			return nil, fmt.Errorf("invalid coordinate length")
		}
		ls = append(ls, c)
	}
	return ls, nil
}
func LineStringToWKT(vertices [][]float64) string {
	if len(vertices) == 0 {
		return "LINESTRING EMPTY"
	}
	var wkt strings.Builder
	wkt.WriteString("LINESTRING (")
	for i, coord := range vertices {
		if len(coord) >= 2 {
			if i > 0 {
				wkt.WriteString(",\n ")
			}
			wkt.WriteString(fmt.Sprintf("%f %f", coord[0], coord[1]))
		}
	}
	wkt.WriteString(")")
	return wkt.String()
}
func PointInLineString(point rtreego.Point, vertices [][]float64) bool {
	if len(point) < 2 || len(vertices) < 2 {
		return false
	}
	//utils.Logger.Info("PointInLineString", zap.String("point", PointToWKT(point)), zap.String("vertices", LineStringToWKT(vertices)))
	//utils.Logger.Info(fmt.Sprintf(`SELECT ST_Contains(ST_GeomFromText('%s'),\nST_GeomFromText('%s')) AS is_point_on_line;`, LineStringToWKT(vertices), PointToWKT(point)))
	x, y := point[0], point[1]
	n := len(vertices)
	inside := false

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		vi := vertices[i]
		vj := vertices[j]

		// 跳过无效顶点
		if len(vi) < 2 || len(vj) < 2 {
			continue
		}

		xi, yi := vi[0], vi[1]
		xj, yj := vj[0], vj[1]

		// 检查点是否在顶点上
		if (xi == x && yi == y) || (xj == x && yj == y) {
			return true
		}

		// 检查点是否在边上（但不在顶点上）
		if (yi == yj && yi == y) && ((xi <= x && x <= xj) || (xj <= x && x <= xi)) {
			return true
		}

		// 跳过水平边
		if yi == yj {
			continue
		}

		// 检查边是否跨越射线
		if (yi <= y && yj > y) || (yj <= y && yi > y) {
			// 计算交点X坐标
			slope := (xj - xi) / (yj - yi)
			xIntersect := xi + (y-yi)*slope

			// 处理浮点精度
			if math.Abs(x-xIntersect) < 1e-9 {
				return true
			}

			if x < xIntersect {
				inside = !inside
			}
		}
	}
	return inside
}

func PointToWKT(point rtreego.Point) string {
	if len(point) < 2 {
		return "POINT EMPTY"
	}
	return fmt.Sprintf("POINT (%f %f)", point[0], point[1])
}

// 单位百万分之一度
type GaeaPolygonSearchTree struct {
	rtree *rtreego.Rtree
}

func NewSearchTreeByPolygons() *GaeaPolygonSearchTree {
	// 初始化R树（二维）
	rt := rtreego.NewTree(2, 25, 50)
	return &GaeaPolygonSearchTree{
		rtree: rt,
	}
}
func (t *GaeaPolygonSearchTree) InsertWktAsPolygon(id string, name string, wktAsPolygon string, buffer float64) {
	utils.Logger.Info("InsertWktAsPolygon", zap.String("wktAsPolygon", wktAsPolygon))
	polygonObj, err := NewPolygonFromWKT(id, name, wktAsPolygon)
	if err != nil {
		utils.Logger.Error("failed to NewPolygonFromWKB", zap.Error(err))
		return
	}
	t.InsertPolygon(polygonObj)
}

func (t *GaeaPolygonSearchTree) InsertPolygon(polygon *Polygon) {
	t.rtree.Insert(polygon)
}
func (t *GaeaPolygonSearchTree) SearchIntersect(x, y float64, anyOrAll bool) []*Polygon {
	pt := rtreego.Point{x, y}
	tinyRect := pt.ToRect(0.01)
	results := t.rtree.SearchIntersect(tinyRect, func(results []rtreego.Spatial, object rtreego.Spatial) (refuse, abort bool) {
		return false, false
	})

	// 精确判断
	var founds []*Polygon
	for _, obj := range results {
		if poly, ok := obj.(*Polygon); ok {
			if PointInLineString(pt, poly.Vertices) {
				founds = append(founds, poly)
				if anyOrAll {
					break
				}
			}
		}
	}
	return founds
}
