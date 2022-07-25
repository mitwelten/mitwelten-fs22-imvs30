package utils

//is a point with float values
type FloatPoint struct {
	X float64
	Y float64
}

type Tuple[T any] struct {
	T1 T
	T2 T
}

func Max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func Abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func Min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}
