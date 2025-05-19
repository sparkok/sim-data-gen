package entity

import "math"

func Almost(realValue float64, expectValue float64) bool {
	return (math.Abs(realValue-expectValue) < 0.00001)
}

func AlmostOne(value float64) bool {
	return Almost(value, 1)
}

func AlmostZero(value float64) bool {
	return Almost(value, 0)
}
