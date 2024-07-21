package utils

import (
	"math"
	"strconv"
)

func RoundFloat(f float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(f*pow) / pow
}

// 任意类型转int
func ToInt(i interface{}) (v int) {
	switch t := i.(type) {
	case int:
		v = t
	case int8:
		v = int(t)
	case int16:
		v = int(t)
	case int32:
		v = int(t)
	case int64:
		v = int(t)
	case uint:
		v = int(t)
	case uint8:
		v = int(t)
	case uint16:
		v = int(t)
	case uint32:
		v = int(t)
	case uint64:
		v = int(t)
	case float64:
		v = int(t)
	case string:
		vv, _ := strconv.ParseInt(t, 10, 64)
		v = int(vv)
	case []uint8:
		vv, _ := strconv.ParseInt(Ui8ToA(i), 10, 64)
		v = int(vv)
	}

	return v
}
