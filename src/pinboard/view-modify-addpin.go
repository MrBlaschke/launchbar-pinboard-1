package main

import (
	. "github.com/nbjahan/go-launchbar"
	"github.com/nbjahan/go-pinboard"
)

const (
	toread int = 1 << iota
	shared
	replace
)

func save(c *Context) Items {
	pin := pinboard.New(c.Config.GetString("token"))
	url := c.Config.GetString("in-url")
	title := c.Config.GetString("in-title")
	desc := c.Config.GetString("in-desc")
	flags := c.Config.GetInt("in-flags")
	tags := c.Config.GetString("in-tags")

	err := pin.AddPostWith(url, title).
		Shared(flags&shared == shared).
		Toread(flags&toread == toread).
		Replace(flags&replace == replace).
		Description(desc).
		Tag(tags).
		Do()
	if err != nil {
		handleError(err)
		return nil
	}

	// FIXME: should we reset in-* ?
	c.Config.Delete("in-title")
	c.Config.Delete("in-url")
	c.Config.Delete("in-desc")
	c.Config.Delete("in-tags")
	c.Config.Delete("in-flags")

	c.Config.Set("refresh", refreshRecent|refreshTags|refreshAll)

	post := &pinboard.Post{}
	post.Description = desc
	post.Title = title
	post.URL = url
	post.Shared = flags&shared == shared
	post.Toread = flags&toread == toread
	post.Tag = tags

	// TODO: add new post to cache

	c.Config.Set("view", "main")

	// post, _ := pinboard.ParsePost(c.Config.Get("in-post"))
	// post, _ := pinboard.ParsePost(p)
	v := pb.NewView("")
	// if post != nil {
	v.Items.Add(newPinItem(post))
	// } else {
	// v.NewItem(title).SetURL(url)
	// }
	v.NewItem("Pin added successfully").SetIcon("at.obdev.LaunchBar:GreenCheckmark").SetAction("")
	v.NewItem("Back to main").SetIcon("BackTemplate").Run("showView", "main")

	return v.Items
}

func init() {
	var i *Item
	v := pb.NewView("modify-addpin")

	i = v.NewItem("Enter Title")
	i.SetIcon("TitleTemplate")
	i.SetActionRunsInBackground(false)
	i.SetSubtitle(pb.Config.GetString("in-title"))
	i.SetRender(func(c *Context) {
		if c.Config.GetString("in-title") != "" {
			// if c.Config.GetString("in-tags") == "" {
			c.Self.SetOrder(2)
			// }
		}
		// if !c.Input.IsEmpty() {
		// 	c.Self.SetSubtitle(c.Input.String())
		// }
	})
	i.SetRun(func(c *Context) Items {
		title := c.Input.String()
		if title == "" {
			title = c.Config.GetString("in-url")
		}
		c.Config.Set("in-title", title)
		if c.Action.IsControlKey() {
			return save(c)
		}
		c.Action.ShowView("modify-addpin")
		return nil
	})

	i = v.NewItem("Enter Pin Tags")
	i.SetIcon("TagsTemplate")
	i.SetActionRunsInBackground(false)
	i.SetSubtitle(pb.Config.GetString("in-tags"))
	i.SetRun(func(c *Context) Items {
		c.Config.Set("in-tags", c.Input.String())
		if c.Action.IsControlKey() {
			return save(c)
		}
		c.Action.ShowView("modify-addpin")
		return nil
	})

	i = v.NewItem("Enter Pin Description")
	i.SetIcon("DescriptionTemplate")
	i.SetActionRunsInBackground(false)
	i.SetSubtitle(pb.Config.GetString("in-desc"))
	i.SetRun(func(c *Context) Items {
		desc := c.Input.String()
		c.Config.Set("in-desc", desc)
		if c.Action.IsControlKey() {
			return save(c)
		}
		c.Action.ShowView("modify-addpin")
		return nil
	})

	i = v.NewItem("Toread")
	i.SetRender(func(c *Context) {
		on := (c.Config.GetInt("in-flags")&toread == toread)
		if on {
			c.Self.SetIcon("OnTemplate")
			c.Self.SetSubtitle("currently: on")
		} else {
			c.Self.SetIcon("OffTemplate")
			c.Self.SetSubtitle("currently: off")
		}
	})
	i.SetRun(func(c *Context) {
		c.Config.Set("in-flags", c.Config.GetInt("in-flags")^toread)
		c.Action.ShowView("modify-addpin")
	})

	i = v.NewItem("Shared")
	i.SetRender(func(c *Context) {
		on := (c.Config.GetInt("in-flags")&shared == shared)
		if on {
			c.Self.SetIcon("OnTemplate")
			c.Self.SetSubtitle("currently: on")
		} else {
			c.Self.SetIcon("OffTemplate")
			c.Self.SetSubtitle("currently: off")
		}
	})
	i.SetRun(func(c *Context) {
		c.Config.Set("in-flags", c.Config.GetInt("in-flags")^shared)
		c.Action.ShowView("modify-addpin")
	})

	i = v.NewItem("Replace")
	i.SetMatch(func(c *Context) bool {
		return c.Config.Get("in-post") == nil
	})

	i.SetRender(func(c *Context) {
		on := (c.Config.GetInt("in-flags")&replace == replace)
		if on {
			c.Self.SetIcon("OnTemplate")
			c.Self.SetSubtitle("currently: on")
		} else {
			c.Self.SetIcon("OffTemplate")
			c.Self.SetSubtitle("currently: off")
		}
	})
	i.SetRun(func(c *Context) {
		c.Config.Set("in-flags", c.Config.GetInt("in-flags")^replace)
		c.Action.ShowView("modify-addpin")
	})

	i = v.NewItem("Save to Pinboard").SetActionRunsInBackground(false)
	i.SetSubtitle(pb.Config.GetString("in-url"))
	i.SetIcon("SaveTemplate")
	i.SetRun(func(c *Context) Items {
		return save(c)
	})
	i = v.NewItem("Back").SetSubtitle("hold CTRL to cancel").SetOrder(99999).SetIcon("BackTemplate")
	i.SetRun(func(c *Context) {
		if c.Action.IsControlKey() {
			c.Action.ShowView("main")
		} else {
			c.Action.ShowView("modify")
		}

	})

}
