package services

import (
	. "sim_data_gen/utils"
	"fmt"
	"github.com/twpayne/go-geom"
	ProjLib "github.com/xeonx/proj4"
	"go.uber.org/zap"
	"math"
	"sync"
)

// EPSGCode is the enum type for coordinate systems
type EPSGCode int
type ProjService struct {
}

// 库的地址在 https://github.com/xeonx/proj4/blob/master/geocent.c
const (
	Deg2Rad                              = math.Pi / 180.0
	Rad2Deg                              = 180.0 / math.Pi
	EPSG3395                    EPSGCode = 3395
	EPSG3857                             = 3857
	EPSG4527                             = 4527 // ground xy in changzhou
	EPSG4087                             = 4087
	EPSG4326                             = 4326
	WebMercator                          = EPSG3857
	WorldEquidistantCylindrical          = EPSG4087
	WorldMercator                        = EPSG3395
	WGS84                                = EPSG4326
	LocalEPSG                            = 9999
)

func Local2MercatorEarth(x, y float64) (float64, float64, error) {
	earth, err := Transform2D(LocalEPSG, WebMercator, [3]float64{x, y, 0})
	if err != nil {
		return x, y, err
	}
	return earth[0], earth[1], nil
}
func LngLat2Earth(input [2]float64) ([2]float64, error) {
	earth, err := Transform2D(EPSG4326, WebMercator, [3]float64{input[0], input[1], 0})
	result := [2]float64{earth[0], earth[1]}
	return result, err
}
func Earth2LngLat(input [2]float64) ([2]float64, error) {
	lnglat, err := Transform2D(WebMercator, EPSG4326, [3]float64{input[0], input[1], 0})
	result := [2]float64{lnglat[0], lnglat[1]}
	return result, err
}

func Transform2DCoords(fromSrid int, toSrid int, input [2]float64) ([2]float64, error) {
	lnglat, err := Transform2D(EPSGCode(fromSrid), EPSGCode(toSrid), [3]float64{input[0], input[1], 0})
	result := [2]float64{lnglat[0], lnglat[1]}
	return result, err
}

// ensure only one person is updating our cache of ProjObjs at a time
var cacheLock = sync.Mutex{}

func Transform2D(from EPSGCode, dest EPSGCode, input [3]float64) (result [3]float64, err error) {
	var (
		projSrc        *ProjKit
		projDst        *ProjKit
		transformation ProjLib.Transformation
		x              = []float64{input[0]}
		y              = []float64{input[1]}
		z              = []float64{input[2]} //height - meter for IsGeoCent
	)
	//result = [3]float64{0, 0, 0}
	defer cacheLock.Unlock()
	cacheLock.Lock()
	projSrc, err = newProjKit(from)
	if err != nil {
		return
	}
	projDst, err = newProjKit(dest)
	if err != nil {
		return
	}
	//output, err = utils.TransformCgcs2000ToEsp4326(input)
	if projSrc.IsLatLong {
		x[0] = x[0] * Deg2Rad
		y[0] = y[0] * Deg2Rad
	}
	//参考  https://en.wikipedia.org/wiki/Geocentric_model
	if projSrc.IsGeoCent {
		Logger.Error("TO DO,result maybe is error")
		//x[0] = x[0] * Deg2Rad
		//y[0] = y[0] * Deg2Rad
	}

	transformation, err = ProjLib.NewTransformation(projSrc.ProjObj, projDst.ProjObj)
	if err != nil {
		return
	}

	err = transformation.TransformRaw(x, y, z)
	if err != nil {
		Logger.Error(fmt.Sprintf("failed to TransformRaw(%f, %f, %f)", x, y, z))
		return
	}
	if projDst.IsLatLong {
		x[0] = x[0] * Rad2Deg
		y[0] = y[0] * Rad2Deg
	}
	if projDst.IsGeoCent {
		Logger.Error("TO DO,result maybe is error")
		//x[0] = x[0] * Deg2Rad
		//y[0] = y[0] * Deg2Rad
	}
	//Logger.Info(fmt.Sprintf("transform [ESPG:%d] %f,%f,%f => [ESPG:%d] %f,%f,%f", projSrc.EspgCode, input[0], input[1], input[2], projDst.EspgCode, x[0], y[0], z[0]))
	result[0] = x[0]
	result[1] = y[0]
	result[2] = z[0]
	return
}

//---------------------------------------------------------------------------

// conversion holds the objects needed to perform a conversion
type ProjKit struct {
	EspgCode  EPSGCode
	ProjObj   *ProjLib.Proj
	IsLatLong bool
	IsGeoCent bool
}

var ProjKits = map[EPSGCode]*ProjKit{}

// 参数含义的中文说明 https://www.cnblogs.com/eshinex/p/10299947.html#
/**
+a	 椭球体长半轴长度
+axis	 轴防线
+b	 椭球体短半轴长度
+ellps	 椭球体名称，在cmd中输入：proj -le 查看支持哪些椭球体
+k	 比例系数（比例因子），旧版本，不赞成使用
+k_0	 比例系数（比例因子）
+lat_0	 维度起点
+lon_0	 中央经线
+lon_wrap	 中央经线的包装参数（详见下面的说明）
+no_defs
不要使用proj库中的缺省定义文件。
在linux中路径为：/usr/share/proj/proj_def.dat
windows中为安装路径下的：E:\SvnWorkspace\LY_WEB_GIS\branches\Documents\ms4w-mapserver-for-wimdows\release-1911-x64-gdal-2-3-3-mapserver-7-2-1\bin\proj\SHARE\proj_def.dat
标红处是我安装MapServer是自带安装Proj的路径
+over	 允许经度输出在-180到180范围之外，禁用wrapping(详见下面的说明)
+pm	 备用本初子午线(通常是一个城市名称，见下文)
+proj	 投影名称，在cmd中输入：proj -l 查看数据支持
+units	 水平单位，meters（米）、 US survey feet, etc（英尺等 us-ft）.
+vunits	 垂直单位
+x_0	 东（伪）偏移量
+y_0	 北（伪）偏移量
*/
var projStrings = map[EPSGCode]string{
	EPSG4326: "+proj=longlat +datum=WGS84 +no_defs +type=crs",
	EPSG3395: "+proj=merc +lon_0=0 +k=1 +x_0=0 +y_0=0 +datum=WGS84",                                                // TODO: support +units=m +no_defs
	EPSG3857: "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0",                     // TODO: support +units=m +nadgrids=@null +wktext +no_defs
	EPSG4087: "+proj=eqc +lat_ts=0 +lat_0=0 +lon_0=0 +x_0=0 +y_0=0 +datum=WGS84",                                   // TODO: support +units=m +no_defs
	EPSG4527: "+proj=tmerc +lat_0=0 +lon_0=117 +k=1 +x_0=39500000 +y_0=0 +ellps=GRS80 +units=m +no_defs +type=crs", //refer to //https://epsg.io/4527,transform([39643479.035, 3208227.221],"EPSG:4527","EPSG:4326") => 118.47216855985951,28.982538274117232
}

// newConversion creates a conversion object for the destination systems. If
// such a conversion already exists in the cache, use that.
func newProjKit(espgCode EPSGCode) (*ProjKit, error) {
	str, ok := projStrings[espgCode]
	if !ok {
		return nil, fmt.Errorf("epsg code is not a supported projection")
	}

	projKit, ok := ProjKits[espgCode]
	if ok {
		return projKit, nil
	}

	// need to build it
	var err error
	var projObj *ProjLib.Proj

	projObj, err = ProjLib.InitPlus(str)
	if err != nil {
		return nil, err
	}

	projKit = &ProjKit{
		EspgCode:  espgCode,
		ProjObj:   projObj,
		IsLatLong: projObj.IsLatLong(),
		IsGeoCent: projObj.IsGeoCent(),
	}

	// cache it
	ProjKits[espgCode] = projKit

	return projKit, nil
}

var projService *ProjService

func InitProjKits() {
	for code, _ := range projStrings {
		newProjKit(code)
	}
}
func AddProjString(code EPSGCode, projString string) {
	projStrings[code] = projString
}
func CreateProjService(varNamePrex string) {
	if projService != nil {
		return
	}
	projService = new(ProjService)
	var localProjString string
	if localProjString = GetConfig().String(fmt.Sprintf("%s.%s", varNamePrex, "LocalEPSG.ProjString"), "+proj=tmerc +lat_0=0 +lon_0=117 +k=1 +x_0=39500000 +y_0=0 +ellps=GRS80 +units=m +no_defs +type=crs"); localProjString == "" {
		Logger.Error(fmt.Sprintf("please set %s.url to  the voronoi service", varNamePrex))
		return
	}
	AddProjString(LocalEPSG, localProjString)
	InitProjKits()
}
func ExtractCoordsFromPolygonForFloats(polygon *geom.Polygon, srcSrid int) (coords [][][]float64) {

	ret := make([][][]float64, len(polygon.Coords()), len(polygon.Coords()))
	var indexOfPoint int
	for i, coord := range polygon.Coords() {
		xys := make([][]float64, len(coord), len(coord))
		ret[i] = xys
		for j, xy := range coord {
			if isRightHand4Geom(coord) {
				indexOfPoint = j
			} else {
				indexOfPoint = len(xys) - j - 1
			}
			xys[indexOfPoint] = make([]float64, 2, 2)
			if srcSrid != 4326 && srcSrid != 0 {
				//如果没有包含srid(srcSrid = 0)则认为是4326
				src := [2]float64{xy.X(), xy.Y()}
				newXY, err := Transform2DCoords(srcSrid, 4326, src)
				if err == nil {
					xys[indexOfPoint][0] = newXY[0]
					xys[indexOfPoint][1] = newXY[1]
				} else {
					Logger.Error("convert srcSrid error", zap.Int("srcSrid", srcSrid), zap.Float64("x", xy.X()), zap.Float64("y", xy.Y()))
					xys[indexOfPoint][0] = 0
					xys[indexOfPoint][1] = 0
				}

			} else {
				xys[indexOfPoint][0] = xy.X()
				xys[indexOfPoint][1] = xy.Y()
			}
		}
	}
	return ret
}
func ExtractCoordsFromLineStringForFloats(srcCoords [][]float64, srcSrid int) (coords [][]float64) {
	coords = make([][]float64, len(srcCoords), len(srcCoords))
	var indexOfPoint int
	for j, xy := range srcCoords {
		if isRightHand(srcCoords) {
			indexOfPoint = j
		} else {
			indexOfPoint = len(srcCoords) - j - 1
		}
		coords[indexOfPoint] = make([]float64, 2, 2)
		if srcSrid != 4326 && srcSrid != 0 {
			//如果没有包含srid(srcSrid = 0)则认为是4326
			src := [2]float64{xy[0], xy[1]}
			newXY, err := Transform2DCoords(srcSrid, 4326, src)
			if err == nil {
				coords[indexOfPoint][0] = newXY[0]
				coords[indexOfPoint][1] = newXY[1]
			} else {
				Logger.Error("convert srcSrid error", zap.Int("srcSrid", srcSrid), zap.Float64("x", xy[0]), zap.Float64("y", xy[1]))
				coords[indexOfPoint][0] = 0
				coords[indexOfPoint][1] = 0
			}

		} else {
			coords[indexOfPoint][0] = xy[0]
			coords[indexOfPoint][1] = xy[1]
		}
	}
	return coords
}

func isRightHand4Geom(coords []geom.Coord) bool {
	coord0 := coords[0]
	coord1 := coords[1]
	if coord0.X() <= coord1.X() {
		if coord0.Y() <= coord1.Y() {
			return true
		} else {
			return false
		}
	} else {
		if coord0.Y() <= coord1.Y() {
			return false
		} else {
			return true
		}
	}
	return false
}

func isRightHand(coords [][]float64) bool {
	coord0 := coords[0]
	coord1 := coords[1]
	if coord0[0] <= coord1[0] {
		if coord0[1] <= coord1[1] {
			return true
		} else {
			return false
		}
	} else {
		if coord0[1] <= coord1[1] {
			return false
		} else {
			return true
		}
	}
	return false
}
