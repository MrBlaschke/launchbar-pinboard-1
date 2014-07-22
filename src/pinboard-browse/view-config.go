package main

import (
	"fmt"
	"time"

	. "github.com/nbjahan/go-launchbar"
)

func init() {
	v := pb.NewView("config")

	i := v.NewItem("Debug")
	i.SetRender(func(c *Context) {
		if c.Config.GetBool("debug") {
			c.Self.SetIcon("OnTemplate")
			c.Self.SetSubtitle("currently: on")
		} else {
			c.Self.SetIcon("OffTemplate")
			c.Self.SetSubtitle("currently: off")
		}
	})
	i.SetRun(func(c *Context) Items {
		c.Self.SetOrder(-1)
		c.Config.Set("debug", !c.Config.GetBool("debug"))
		c.Action.ShowView("config")
		return nil
	})

	i = v.NewItem("Set Cache Time (Minutes)")
	i.SetIcon("LoadingTemplate")
	i.SetSubtitle(fmt.Sprintf("currently: %v", pb.Config.GetTimeDuration("cachelife")))
	i.SetMatch(func(c *Context) bool {
		return c.Input.IsNumber() || c.Input.IsEmpty()
	})

	i.SetRender(func(c *Context) {
		if c.Input.IsNumber() {
			c.Self.SetOrder(-1)
			c.Self.SetSubtitle(fmt.Sprintf("%v", time.Duration(c.Input.Float64()*float64(time.Minute))))
		}
	})
	i.SetRun(func(c *Context) {
		if c.Input.IsNumber() {
			c.Config.Set("cachelife", time.Duration(c.Input.Float64()*float64(time.Minute)))
		}
		c.Action.ShowView("config")
	})

	i = v.NewItem("Back").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(ShowViewFunc("main"))
}
