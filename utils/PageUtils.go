package utils

import "fmt"
/**
	返回limit和offset的sql子句
 */
func MakeLimitOffset(curPage int,pageSize int)string {
	return fmt.Sprintf(" limit %d offset %d",pageSize,pageSize * (curPage - 1))
}
