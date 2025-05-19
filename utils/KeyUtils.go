package utils

import (
	"strings"
)
import "strconv"
import "github.com/google/uuid"
func SplitKey(value string)[]string  {
	return strings.Split(value,"--")
}
func UnionKey(keys ... interface{})string  {
	txt := ""
	for key:= range keys {
		txt += "--" + string(key)
	}
	return txt
}
func String2Int(value string) int  {
	v64,_ := strconv.ParseInt(value,10,64)
	return  int(v64)
}
func StringPtr2IntPtr(value *string) *int  {
	if(value == nil) {
		return nil
	}
	v64,_ := strconv.Atoi(*value)
	return  &v64
}

func String2Float(value string) float64  {
	v,_ := strconv.ParseFloat(value,64)
	return v
}
/**
生成uuid
 */
func Uuid() *string {
	uuid,_ := uuid.NewUUID()
	uuidTxt := uuid.String()

	return &uuidTxt
}
func AutoInt() *int {
	//value := 0;
	return nil
}
