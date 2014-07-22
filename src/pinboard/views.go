package main

import (
	"fmt"
	"time"

	. "github.com/nbjahan/go-launchbar"
)

func updateCache() {

	done := make(chan struct{}, 3)
	v := pb.Config.GetString("view")
	pb.Config.Set("view", "working")
	go func() { getAllPosts(true); done <- struct{}{} }()
	go func() { getAllTags(true); done <- struct{}{} }()
	go func() { getAllRecent(true); done <- struct{}{} }()
	<-done
	<-done
	<-done
	// c.Action.ShowView(v)
	pb.Config.Set("view", v)
}
func init() {
	var i *Item
	v := pb.NewView("*")
	i = v.NewItem("Cahce is obsolete!")
	i.SetSubtitle("Enter to refresh (⌃Enter to ignore)")
	i.SetActionRunsInBackground(false)
	i.SetOrder(9997)
	i.SetIcon("Alert")
	i.SetMatch(MatchIfTrueFunc(pb.Config.GetInt("refresh") > 0))
	i.SetRun(func(c *Context) {
		if c.Action.IsControlKey() {
			c.Config.Set("refresh", 0)
			c.Action.ShowView(c.Config.GetString("view"))
			return
		}
		updateCache()
	})

	v.NewItem(fmt.Sprintf("Executed in: %v", time.Since(start))).
		SetSubtitle(fmt.Sprintf("%0.3f seconds", float64(time.Since(start))/float64(time.Second))).
		// SetIcon("LoadingTemplate").
		SetIcon("com.apple.Safari:ToolbarWebInspectorTemplate").
		SetMatch(MatchIfTrueFunc(pb.Config.GetBool("debug"))).
		SetOrder(9998)
}
