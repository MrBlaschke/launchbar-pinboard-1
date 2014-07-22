package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	. "github.com/nbjahan/go-launchbar"
	"github.com/nbjahan/go-pinboard"
)

func getPageInfo(u string) map[string]string {
	//TODO: getPageInfo: don't download the whole page it could be large binary file!
	d, _ := goquery.NewDocument(u)
	info := make(map[string]string)
	info["title"] = strings.TrimSpace(d.Find("title").Text())
	if desc, found := d.Find("meta[name=description]").Attr("content"); found {
		info["description"] = strings.TrimSpace(desc)
	}
	return info
}
func getTag(u string) string {
	pin := pinboard.New(pb.Config.GetString("token"))
	if popular, _, err := pin.GetSuggestedTags(u); err == nil {
		if len(popular) > 0 {
			return strings.Join(popular, " ")
		}
	}
	return ""
}

func init() {
	var i *Item
	v := pb.NewView("modify")

	i = v.NewItem("Pinboard: Quick Add Pin")
	i.SetSubtitle("⌃=Autofill, ⇧=Overwrite (You can enter a title)")
	i.SetIcon("AddTemplate").SetActionRunsInBackground(false)
	i.SetMatch(func(c *Context) bool {
		return c.Config.Get("in-post") == nil
	})
	i.SetRun(func(c *Context) Items {
		url := c.Config.GetString("in-url")
		title := c.Input.String()
		tag := ""
		desc := ""
		if c.Action.IsControlKey() {
			info := getPageInfo(url)
			title = info["title"]
			desc = info["description"]
			tag = getTag(url)
		}

		overwrite := false
		if c.Action.IsShiftKey() {
			overwrite = true
		}

		if title == "" {
			title = c.Config.GetString("in-title")
			if title == "" {
				title = url
			}
		}

		pin := pinboard.New(c.Config.GetString("token"))
		err := pin.AddPostWith(url, title).
			Replace(overwrite).
			// TODO: shared should be an option
			Shared(c.Config.GetBool("public-post")).
			Toread(true).
			Tag(tag).
			Description(desc).
			Do()
		if err != nil {
			// TODO: replace if exists
			// if err.Error() == "item already exists"
			// give a change to replace
			c.Logger.Printf("Error on Quick Add: %v", err)
			handleError(err)
			return nil
		}

		c.Config.Set("refresh", refreshRecent|refreshTags|refreshAll)

		v := pb.NewView("")
		// v.NewItem(title).SetURL(url)
		post := &pinboard.Post{}
		post.Description = desc
		post.Title = title
		post.URL = url
		post.Shared = c.Config.GetBool("public-post")
		post.Toread = true
		post.Tag = tag

		// FIXME: should we reset in-* ?
		c.Config.Delete("in-title")
		c.Config.Delete("in-url")
		c.Config.Delete("in-desc")
		c.Config.Delete("in-tags")
		c.Config.Delete("in-flags")

		// TODO: add new post to cache

		c.Config.Set("view", "main")

		v.Items.Add(newPinItem(post))
		v.NewItem("Pin added successfully").SetIcon("at.obdev.LaunchBar:GreenCheckmark").SetAction("")
		v.NewItem("Back").SetIcon("BackTemplate").Run("showView", "main")
		return v.Items
		// else {
		// 	handleSuccess("Pin added successfully", url)
		// }
	})

	i = v.NewItem("Pinboard: Add/Edit Pin")
	i.SetActionReturnsItems(true)
	i.SetSubtitle(pb.Config.GetString("in-url"))
	i.SetIcon("AddTemplate").SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		if c.Config.Get("in-post") != nil {
			c.Self.SetTitle("Pinboard: Edit Pin")
			c.Self.SetIcon("EditTemplate")
		}
	})
	i.SetRun(func(c *Context) Items {
		url := c.Config.GetString("in-url")
		// title := c.Config.GetString("in-title")
		if data := c.Config.Get("in-post"); data != nil {
			if post, err := pinboard.ParsePost(data); err == nil {
				c.Config.Set("in-title", post.Title)
				c.Config.Set("in-desc", post.Description)
				flags := replace
				if post.Toread {
					flags |= toread
				}
				if post.Shared {
					flags |= shared
				}
				c.Config.Set("in-flags", flags)
				c.Config.Set("in-tags", post.Tag)
			} else {
				c.Config.Delete("in-desc")
				c.Config.Delete("in-tags")
			}
		}
		if c.Action.IsControlKey() {
			c.Config.Set("in-title", getPageInfo(url)["title"])
			c.Config.Set("in-desc", getPageInfo(url)["description"])
			c.Config.Set("in-tags", getTag(url))
		}
		c.Action.ShowView("modify-addpin")
		return nil
	})

	// i = v.NewItem("Pinboard: Add with Info")
	// i.SetIcon("AddTemplate")
	//

	i = v.NewItem("Delete Pin").SetIcon("TrashTemplate")
	i.SetSubtitle(pb.Config.GetString("in-url"))
	i.SetActionRunsInBackground(false)
	i.SetRun(func(c *Context) Items {
		url := c.Config.GetString("in-url")
		pin := pinboard.New(c.Config.GetString("token"))
		err := pin.DeletePost(url)
		if err != nil {
			handleError(err)
			return nil
		}
		post, _ := pinboard.ParsePost(c.Config.Get("in-post"))
		v := pb.NewView("")
		if post != nil {
			v.Items.Add(newPinItem(post))
		} else {
			title := c.Config.GetString("in-title")
			if title == "" {
				title = url
			}
			v.NewItem(title).SetURL(url)
		}
		// v.NewItem(title).SetURL(url)
		v.NewItem("Pin removed successfully").SetIcon("at.obdev.LaunchBar:GreenCheckmark").SetAction("")
		// v.NewItem("Back").SetIcon("BackTemplate").SetFuncName("showView").SetArg("main")
		v.NewItem("Back to main").SetIcon("BackTemplate").Run("showView", "main")
		return v.Items
		// } else {
		// handleSuccess("Pin removed successfully")
		// }
	})

	// post, _ := pinboard.ParsePost(pb.Config.Get("in-post"))
	// if post != nil {
	// v.AddItem(newPinItem(post))
	// } else {
	i = v.NewItem(pb.Config.GetString("in-title"))
	i.SetAction("")
	i.SetSubtitle(pb.Config.GetString("in-url"))
	i.SetURL(pb.Config.GetString("in-url"))
	i.SetRender(func(c *Context) {
		if c.Self.Item().Title == "" {
			c.Self.SetTitle(c.Self.Item().URL)
		}
	})
	// }
	i = v.NewItem("Back to main").SetOrder(99999)
	i.SetIcon("BackTemplate")
	i.SetRun(ShowViewFunc("main"))

}
