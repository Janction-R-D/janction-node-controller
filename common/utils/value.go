package utils

func TypeValue[T any](value *T) T {
	if value == nil {
		var v T
		return v
	}
	return *value
}

func TypePtr[T any](value T) *T {
	return &value
}

func TypeSliceValue[T any](a []*T) []T {
	if a == nil {
		return nil
	}

	res := make([]T, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] != nil {
			res[i] = *a[i]
		}
	}
	return res
}

func TypeSlicePtr[T any](a []T) []*T {
	if len(a) == 0 {
		return []*T{}
	}

	res := make([]*T, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = &a[i]
	}
	return res
}

func InterfaceValue[T any](value interface{}) T {
	v, ok := value.(T)
	if !ok {
		var vv T
		return vv
	}
	return v
}
