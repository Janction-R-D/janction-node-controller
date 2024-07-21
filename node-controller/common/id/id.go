package id

import (
	"encoding/json"
	"strconv"
)

type Type int8

const (
	Type_Zero Type = iota
	Type_String
	Type_Int
)

type ID struct {
	val interface{}
}

// FromInt creates an ID object with an int64 value.
func FromInt(val int64) ID {
	return ID{
		val: val,
	}
}

// FromString creates an ID object with a Type_String value.
func FromString(val string) ID {
	return ID{
		val: val,
	}
}

// Parse the given Type_String and try to convert it to an integer before
// setting it as a Type_String value.
func Parse(val string) ID {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return FromString(val)
	}
	return FromInt(i)
}

// Type return type.
func (id *ID) Type() Type {
	if id == nil || id.val == nil {
		return Type_Zero
	}
	if _, ok := id.val.(int64); ok {
		return Type_Int
	}
	return Type_String
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *ID) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		var str string
		err := json.Unmarshal(value, &str)
		if err == nil {
			id.val = str
		}
		return err
	}
	var intVal int64
	err := json.Unmarshal(value, &intVal)
	if err == nil {
		id.val = intVal
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (x ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.val)
}

// Type_String returns the Type_String value, or the Itoa of the int value.
func (id ID) String() string {
	if id.val == nil {
		return ""
	}

	if v, ok := id.val.(string); ok {
		return v
	}
	return strconv.FormatInt(id.val.(int64), 10)
}

// IntValue returns the IntVal if type Int, or if
// it is a Type_String, will attempt a conversion to int,
// returning 0 if a parsing error occurs.
func (id ID) IntValue() int64 {
	if v, ok := id.val.(int64); ok {
		return v
	}
	if v, ok := id.val.(string); ok {
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	}
	return 0
}

// IsZero return true if is not empty.
func (id ID) IsZero() bool {
	if id.val == nil {
		return true
	}
	if v, ok := id.val.(string); ok {
		return v == ""
	}
	if v, ok := id.val.(int64); ok {
		return v == 0
	}

	return true
}
