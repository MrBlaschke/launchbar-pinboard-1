package main

import (
	"fmt"
	"time"

	. "github.com/nbjahan/go-launchbar"
)

var start = time.Now()
var pb *Action

func init() {
	pb = NewAction("Pinboard Browse", ConfigValues{
		"actionDefaultScript": "pinboard",
		"cachelife":           60 * time.Minute,
		"debug":               false,
	})
}

func main() {
	pb.Init()
	// pb.Logger.Println("in:" + pb.Input.Raw())
	out := pb.Run()
	// nice := out
	// js, err := sjson.NewJson([]byte(out))
	// if err == nil {
	// 	b, err := js.EncodePretty()
	// 	if err == nil {
	// 		nice = string(b)
	// 	}
	// }

	// pb.Logger.Println("out:", string(nice))
	fmt.Println(out)

}
