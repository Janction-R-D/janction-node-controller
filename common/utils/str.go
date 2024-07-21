package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func StrToBytes(s string) []byte {
	if s == "" {
		return []byte{}
	}
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func BytesToStr(b []byte) string {
	if b == nil || len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

func IntToCN(num int) string {
	// 1、数字为0
	if num == 0 {
		return "零"
	}
	var ans string
	// 数字
	szdw := []string{"十", "百", "千", "万", "十万", "百万", "千万", "亿"}
	// 数字单位
	sz := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	res := make([]string, 0)

	// 数字单位角标
	idx := -1
	for num > 0 {
		// 当前位数的值
		x := num % 10
		// 2、数字大于等于10
		// 插入数字单位，只有当数字单位角标在范围内，且当前数字不为0 时才有效
		if idx >= 0 && idx < len(szdw) && x != 0 {
			res = append([]string{szdw[idx]}, res...)
		}
		// 3、数字中间有多个0
		// 当前数字为0，且后一位也为0 时，为避免重复删除一个零文字
		if x == 0 && len(res) != 0 && res[0] == "零" {
			res = res[1:]
		}
		// 插入数字文字
		res = append([]string{sz[x]}, res...)
		num /= 10
		idx++
	}
	// 4、个位数为0
	if len(res) > 1 && res[len(res)-1] == "零" {
		res = res[:len(res)-1]
	}
	// 合并字符串
	for i := 0; i < len(res); i++ {
		ans = ans + res[i]
	}
	return ans
}

func BytesIndexToRuneIndex(str string, start, end int) (left, right int) {
	left = len([]rune(str[:start]))
	right = len([]rune(str[:end]))
	return left, right
}

func ToString(i interface{}) (v string) {
	switch t := i.(type) {
	case string:
		v = t
	case int:
		v = strconv.Itoa(t)
	case int8:
		v = strconv.Itoa(int(t))
	case int16:
		v = strconv.Itoa(int(t))
	case int32:
		v = strconv.Itoa(int(t))
	case int64:
		v = strconv.Itoa(int(t))
	case uint:
		v = strconv.Itoa(int(t))
	case uint8:
		v = strconv.Itoa(int(t))
	case uint16:
		v = strconv.Itoa(int(t))
	case uint32:
		v = strconv.Itoa(int(t))
	case uint64:
		v = strconv.Itoa(int(t))
	case float32, float64:
		v = fmt.Sprintf("%v", t)
	case []uint8:
		v = Ui8ToA(t)
	}

	return v
}

// []uint8 转字符串
func Ui8ToA(i interface{}) string {
	if v, ok := i.(string); ok {
		return v
	}

	return string(Ui8ToB(i))
}

// []uint8 转字符串字节
func Ui8ToB(i interface{}) (b []byte) {
	if v, ok := i.([]uint8); ok {
		for _, val := range v {
			b = append(b, val)
		}
	}

	return b
}

func JoinAndWrapper(array []string, sep, wrap string) string {
	var arrayCap int
	for _, s := range array {
		arrayCap += len(s) + len(sep) + len(wrap)*2
	}

	var buf = bytes.Buffer{}
	buf.Grow(arrayCap)

	for _, s := range array {
		buf.WriteString(wrap)
		buf.WriteString(s)
		buf.WriteString(wrap)
		buf.WriteString(sep)
	}
	return buf.String()[:buf.Len()-len(sep)]
}

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
	charStr = letters + digits
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charStr[rand.Intn(len(charStr))]
	}
	return string(result)
}

func StrToIntIgnoreErr(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}

func CheckCSVStr(s string) bool {
	invalidPrefixes := []string{"=", "+", "-", "@"}
	for _, prefix := range invalidPrefixes {
		if strings.HasPrefix(s, prefix) {
			return false
		}
	}

	// Check for tab and carriage return
	if len(s) > 0 && (s[0] == 0x09 || s[0] == 0x0D) {
		return false
	}

	return true
}
