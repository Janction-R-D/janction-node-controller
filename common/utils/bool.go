package utils

func ConvertBoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ReturnByBool(ok bool, trueDesc, falseDesc string) string {
	if ok {
		return trueDesc
	}
	return falseDesc
}

func ReturnIfOk[T any](ok bool, okVal, defaultVal T) (v T) {
	if ok {
		return okVal
	}
	return defaultVal
}

func ReturnPtrIfOk[T any](ok bool, okVal, defaultVal *T) (v *T) {
	if ok {
		return okVal
	}
	return defaultVal
}

func ConvertIntToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}
