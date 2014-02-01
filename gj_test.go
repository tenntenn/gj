package gj

import (
	"encoding/json"
	"strings"
	"testing"
)

var jsonStrs = []string{
	`{"hoge": 100, "data": 1, "children":[{"data":2, "children":[]}, {"data":3, "children":[]}]}`,
	`{}`,
	`100`,
	`[]`,
	`[100]`,
	`"hoge"`,
}

func TestNewJsonValue(t *testing.T) {
	for _, jsonStr := range jsonStrs {
		_, err := NewJsonValue([]byte(jsonStr))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestNewJsonValueFromDecoder(t *testing.T) {
	for _, jsonStr := range jsonStrs {
		d := json.NewDecoder(strings.NewReader(jsonStr))
		_, err := NewJsonValueFromDecoder(d)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSlice(t *testing.T) {
	v, _ := NewJsonValue([]byte(jsonStrs[0]))
	if _, ok := v.Object().Get("children").Array().Slice().TryArray(); !ok {
		t.Error()
	}
}
