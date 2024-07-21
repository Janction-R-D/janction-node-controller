package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func JoinNums(nums []int, split string) string {
	if len(nums) == 0 {
		return ""
	}

	strSlice := make([]string, len(nums))
	for i, n := range nums {
		strSlice[i] = fmt.Sprint(n)
	}

	return strings.Join(strSlice, split)
}

func CheckParam(params string) bool {
	reg, _ := regexp.Compile("^[\u4E00-\u9FA5a-zA-Z0-9\\-_/\\s\\.]+$")
	if reg.MatchString(params) {
		return true
	}
	return false
}

func Indent(src, prefix, indent string) (dst string, err error) {
	var str bytes.Buffer
	err = json.Indent(&str, []byte(src), prefix, indent)
	if err != nil {
		return "", err
	}
	return str.String(), nil
}
