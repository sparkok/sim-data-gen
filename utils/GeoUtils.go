package utils

import "math"

func CalculateDistanceAsMeter(lon1Micro, lat1Micro, lon2Micro, lat2Micro float64) float64 {
	// 单位转换：百万分之一度 → 普通度数
	lon1 := float64(lon1Micro) / 1e6
	lat1 := float64(lat1Micro) / 1e6
	lon2 := float64(lon2Micro) / 1e6
	lat2 := float64(lat2Micro) / 1e6

	// 转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// 计算差值
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine公式计算
	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 地球半径（单位：米）
	earthRadius := 6371.0 * 1000

	return earthRadius * c
}
