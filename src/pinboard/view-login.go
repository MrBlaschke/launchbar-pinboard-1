package main

import (
	"strings"

	. "github.com/nbjahan/go-launchbar"

	"github.com/nbjahan/go-pinboard"
)

func init() {
	var i *Item
	v := pb.NewView("login")

	i = v.NewItem("Pinboard: Enter Auth Token")
	i.SetSubtitle("user:token")
	i.SetIcon("TokenTemplate")
	i.SetActionRunsInBackground(false)
	i.SetMatch(MatchIfFalseFunc(pb.Config.GetBool("loggedin")))
	i.SetRender(func(c *Context) {
		if c.Input.IsString() && (strings.Contains(c.Input.String(), ":") || c.Input.IsEmpty()) {
			c.Self.SetOrder(0)
		} else {
			c.Self.SetOrder(2)
		}
	})
	i.SetRun(func(c *Context) {
		if c.Config.GetBool("loggedin") {
			return
		}
		c.Config.Set("view", "working")
		p := pinboard.New(c.Input.String())
		_, err := p.GetUserSecret()
		if err != nil {
			c.Config.Set("error", err.Error())
			c.Action.ShowView("error")
			return
		}
		c.Config.Set("loggedin", true)
		c.Config.Set("token", c.Input.String())
		updateCache(c)
		// c.Action.ShowView("main")
	})
	i = v.NewItem("Pinboard: Enter Username")
	i.SetIcon("UserTemplate")
	i.SetMatch(MatchIfFalseFunc(pb.Config.GetBool("loggedin") || pb.Config.GetString("username") != ""))
	i.SetRender(func(c *Context) {
		if c.Input.IsString() && !c.Input.IsEmpty() {
			c.Self.SetSubtitle(c.Input.String())
		}
	})
	i.SetRun(func(c *Context) {
		if c.Config.GetBool("loggedin") || c.Config.GetString("username") != "" {
			return
		}
		c.Config.Set("username", c.Input.String())
		c.Action.ShowView("enterpassword")
	})

	i = v.NewItem("Find your Token").SetURL("https://pinboard.in/settings/password").SetAction("")

	i = v.NewItem("Pinboard: Preferences").SetOrder(9998)

	i.SetIcon("DebugTemplate").SetRun(ShowViewFunc("config"))
	i.SetRender(func(c *Context) {
		if !c.Config.GetBool("loggedin") {
			c.Self.SetSubtitle("Enter here to logc.Input.")
		}
	})
}
