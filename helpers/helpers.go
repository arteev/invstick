package helpers

import "math"

func MkSlice(args ...interface{}) []interface{} {
	return args
}

func Mul(a, b int) int {
	return a * b
}

func Add(a, b int) int {
	return a + b
}

func Calcpages(onpage, count int) int {
	if onpage == 0 {
		return 0
	}
	i := count / onpage
	if math.Mod(float64(count), float64(onpage)) > 0 {
		i++
	}
	return i
}

func MkSliceRange(from, count int) (result []int) {
	for i := from; i < from+count; i++ {
		result = append(result, i)
	}
	return
}
