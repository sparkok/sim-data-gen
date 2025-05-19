package utils

import "strings"

func GetMapIntItemByPath(obj interface{}, path string, defaultValue int) int {
	item := GetMapItemByPath(obj, path)
	if item == nil {
		return defaultValue
	}
	return item.(int)
}
func GetMapFloatItemByPath(obj interface{}, path string, defaultValue float64) float64 {
	item := GetMapItemByPath(obj, path)
	if item == nil {
		return defaultValue
	}
	return item.(float64)
}
func GetMapStringItemByPath(obj interface{}, path string, defaultValue string) string {
	item := GetMapItemByPath(obj, path)
	if item == nil {
		return defaultValue
	}
	return item.(string)
}
func GetMapArrayItemByPath(obj interface{}, path string, defaultValue []interface{}) []interface{} {
	item := GetMapItemByPath(obj, path)
	if item == nil {
		return defaultValue
	}
	return item.([]interface{})
}

func GetMapItemByPath(obj interface{}, path string) interface{} {
	names := strings.Split(path, ".")
	var (
		section interface{}
		existed bool
	)
	section = obj
	for _, name := range names {
		switch section.(type) {
		case map[string]interface{}:
			if section, existed = section.(map[string]interface{})[name]; !existed {
				return nil
			}
		default:
			return nil
		}
	}
	return section
}
