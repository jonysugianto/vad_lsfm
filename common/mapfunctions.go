package common

import (

)

func MapArray(input [][]float64, mf func([]float64)[]float64) [][]float64 {
	var size = len(input)
	var result [][]float64
	for i := 0; i < size; i++ {
		var temp = mf(input[i])
		result = append(result, temp)
	}
	return result
}

func MapArray2Single(input [][]float64, mf func([]float64)float64) []float64 {
	var size = len(input)
	var result []float64
	for i := 0; i < size; i++ {
		var temp = mf(input[i])
		result = append(result, temp)
	}
	return result
}

func MapValue(input []float64, mf func(float64)float64) []float64 {
	var size = len(input)
	var result []float64
	for i := 0; i < size; i++ {
		var temp = mf(input[i])
		result = append(result, temp)
	}
	return result
}

func MapArray2Value(input [][]float64, mf func([]float64)float64) []float64 {
	var size = len(input)
	var result []float64
	for i := 0; i < size; i++ {
		var temp = mf(input[i])
		result = append(result, temp)
	}
	return result
}
