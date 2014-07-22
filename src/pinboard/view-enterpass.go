package main

import (
	. "github.com/nbjahan/go-launchbar"

	"github.com/nbjahan/go-pinboard"
)

func init() {
	var i *Item
	v := pb.NewView("enterpassword")

	i = v.NewItem("Pinboard: Enter Password")
	i.SetSubtitle("your password will be used once to get your token")
	i.SetIcon("KeyTemplate")
	i.SetActionRunsInBackground(false)
	i.SetMatch(MatchIfFalseFunc(pb.Config.GetBool("loggedin") || pb.Config.GetString("username") == ""))
	i.SetRun(func(c *Context) {
		if c.Config.GetBool("loggedin") || c.Config.GetString("username") == "" {
			return
		}
		c.Config.Set("view", "working")
		token, err := pinboard.GetAuthToken(c.Config.GetString("username"), c.Input.String())
		if err != nil {
			c.Logger.Printf("token: %q", token)
			c.Logger.Printf("err: %#v\n", err)
			c.Config.Set("error", err.Error())
			c.Action.ShowView("loginfailed")
			return
		}
		c.Config.Set("view", "main")
		c.Config.Set("loggedin", true)
		c.Config.Set("token", token)
		updateCache()
		// c.Action.ShowView("main")
	})

	i = v.NewItem("Back").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(func(c *Context) {
		c.Config.Set("username", "")
		c.Action.ShowView("main")
	})

	v = pb.NewView("loginfailed")
	i = v.NewItem("Error").SetIcon("AlertTemplate")
	i.SetRender(func(c *Context) {
		if e := c.Config.GetString("error"); e != "" {
			c.Self.SetTitle(e)
		}
	})
	i.SetAction("")

	i = v.NewItem("Back").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(func(c *Context) {
		c.Config.Set("error", "")
		c.Action.ShowView("enterpassword")
	})

}
