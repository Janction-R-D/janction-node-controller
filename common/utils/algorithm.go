package utils

type compareType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | uintptr | string
}

func Max[T compareType](a, b T) T {
	if a < b {
		return b
	}
	return a
}

func Min[T compareType](a, b T) T {
	if a < b {
		return a
	}
	return b
}
