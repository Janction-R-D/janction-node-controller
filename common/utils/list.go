package utils

func ExistInList(ele interface{}, slc interface{}) bool {
	if slc == nil {
		return false
	}
	switch ele.(type) {
	case string:
		e := ele.(string)
		s := slc.([]string)
		for _, v := range s {
			if v == e {
				return true
			}
		}
		return false
	case int:
		e := ele.(int)
		s := slc.([]int)
		for _, v := range s {
			if v == e {
				return true
			}
		}
		return false
	case int32:
		e := ele.(int32)
		s := slc.([]int32)
		for _, v := range s {
			if v == e {
				return true
			}
		}
		return false
	case int64:
		e := ele.(int64)
		s := slc.([]int64)
		for _, v := range s {
			if v == e {
				return true
			}
		}
		return false
	case uint32:
		e := ele.(uint32)
		s := slc.([]uint32)
		for _, v := range s {
			if v == e {
				return true
			}
		}
		return false
	default:
		return false
	}
}
