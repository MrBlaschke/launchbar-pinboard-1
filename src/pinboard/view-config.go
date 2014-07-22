package main

import (
	"fmt"
	"path"
	"time"

	. "github.com/nbjahan/go-launchbar"
)

func init() {
	var i *Item
	v := pb.NewView("config")

	i = v.NewItem("Debug")
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

	i = v.NewItem("Pinboard: Logout")
	i.SetIcon("LogoutTemplate")
	i.SetSubtitle("Logged in as " + pb.Config.GetString("username"))
	i.SetMatch(MatchIfTrueFunc(pb.Config.GetBool("loggedin")))
	i.SetRun(func(c *Context) {
		//TODO: cleanup and notify
		c.Config.Set("loggedin", false)
		c.Config.Set("token", "")
		c.Config.Set("username", "")

		c.Config.Delete("totaltags")
		c.Config.Delete("totalposts")
		c.Config.Delete("error")

		c.Cache.Clean()

		c.Action.ShowView("main")
	})

	i = v.NewItem("Paths")
	i.SetActionRunsInBackground(false)
	i.SetActionReturnsItems(true)
	i.SetIcon("at.obdev.LaunchBar:Category")
	i.SetMatch(MatchIfTrueFunc(pb.Config.GetBool("debug")))
	i.SetRun(func(c *Context) Items {
		items := &Items{}
		items.Add(NewItem("error.log").SetPath(path.Join(pb.SupportPath(), "error.log")))
		items.Add(NewItem("config.json").SetPath(path.Join(pb.SupportPath(), "config.json")))
		return *items
	})

	i = v.NewItem("Back").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(ShowViewFunc("main"))
}
