package utils

//is a point with float values
type FloatPoint struct {
	X float64
	Y float64
}

func Max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
