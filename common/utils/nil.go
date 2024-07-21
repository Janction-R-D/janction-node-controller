package utils

import (
	"reflect"
)

func NilBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func NilPtrWithDefault[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

func GetDefaultIfZero[T any](v interface{}, def T) (value T) {
	refV := reflect.ValueOf(v)
	if refV.IsZero() {
		return def
	}
	return v.(T)
}
