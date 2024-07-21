package utils

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

func LessByNumAlphaPinYin(a, b string) bool {
	a1, _ := UTF82GBK(a)
	b1, _ := UTF82GBK(b)
	b1Len := len(b1)
	for idx, chr := range a1 {
		if idx > b1Len-1 {
			return false
		}
		if chr != b1[idx] {
			return chr < b1[idx]
		}
	}
	return a < b
}

// UTF82GBK : transform UTF8 rune into GBK byte array
func UTF82GBK(src string) ([]byte, error) {
	GB18030 := simplifiedchinese.All[0]
	return io.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), GB18030.NewEncoder()))
}

// GBK2UTF8 : transform  GBK byte array into UTF8 string
func GBK2UTF8(src []byte) (string, error) {
	GB18030 := simplifiedchinese.All[0]
	data, err := io.ReadAll(transform.NewReader(bytes.NewReader(src), GB18030.NewDecoder()))
	return string(data), err
}

func FastSort(arr []int) []int {
	if len(arr) == 0 || len(arr) == 1 {
		return arr
	}

	sign := arr[0]
	left := []int{}
	right := []int{}
	for _, v := range arr {
		if v <= sign {
			left = append(left, v)
			continue
		}
		right = append(right, v)
	}

	return append(append(FastSort(left[1:]), sign), FastSort(right)...)
}
