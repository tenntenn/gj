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
	value  interface{}
	codec  *Codec
	parent *Value
	key    string
	index  int
}

func New(data []byte) (*Value, error) {
	var v interface{}
	err := DefaultCodec.Unmarshal(data, &v)
	return &Value{value: v, codec: &DefaultCodec}, err
}

func NewWithCodec(data []byte, codec Codec) (*Value, error) {
	var v interface{}
	err := codec.Unmarshal(data, &v)
	return &Value{value: v, codec: &codec}, err
}

func ValueOf(v interface{}) *Value {
	return &Value{value: v, codec: &DefaultCodec}
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

func (v *Value) ParentKey() string {
	return v.key
}

func (v *Value) ParentIndex() int {
	return v.index
}

func (v *Value) Parent() *Value {
	return v.parent
}

func (v *Value) Isolate() *Value {
	return &Value{value: v.value, codec: v.codec}
}

func (v *Value) IsObject() bool {

	if _, ok := v.value.(map[string]interface{}); ok {
		return true
	}

	if _, ok := v.value.(map[interface{}]interface{}); ok {
		return true
	}

	return false
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

	if f, ok := v.value.(float64); ok {
		return int64(f)
	}

	if n, ok := v.value.(int64); ok {
		return n
	}

	panic("v cannot convert to a int value.")
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
		return &Value{value: a[i], codec: v.codec, index: i, parent: v}
	}

	if s, ok := v.value.([]byte); ok {
		return &Value{value: string(s[i]), codec: v.codec, index: i, parent: v}
	}

	panic("v is not an array (slice) or string.")
}

func (v *Value) Slice(i, j int) *Value {

	if a, ok := v.value.([]interface{}); ok {
		return &Value{value: a[i:j], codec: v.codec, parent: v}
	}

	if s, ok := v.value.([]byte); ok {
		return &Value{value: string(s[i:j]), codec: v.codec, parent: v}
	}

	panic("v is not an array (slice) or string.")
}

func (v *Value) Get(k string) *Value {

	if o, ok := v.value.(map[string]interface{}); ok {
		return &Value{value: o[k], codec: v.codec, key: k, parent: v}
	}

	if o, ok := v.value.(map[interface{}]interface{}); ok {
		return &Value{value: o[k], codec: v.codec, key: k, parent: v}
	}

	panic("v is not an object (map).")
}

func (v *Value) Has(k string) bool {

	if o, ok := v.value.(map[string]interface{}); ok {
		_, has := o[k]
		return has
	}

	if o, ok := v.value.(map[interface{}]interface{}); ok {
		_, has := o[k]
		return has
	}

	panic("v is not an object (map).")
}

func (v *Value) TryGet(k string) (*Value, bool) {

	if o, ok := v.value.(map[interface{}]interface{}); ok {
		if fv, has := o[k]; has {
			return &Value{value: fv, codec: v.codec, key: k, parent: v}, true
		} else {
			return nil, false
		}
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

	if o, ok := v.value.(map[interface{}]interface{}); ok {
		keys := make([]string, 0, len(o))
		for k, _ := range o {
			keys = append(keys, k.(string))
		}
		return keys
	}

	panic("v is not an object (map).")
}

func (v *Value) EachIndex(f func(i int, v *Value) bool) {
	for i := 0; i < v.Len(); i++ {
		if end := f(i, v.Index(i)); end {
			break
		}
	}
}

func (v *Value) EachKey(f func(k string, v *Value) bool) {
	for _, k := range v.Keys() {
		if end := f(k, v.Get(k)); end {
			break
		}
	}
}

func (v *Value) Find(f func(v *Value) (ok, end bool)) <-chan *Value {
	ch := make(chan *Value)
	go func() {
		v.traversal(f, ch)
		close(ch)
	}()
	return ch
}

func (v *Value) traversal(f func(v *Value) (ok, end bool), ch chan<- *Value) {
	switch {
	case v.IsArray():
		v.EachIndex(func(i int, e *Value) bool {
			ok, end := f(e)
			if ok {
				ch <- e
			}

			if !end {
				e.traversal(f, ch)
			}

			return end
		})
	case v.IsObject():
		v.EachKey(func(k string, e *Value) bool {
			ok, end := f(e)
			if ok {
				ch <- e
			}

			if !end {
				e.traversal(f, ch)
			}

			return end
		})
	}
}
