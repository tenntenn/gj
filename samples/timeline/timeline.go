package main

import (
	gj "../../"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	jsonStr, _ := ioutil.ReadAll(os.Stdin)
	v, _ := gj.New(jsonStr)
	traversal(v.Slice(1, v.Len()))
}

func traversal(v *gj.Value) {

	if v.IsArray() {
		v.EachIndex(func(i int, v *gj.Value) {
			traversal(v)
		})
		return
	}

	if v.IsObject() {

		if v.Has("type") && v.Get("type").String() == "ResourceSendRequest" {
			fmt.Println(v.Get("data").Get("url").String())
		}

		if v.Has("children") {
			traversal(v.Get("children"))
		}
	}
}
