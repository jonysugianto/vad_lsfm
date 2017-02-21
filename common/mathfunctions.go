package common

import (
	"math"
)

var DELTA_DISTANCE int = 1

func Abs(data []float64) []float64 {
	var size = len(data)
	var ret []float64
	for i := 0; i < size; i++ {
		ret = append(ret, math.Abs(data[i]))
	}
	return ret
}

func LogNatural(data []float64) []float64 {
	var size = len(data)
	var ret []float64
	for i := 0; i < size; i++ {
		ret = append(ret, math.Log(data[i]))
	}
	return ret
}

func Log2(data []float64) []float64 {
	var size = len(data)
	var ret []float64
	for i := 0; i < size; i++ {
		ret = append(ret, math.Log2(data[i]))
	}
	return ret
}

func LogPositiv(data []float64) []float64 {
	var size = len(data)
	var ret []float64
	for i := 0; i < size; i++ {
		ret = append(ret, math.Log1p(data[i]))
	}
	return ret
}

func DeltaSingleValue(data []float64) []float64 {
	var ret []float64
	var size = len(data)
	for i := 0; i < DELTA_DISTANCE; i++ {
		ret = append(ret, 0)
	}

	for i := DELTA_DISTANCE; i < size; i++ {
		ret = append(ret, data[i]-data[i-DELTA_DISTANCE])
	}

	return ret
}
