package gq

import (
	"encoding/json"
)

type JsonValue struct {
	value interface{}
}

func NewJsonValue(jsonStr []byte) (*JsonValue, error) {
	var v interface{}
	err := json.Unmarshal(jsonStr, &v)
	return Value(v), err
}

func NewJsonValueFromDecoder(decoder *json.Decoder) (*JsonValue, error) {
	var v interface{}
	err := decoder.Decode(&v)
	return Value(v), err
}

func Value(v interface{}) *JsonValue {
	return &JsonValue{v}
}

func (v *JsonValue) Json() []byte {
	// TODO: エラー処理
	jsonStr, _ := json.Marshal(v.value)
	return jsonStr
}

func (v *JsonValue) TryObject() (Object, bool) {
	o, ok := v.value.(map[string]interface{})
	return Object(o), ok
}

func (v *JsonValue) Object() Object {
	o, _ := v.TryObject()
	return o
}

func (v *JsonValue) TryArray() (Array, bool) {
	a, ok := v.value.([]interface{})
	return Array(a), ok
}

func (v *JsonValue) Array() Array {
	a, _ := v.TryArray()
	return a
}

func (v *JsonValue) TryNumber() (Number, bool) {
	n, ok := v.value.(float64)
	return Number(n), ok
}

func (v *JsonValue) Number() Number {
	n, _ := v.TryNumber()
	return n
}

func (v *JsonValue) TryString() (string, bool) {
	s, ok := v.value.(string)
	return s, ok
}

func (v *JsonValue) String() string {
	s, _ := v.TryString()
	return s
}

func (v *JsonValue) TryBool() (bool, bool) {
	b, ok := v.value.(bool)
	return b, ok
}

func (v *JsonValue) Bool() bool {
	b, _ := v.TryBool()
	return b
}

type Object map[string]interface{}

func (o Object) Get(key string) *JsonValue {
	return Value(o[key])
}

type Array []interface{}

func (a Array) Get(index int) *JsonValue {
	return Value(a[index])
}

type Number float64

func (n Number) Int() int {
	return int(n)
}

func (n Number) Float64() float64 {
	return float64(n)
}
