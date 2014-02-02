package gj

import (
	"encoding/json"
	"fmt"
)

type Codec struct {
	Marshal   func(v interface{}) (data []byte, err error)
	Unmarshal func(data []byte, v interface{}) (err error)
}

var (
	JSON         = Codec{jsonMarshal, jsonUnmarshal}
	DefaultCodec = JSON
)

func jsonMarshal(v interface{}) (data []byte, err error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v interface{}) (err error) {
	return json.Unmarshal(data, v)
}

type Value struct {
	value interface{}
	codec *Codec
}

func New(data []byte) (*Value, error) {
	var v interface{}
	err := DefaultCodec.Unmarshal(data, &v)
	return &Value{v, &DefaultCodec}, err
}

func ValueOf(v interface{}) *Value {
	return &Value{v, &DefaultCodec}
}

func (v *Value) Marshal() (data []byte, err error) {
	return v.codec.Marshal(v.value)
}

func (v *Value) Unmarshal(dst interface{}) (err error) {
	var data []byte
	if data, err = v.Marshal(); err != nil {
		return err
	}
	return v.codec.Unmarshal(data, dst)
}

func (v *Value) IsObject() bool {
	_, isObject := v.value.(map[string]interface{})
	return isObject
}

func (v *Value) IsArray() bool {
	_, isArray := v.value.([]interface{})
	return isArray
}

func (v *Value) IsNumber() bool {
	_, isNumber := v.value.(float64)
	return isNumber
}

func (v *Value) Int() int64 {
	n, ok := v.value.(float64)
	if !ok {
		panic("v cannot convert to a int value.")
	}

	return int64(n)
}

func (v *Value) Float() float64 {
	f, ok := v.value.(float64)
	if !ok {
		panic("v cannot convert to a float value.")
	}

	return f
}

func (v *Value) IsString() bool {
	_, isString := v.value.(string)
	return isString
}

func (v *Value) String() string {
	if s, ok := v.value.(string); ok {
		return s
	}

	return fmt.Sprintf("%v", v.value)
}

func (v *Value) IsBool() bool {
	_, isBool := v.value.(bool)
	return isBool
}

func (v *Value) Bool() bool {
	b, ok := v.value.(bool)
	if !ok {
		panic("v cannot convert to a bool value.")
	}

	return b
}

func (v *Value) Index(i int) *Value {

	if a, ok := v.value.([]interface{}); ok {
		return &Value{a[i], v.codec}
	}

	if s, ok := v.value.([]byte); ok {
		return &Value{string(s[i]), v.codec}
	}

	panic("v is not an array (slice) or string.")
}

func (v *Value) Slice(i, j int) *Value {

	if a, ok := v.value.([]interface{}); ok {
		return &Value{a[i:j], v.codec}
	}

	if s, ok := v.value.([]byte); ok {
		return &Value{string(s[i:j]), v.codec}
	}

	panic("v is not an array (slice) or string.")
}

func (v *Value) Get(k string) *Value {

	if o, ok := v.value.(map[string]interface{}); ok {
		return &Value{o[k], v.codec}
	}

	panic("v is not an object (map).")
}

func (v *Value) Has(k string) bool {
	if o, ok := v.value.(map[string]interface{}); ok {
		_, has := o[k]
		return has
	}

	panic("v is not an object (map).")
}

func (v *Value) Len() int {
	if a, ok := v.value.([]interface{}); ok {
		return len(a)
	}

	if s, ok := v.value.([]byte); ok {
		return len(s)
	}

	panic("v is not an array (slice) or string.")

}

func (v *Value) Keys() []string {

	if o, ok := v.value.(map[string]interface{}); ok {
		keys := make([]string, 0, len(o))
		for k, _ := range o {
			keys = append(keys, k)
		}
		return keys
	}

	panic("v is not an object (map).")
}

func (v *Value) EachIndex(f func(i int, v *Value)) {
	for i := 0; i < v.Len(); i++ {
		f(i, v.Index(i))
	}
}

func (v *Value) EachKey(f func(k string, v *Value)) {
	for _, k := range v.Keys() {
		f(k, v.Get(k))
	}
}
