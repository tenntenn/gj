package main

import (
	gj "../../"
	"fmt"
	"io/ioutil"
	"os"
)

func isSendRequest(v *gj.Value) bool {
	return v.ParentKey() == "type" && v.String() == "ResourceSendRequest"
}

func isTimeStamp(v *gj.Value, message string) bool {
	return v.ParentKey() == "type" &&
		v.String() == "TimeStamp" &&
		v.Parent().Get("data").Get("message").String() == message
}

func main() {
	jsonStr, _ := ioutil.ReadAll(os.Stdin)
	v, _ := gj.New(jsonStr)

	startTimes := []*gj.Value{}
	endTimes := []*gj.Value{}

	v.Traversal(func(v *gj.Value) bool {
		switch {
		case isSendRequest(v):
			startTimes = append(startTimes, v.Parent())
		case isTimeStamp(v, "hoge"):
			endTimes = append(endTimes, v.Parent())
		}
		return false
	})

	fmt.Println(startTimes)
	fmt.Println(endTimes)
}
