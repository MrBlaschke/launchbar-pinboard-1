package main

import . "github.com/nbjahan/go-launchbar"

func init() {
	var i *Item
	v := pb.NewView("success")

	i = v.NewItem("Success")
	i.SetIcon("SuccessTemplate")
	i.SetIcon("at.obdev.LaunchBar:GreenCheckmark")
	i.SetRender(func(c *Context) {
		e := c.Config.GetString("success")
		if e != "" {
			c.Self.SetTitle(e)
		}
	})
	i.SetAction("")

	i = v.NewItem("Back").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(func(c *Context) {
		c.Config.Set("success", "")
		c.Action.ShowView("main")
	})
}
