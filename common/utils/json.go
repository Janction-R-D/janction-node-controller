package utils

import (
	"bytes"
	"encoding/json"
)

func MarshalIgnoreErr(v interface{}) string {
	if v == nil {
		return ""
	}
	data, _ := json.Marshal(v)
	return string(data)
}

func UnmarshalIgnoreErr(src string, dst interface{}) {
	_ = json.Unmarshal(StrToBytes(src), dst)
}

func UnmarshalToInt32Slice(v string) (res []int32, err error) {
	if v == "" || v == "[]" {
		return nil, nil
	}
	err = json.Unmarshal([]byte(v), &res)
	return res, err
}

func UnmarshalToUint32Slice(v string) (res []uint32, err error) {
	if v == "" || v == "[]" {
		return nil, nil
	}
	err = json.Unmarshal([]byte(v), &res)
	return res, err
}

func UnmarshalWithoutErr[T any](value string) T {
	var v T
	if value == "" {
		return v
	}
	_ = json.Unmarshal(StrToBytes(value), &v)
	return v
}

func Unmarshal[T any](value string) (v T, err error) {
	if value == "" {
		return v, nil
	}
	err = json.Unmarshal(StrToBytes(value), &v)
	return v, err
}

func Indent(src, prefix, indent string) (dst string, err error) {
	var str bytes.Buffer
	err = json.Indent(&str, []byte(src), prefix, indent)
	if err != nil {
		return "", err
	}
	return str.String(), nil
}

func UnmarshalIgnoreEmpty(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal([]byte(data), v)
}
