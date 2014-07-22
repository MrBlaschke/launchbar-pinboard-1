package main

import . "github.com/nbjahan/go-launchbar"

func init() {
	var i *Item
	v := pb.NewView("working")

	i = v.NewItem("I'm Working. Please wait...").SetIcon("LoadingTemplate")
	i.SetActionRunsInBackground(false)
	i.SetAction("")

	i = v.NewItem("Cancel").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(func(c *Context) {
		//TODO: may I kill the job?
		c.Action.ShowView("main")
	})
}
