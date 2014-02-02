package main

import (
	gj "../../"
	"fmt"
	"github.com/vmihailenco/msgpack"
)

func marshal(v interface{}) (data []byte, err error) {
	return msgpack.Marshal(v)
}

func unmarshal(data []byte, v interface{}) (err error) {
	return msgpack.Unmarshal(data, v)
}

func main() {
	codec := gj.Codec{marshal, unmarshal}
	m := map[string]interface{}{
		"hoge": 100,
		"piyo": 2.5,
	}

	data, _ := msgpack.Marshal(m)
	v, _ := gj.NewWithCodec(data, codec)

	fmt.Println(v.Get("hoge").Int())
	fmt.Println(v.Get("piyo").Float())
}
