package utils

import (
	"context"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/wkt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// refer to "github.com/twpayne/go-geom/encoding/geojson"
type Geometry struct {
	Srid int
	//Geom geom.T
	Wkt string
}

// gorm 创建表使用
func (g Geometry) GormDataType() string {
	return "geometry"
	//return "Geometry USING geom::Geometry"
}
func (g *Geometry) Scan(input interface{}) error {
	if input == nil {
		g.Wkt = ""
	}
	var gem geom.T
	var err error
	var srid int
	switch v := input.(type) {
	case []byte:
		gem, srid, err = g.Decode4Spatialite(v)
		if err != nil {
			return errors.Unwrap(err)
		}
		g.Wkt, err = wkt.Marshal(gem)
		g.Srid = srid
	case string:
		//postgis
		gem, err = ewkbhex.Decode(v)
		if err != nil {
			return errors.Unwrap(err)
		}
		g.Wkt, err = wkt.Marshal(gem)
		g.Srid = gem.SRID()
	default:
		return errors.New("invalid geometry")
	}

	return nil
}

func (g Geometry) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_GeomFromEWKT(?)",
		Vars: []interface{}{fmt.Sprintf("SRID=%d;%s", g.Srid, g.Wkt)},
	}
}
func (g Geometry) Value() (driver.Value, error) {
	//value, err := ewkbhex.Encode(g.Geom, binary.BigEndian)
	//if err != nil {
	//	GetTargetLogger().Error(err.Error())
	//}
	//return value, err
	var (
		gt  geom.T
		err error
	)
	if gt, err = wkt.Unmarshal(g.Wkt); err != nil {
		Logger.Error("failed to GeometryFromWkt", zap.Error(err))
		return nil, err
	}
	if g.Srid == 0 {
		//如果发现srid=0的就设置成 4326
		g.Srid = 4326
	}
	gt, _ = geom.SetSRID(gt, g.Srid)
	//SRID 不是wkt和wkb的官方标准,是ewkb的标准
	value, err := ewkbhex.Encode(gt, binary.BigEndian)
	if err != nil {
		Logger.Error("geom", zap.Error(err))
	}
	return value, err
}
func (g Geometry) extractWKB(data []byte) ([]byte, error) {
	// 固定头部长度：START(1) + ENDIAN(1) + SRID(4) + MBR(32) = 38 字节
	// 但实际 MBR 是否存在需根据标志位判断（此处假设存在 MBR）
	wkbStartIndex := 38 // 跳过前 38 字节（START + ENDIAN + SRID + MBR）

	// 检查数据长度
	if len(data) < wkbStartIndex+1 {
		return nil, fmt.Errorf("data too short for WKB extraction")
	}

	// 提取 WKB 数据（排除末尾的 0xFE）
	wkbData := data[wkbStartIndex : len(data)-1]
	return wkbData, nil
}
func (g Geometry) Decode4Spatialite(blob []byte) (gem geom.T, srid int, err error) {
	/**
	根据提供的HEX报文和搜索结果，结合SpatiaLite的BLOB-Geometry格式，以下是按字节解析的详细分析：

	---

	### **1. 起始标志与字节序**
	- **字节0**: `00`
	  起始标志（START），表示这是SpatiaLite的几何BLOB数据。
	- **字节1**: `01`
	  ENDIAN标志，表示数据采用**小端模式（Little Endian）**，后续多字节数据需按小端解析。

	---

	### **2. SRID（空间参考系统标识符）**
	- **字节2-5**: `11 0F 00 00`
	  转换为十进制：`0x00000F11 = 3857`
	  表示使用EPSG:3857投影坐标系（Web Mercator），常用于地图服务。

	---

	### **3. MBR（最小外接矩形）**
	- **字节6-13（MBR_MIN_X）**: `B5 09 E7 43 F2 D3 68 41`
	  小端转双精度浮点数：`4168D3F2.43E709B5` → 约为 **-73.99**（经度）。
	- **字节14-21（MBR_MIN_Y）**: `AE 79 8A 8E 6B 29 47 41`
	  解析为约为 **40.71**（纬度）。
	- **字节22-29（MBR_MAX_X）**: `D4 66 28 E2 04 D4 68 41`
	  解析为约为 **-73.98**（经度）。
	- **字节30-37（MBR_MAX_Y）**: `7A B1 F6 3B C1 29 47 41`
	  解析为约为 **40.72**（纬度）。
	  **MBR范围**：覆盖纽约市曼哈顿岛的某区域。

	---

	### **4. 几何类型与结构**
	- **字节38**: `7C`
	  MBR结束标志（MBR_END），固定为0x7C。
	- **字节39-42**: `03 00 00 00`
	  CLASS TYPE为3（POLYGON），表示这是一个多边形几何对象。

	---

	### **5. 几何数据（多边形坐标）**
	- **字节43开始**：`AD 1D 9B 09 F3 D3 68 41 0C C8 DF 98 8F 29 47 41 ...`
	  根据多边形类型，后续数据包含环（Ring）的数量和坐标序列。
	  - **环数量**：通常为4字节整数，假设此处为`00000001`（1个环）。
	  - **坐标序列**：每对双精度浮点数表示一个顶点。例如：
	    - 第一个顶点：`AD1D9B09F3D36841`（X）和 `0CC8DF988F294741`（Y）
	      解析为坐标（-73.989, 40.712）。
	    - 后续顶点需按8字节小端模式逐对解析，最终构成多边形边界。

	---

	### **6. 扩展特性与校验**
	- **压缩几何体**：若CLASS TYPE为1000003（压缩多边形），则数据包含压缩算法标志，但此处未见相关标记。
	- **校验值**：HEX末尾的`FE`可能为Intel HEX文件的校验和，但在此BLOB中更可能是几何数据的一部分，需结合具体结构判断。

	---

	### **总结**
	该HEX报文描述一个基于EPSG:3857坐标系的多边形几何对象，覆盖纽约曼哈顿区域，包含多个坐标点。解析过程需注意：
	1. **小端模式**：所有多字节数据（如SRID、浮点数）需按小端解析。
	2. **MBR有效性**：MBR范围需与实际坐标一致，确保空间查询效率。
	3. **几何结构**：多边形需闭合（首尾点相同），且无自相交，否则可能导致GIS处理错误。

	完整坐标列表可通过逐对提取浮点数获得，适用于地图渲染或空间分析。
	*/
	// ExtractWKBFromSpatiaLite 从SpatiaLite BLOB提取标准WKB格式
	// 参考SpatiaLite BLOB结构定义
	// 使用binary包处理字节序

	// 验证基本结构
	if len(blob) < 39 {
		return nil, 0, fmt.Errorf("blob长度不足，最小需要39字节")
	}
	var order binary.ByteOrder
	byteOrder := blob[1]
	if byteOrder != 0x00 && byteOrder != 0x01 {
		return nil, srid, fmt.Errorf("无效的字节序标志（期待0x00或0x01，实际0x%x）", byteOrder)
	}

	switch byteOrder {
	case 0x00:
		order = binary.BigEndian
	case 0x01:
		order = binary.LittleEndian
	default:
		return nil, 0, fmt.Errorf("invalid byte order: %x", byteOrder)
	}
	if blob[0] != 0x00 {
		return nil, 0, fmt.Errorf("无效的SpatiaLite起始标志（期待0x00，实际0x%x）", blob[0])
	}
	if blob[38] != 0x7C {
		return nil, 0, fmt.Errorf("无效的MBR结束标志（期待0x7C，实际0x%x）", blob[38])
	}

	srid = int(order.Uint32(blob[2:6]))
	// 获取字节序标志

	// 构造标准WKB（字节序标志 + WKB核心数据）
	wkb := make([]byte, 1+len(blob[39:]))
	wkb[0] = byteOrder
	copy(wkb[1:], blob[39:])
	gem, err = ewkb.Unmarshal(wkb)
	return gem, srid, err
}
