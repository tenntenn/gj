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
	for v := range v.Find(func(v *gj.Value) (ok, end bool) {
		return v.ParentKey() == "type" && v.String() == "ResourceSendRequest", false
	}) {
		fmt.Println(v.Parent().Get("data").Get("url"))
	}
}
