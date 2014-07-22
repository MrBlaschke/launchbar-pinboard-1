package main

import . "github.com/nbjahan/go-launchbar"

func init() {
	var i *Item
	v := pb.NewView("error")

	i = v.NewItem("Error")
	i.SetIcon("AlertTemplate")
	i.SetRender(func(c *Context) {
		if e := c.Config.GetString("error"); e != "" {
			c.Self.SetTitle(e)
			if sub := c.Config.GetString("error-desc"); sub != "" {
				c.Self.SetSubtitle(sub)
			}
		}
	})
	i.SetAction("")

	i = v.NewItem("Back to main").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(func(c *Context) {
		c.Config.Delete("error")
		c.Config.Delete("error-desc")
		c.Action.ShowView("main")
	})
}
