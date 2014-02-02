package gj

import (
	"testing"
)

var jsonStrs = []string{
	`{"hoge": {"foo":100, "bar":200}, "data": 1, "children":[{"data":2, "children":[]}, {"data":3, "children":[]}]}`,
	`{"hoge":100, "piyo":200}`,
	`100`,
	`[]`,
	`[100, 200]`,
	`"hoge"`,
	`false`,
}

func TestNew(t *testing.T) {
	for _, jsonStr := range jsonStrs {
		_, err := New([]byte(jsonStr))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestMarshal(t *testing.T) {
	v, _ := New([]byte(jsonStrs[2]))
	data, err := v.Marshal()
	if err != nil {
		t.Error(err)
	}

	if string(data) != "100" {
		t.Error()
	}
}

func TestUnmarshal(t *testing.T) {

	var dst struct {
		Foo int `json:"foo"`
		Bar int `json:"bar"`
	}

	v, _ := New([]byte(jsonStrs[0]))
	if err := v.Get("hoge").Unmarshal(&dst); err != nil {
		t.Error(err)
	}

	if dst.Foo != 100 || dst.Bar != 200 {
		t.Error()
	}
}

func TestIsObject(t *testing.T) {
	isObject := []bool{true, true, false, false, false, false, false}
	for i, jsonStr := range jsonStrs {
		v, _ := New([]byte(jsonStr))
		if actual := v.IsObject(); actual != isObject[i] {
			t.Error("expect ", isObject[i], "but actual ", actual)
		}
	}
}

func TestIsArray(t *testing.T) {
	isArray := []bool{false, false, false, true, true, false, false}
	for i, jsonStr := range jsonStrs {
		v, _ := New([]byte(jsonStr))
		if actual := v.IsArray(); actual != isArray[i] {
			t.Error("expect ", isArray[i], "but actual ", actual)
		}
	}
}

func TestIsNumber(t *testing.T) {
	isNumber := []bool{false, false, true, false, false, false, false}
	for i, jsonStr := range jsonStrs {
		v, _ := New([]byte(jsonStr))
		if actual := v.IsNumber(); actual != isNumber[i] {
			t.Error("expect ", isNumber[i], "but actual ", actual)
		}
	}
}

func TestInt(t *testing.T) {
	v, _ := New([]byte(jsonStrs[2]))
	if v.Int() != 100 {
		t.Error()
	}
}

func TestFloat(t *testing.T) {
	v, _ := New([]byte(jsonStrs[2]))
	if v.Float() != float64(100) {
		t.Error()
	}
}

func TestIsString(t *testing.T) {
	isString := []bool{false, false, false, false, false, true, false}
	for i, jsonStr := range jsonStrs {
		v, _ := New([]byte(jsonStr))
		if actual := v.IsString(); actual != isString[i] {
			t.Error("expect ", isString[i], "but actual ", actual)
		}
	}
}

func TestString(t *testing.T) {
	v, _ := New([]byte(jsonStrs[5]))
	if v.String() != "hoge" {
		t.Error()
	}
}

func TestIsBool(t *testing.T) {
	isBool := []bool{false, false, false, false, false, false, true}
	for i, jsonStr := range jsonStrs {
		v, _ := New([]byte(jsonStr))
		if actual := v.IsBool(); actual != isBool[i] {
			t.Error("expect ", isBool[i], "but actual ", actual)
		}
	}
}

func TestBool(t *testing.T) {
	v, _ := New([]byte(jsonStrs[6]))
	if v.Bool() != false {
		t.Error()
	}
}

func TestIndex(t *testing.T) {
	v, _ := New([]byte(jsonStrs[4]))
	if v.Index(0).Int() != 100 {
		t.Error()
	}
}

func TestSlice(t *testing.T) {
	v, _ := New([]byte(jsonStrs[0]))
	if v.Get("children").Slice(1, 2).Index(0).Get("data").Int() != 3 {
		t.Error()
	}
}

func TestGet(t *testing.T) {
	v, _ := New([]byte(jsonStrs[0]))
	if v.Get("data").Int() != 1 {
		t.Error()
	}
}

func TestHas(t *testing.T) {
	v, _ := New([]byte(jsonStrs[0]))
	if !v.Has("data") {
		t.Error()
	}

	if v.Has("foo") {
		t.Error()
	}
}

func TestLen(t *testing.T) {
	v, _ := New([]byte(jsonStrs[4]))
	if v.Len() != 2 {
		t.Error()
	}
}

func TestKeys(t *testing.T) {
	keys := map[string]bool{"hoge": false, "children": false, "data": false}
	count := len(keys)
	v, _ := New([]byte(jsonStrs[0]))
	for _, k := range v.Keys() {
		if _, ok := keys[k]; ok {
			count--
			keys[k] = true
		}
	}

	if count != 0 {
		t.Error()
	}
}

func TestEachIndex(t *testing.T) {
	v, _ := New([]byte(jsonStrs[4]))
	values := []int64{100, 200}
	v.EachIndex(func(i int, v *Value) {
		if actual := v.Int(); values[i] != actual {
			t.Error("expect ", values[i], "but actual ", actual)
		}
	})
}

func TestEachKey(t *testing.T) {
	v, _ := New([]byte(jsonStrs[1]))
	values := map[string]int64{
		"hoge": 100,
		"piyo": 200,
	}
	v.EachKey(func(k string, v *Value) {
		if actual := v.Int(); values[k] != actual {
			t.Error("expect ", values[k], "but actual ", actual)
		}
	})
}
