package main

import (
	"fmt"
	"strings"

	. "github.com/nbjahan/go-launchbar"
)

func init() {
	var i *Item
	v := pb.NewView("main")

	i = v.NewItem("Pinboard: Recent")
	i.SetIcon("Pinboard")
	i.SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		c.Self.SetSubtitle(fmt.Sprintf("Last %d Bookmarks", c.Config.GetInt("recent-count")))
	})
	i.SetRun(func(c *Context) Items {
		items := &Items{}

		posts, err := getAllRecent(c.Action.IsControlKey())
		if err != nil {
			c.Logger.Printf("Error: getAllRecent %v", err)
			handleError(err)
			return nil
		}
		for _, post := range posts {
			items.Add(newPinItem(post))
		}
		return *items
	})

	i = v.NewItem("Pinboard: Unread")
	i.SetIcon("Pinboard")
	i.SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		c.Self.SetSubtitle(fmt.Sprintf("Last %d Unread Bookmarks", c.Config.GetInt("unread-count")))
	})
	i.SetRun(func(c *Context) Items {
		items := &Items{}

		posts, err := getAllPosts(c.Action.IsControlKey())
		if err != nil {
			c.Logger.Println("err", err)
			handleError(err)
			return nil
		}
		count := 0
		for _, post := range posts {
			if count == pb.Config.GetInt("unread-count") {
				break
			}
			if post.Toread {
				items.Add(newPinItem(post))
				count++
			}
		}
		return *items
	})

	i = v.NewItem("Pinboard: Search Posts")
	i.SetIcon("Pinboard")
	i.SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		if c.Input.IsString() && !c.Input.IsEmpty() {
			c.Self.SetSubtitle(fmt.Sprintf("query: %s", c.Input.String()))
			c.Self.SetOrder(-1)
		}
	})
	i.SetMatch(func(c *Context) bool {
		return c.Input.IsString()
	})
	i.SetRun(func(c *Context) Items {
		q := strings.TrimSpace(c.Input.String())
		if q == "" {
			return nil
		}
		posts, err := getAllPosts(c.Action.IsControlKey())
		items := &Items{}
		if err != nil {
			c.Logger.Println("err", err)
			handleError(err)
			return nil
		}
		matchedPosts := searchPosts(posts, q)
		for _, post := range matchedPosts {
			items.Add(newPinItem(post))
		}

		return *items
	})

	i = v.NewItem("Pinboard: Tags")
	i.SetIcon("Pinboard")
	i.SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		c.Self.SetSubtitle(fmt.Sprintf("Browse all %d tags", c.Config.GetInt("totaltags")))
	})
	i.SetRun(func(c *Context) Items {
		items := &Items{}
		tags, err := getAllTags(c.Action.IsControlKey())
		if err != nil {
			c.Logger.Println("err", err)
			handleError(err)
			return nil
		}
		for _, tag := range tags {
			i := NewItem(tag.Tag)
			i.SetAction(c.Config.GetString("actionDefaultScript"))
			i.SetSubtitle(tag.Count)
			i.SetIcon("TagTemplate")
			i.Run("showTags", tag.Tag)
			i.SetActionReturnsItems(true)
			i.SetActionRunsInBackground(false)
			items.Add(i)
		}
		return *items

	})

	i = v.NewItem("Pinboard: All Posts")
	i.SetIcon("Pinboard")
	i.SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		c.Self.SetSubtitle(fmt.Sprintf("Browse all %d Pins", c.Config.GetInt("totalposts")))
	})
	i.SetRun(func(c *Context) Items {
		posts, err := getAllPosts(c.Action.IsControlKey())
		if err != nil {
			c.Logger.Println("err", err)
			handleError(err)
			return nil
		}
		items := &Items{}
		for _, post := range posts {
			items.Add(newPinItem(post))
		}
		return *items
	})

	i = v.NewItem("My Pinboard")
	i.SetAction("")
	i.SetURL(fmt.Sprintf("https://pinboard.in/u:%s", strings.SplitN(pb.Config.GetString("token"), ":", 2)[0]))
	i.SetMatch(func(c *Context) bool { return c.Config.GetBool("loggedin") })

	i = v.NewItem("Pinboard: Preferences")
	i.SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/ToolbarAdvanced.icns")
	i.SetOrder(9998)

	i.SetRun(ShowViewFunc("config"))
	i.SetRender(func(c *Context) {
		if !c.Config.GetBool("loggedin") {
			c.Self.SetSubtitle("Enter here to logc.Input.")
		}
	})

}
