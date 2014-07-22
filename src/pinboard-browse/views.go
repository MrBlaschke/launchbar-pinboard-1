package main

import (
	"fmt"
	"time"

	. "github.com/nbjahan/go-launchbar"
)

func init() {
	pb.NewView("*").
		NewItem(fmt.Sprintf("Executed in: %v", time.Since(start))).
		SetSubtitle(fmt.Sprintf("%0.3f seconds", float64(time.Since(start))/float64(time.Second))).
		SetIcon("LoadingTemplate").
		SetMatch(MatchIfTrueFunc(pb.Config.GetBool("debug"))).
		SetOrder(9998)
}
