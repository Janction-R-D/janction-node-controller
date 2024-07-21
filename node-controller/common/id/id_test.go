package id

import (
	"encoding/json"
	"reflect"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestFromInt(t *testing.T) {
	i := FromInt(93)
	if i.val != int64(93) {
		t.Errorf("Expected IntVal=93, got %+v", i)
	}
}

func TestFromString(t *testing.T) {
	i := FromString("76")
	if i.val != "76" {
		t.Errorf("Expected StrVal=\"76\", got %+v", i)
	}
}

type IDHolder struct {
	IOrS ID `json:"val"`
}

func TestIntStrUnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		result ID
	}{
		{"{\"val\": 123}", FromInt(123)},
		{"{\"val\": \"123\"}", FromString("123")},
	}

	for _, c := range cases {
		var result IDHolder
		if err := json.Unmarshal([]byte(c.input), &result); err != nil {
			t.Errorf("Failed to unmarshal input '%v': %v", c.input, err)
		}
		if result.IOrS.String() != c.result.String() {
			t.Errorf("Failed to unmarshal input '%v': expected %+v, got %+v", c.input, c.result, result)
		}
	}
}

func TestIntStrMarshalJSON(t *testing.T) {
	cases := []struct {
		input  ID
		result string
	}{
		{FromInt(123), "{\"val\":123}"},
		{FromString("123"), "{\"val\":\"123\"}"},
	}

	for _, c := range cases {
		input := IDHolder{c.input}
		result, err := json.Marshal(&input)
		if err != nil {
			t.Errorf("Failed to marshal input '%v': %v", input, err)
		}
		if string(result) != c.result {
			t.Errorf("Failed to marshal input '%v': expected: %+v, got %q", input, c.result, string(result))
		}
	}
}

func TestIntStrMarshalJSONUnmarshalYAML(t *testing.T) {
	cases := []struct {
		input ID
	}{
		{FromInt(123)},
		{FromString("123")},
	}

	for _, c := range cases {
		input := IDHolder{c.input}
		jsonMarshalled, err := json.Marshal(&input)
		if err != nil {
			t.Errorf("1: Failed to marshal input: '%v': %v", input, err)
		}

		var result IDHolder
		err = yaml.Unmarshal(jsonMarshalled, &result)
		if err != nil {
			t.Errorf("2: Failed to unmarshal '%+v': %v", string(jsonMarshalled), err)
		}

		if !reflect.DeepEqual(input, result) {
			t.Errorf("3: Failed to marshal input '%+v': got %+v", input, result)
		}
	}
}

func TestID_IsZero(t *testing.T) {
	if !FromInt(0).IsZero() {
		t.Error("zero int should is zero")
	}
	if !FromString("").IsZero() {
		t.Error("empty string should is zero")
	}
	if !(ID{}).IsZero() {
		t.Error("empty should is zero")
	}

	if FromInt(1).IsZero() {
		t.Error("none zero int should not is zero")
	}
	if FromString("1").IsZero() {
		t.Error("none empty string should not is zero")
	}
}

func TestParse(t *testing.T) {
	id := Parse("123")
	if _, ok := id.val.(int64); !ok || id.IntValue() != int64(123) {
		t.Error("should not ")
	}
	id = Parse("abc")
	if _, ok := id.val.(string); !ok || id.String() != "abc" {
		t.Error("should not ")
	}
}

func TestID_Compare(t *testing.T) {
	if FromInt(123) != FromInt(123) {
		t.Error("compare failed")
	}
}
