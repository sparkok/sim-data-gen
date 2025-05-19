package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

func RefInt(value int) *int {
	return &value
}
func RefInt64(value uint64) *uint64 {
	return &value
}
func RefFloat(value float64) *float64 {
	return &value
}
func RefTime(value time.Time) *time.Time {
	return &value
}
func RefString(value string) *string {
	return &value
}
func RefBool(value bool) *bool {
	return &value
}
func RefFloat64(value float64) *float64 {
	return &value
}
func StringPtrAsTxt(value *string) string {
	if value == nil {
		return "nil"
	}
	return *value
}
func EverOk(value time.Time, _ error) time.Time {
	return value
}
func AssertNotNull(ptr interface{}, message string) interface{} {
	if ptr == nil {
		Logger.Error(message)
		return ptr
	}
	return ptr
}
func RefBool2Int(bootPtr *bool) *int {
	if bootPtr == nil {
		return nil
	}
	if *bootPtr == true {
		return RefInt(1)
	} else {
		return RefInt(0)
	}

}
func Utc() int {
	var utc = time.Now().Unix()
	return int(utc)
}

var location *time.Location

func ToBeijing(value time.Time) time.Time {
	if location == nil {
		//location, _ = time.LoadLocation("Asia/Shanghai")
		location, _ = time.LoadLocation("UTC")
	}

	// 加载目标时区
	// 将 UTC 时间转换为目标时区的时间
	return value.In(location)
}

// 系统使用的时区,这里设置为中国的东八区
var SystemZone = time.FixedZone("CST", 8*60*60)

func Change2TheDay(theTime time.Time, theDay string) (rst time.Time) {
	rst, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %02d:%02d:%02d", theDay, theTime.Hour(), theTime.Minute(), theTime.Second()), SystemZone)
	return rst
}

func Obj2JsonTxt(obj interface{}) string {
	bin, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	var res = string(bin)
	return res
}
func FloatPtrValue(ValuePtr *float64, defaultValue float64) float64 {
	if ValuePtr == nil {
		return defaultValue
	}
	return *ValuePtr
}
func CopyRefInt(value *int) *int {
	if value == nil {
		return nil
	}
	tmpValue := *value
	return &tmpValue
}
func CopyRefInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	tmpValue := *value
	return &tmpValue
}
func CopyRefFloat(value *float64) *float64 {
	if value == nil {
		return nil
	}
	tmpValue := *value
	return &tmpValue
}
func CopyRefString(value *string) *string {
	if value == nil {
		return nil
	}
	tmpValue := *value
	return &tmpValue
}
func CopyRefBool(value *bool) *bool {
	if value == nil {
		return nil
	}
	tmpValue := *value
	return &tmpValue
}
func CopyRefFloat64(value *float64) *float64 {
	if value == nil {
		return nil
	}
	tmpValue := *value
	return &tmpValue
}
