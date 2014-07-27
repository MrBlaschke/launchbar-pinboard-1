package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/DHowett/go-plist"
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

		c.Cache.Delete("my-recent")
		c.Cache.Delete("my-tags")
		c.Cache.Delete("my-posts")

		c.Action.ShowView("main")
	})

	i = v.NewItem("Check for Action Updates ...")
	i.SetRender(func(c *Context) {
		version := pb.Config.GetString("newversion")
		if version != "" {
			c.Self.SetSubtitle(fmt.Sprintf("latest version: v%s (I'm: v%s)", version, pb.Version()))
		} else {
			c.Self.SetSubtitle("click to check")
		}
	})
	i.SetActionRunsInBackground(false)
	i.SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/ToolbarDownloadsFolderIcon.icns")
	i.SetRun(func(c *Context) Items {
		v := pb.NewView("")
		version := ""
		if resp, err := http.Get("https://raw.githubusercontent.com/nbjahan/launchbar-pinboard/master/src/Info.plist"); err == nil {
			defer resp.Body.Close()
			if data, err := ioutil.ReadAll(resp.Body); err == nil {
				var v map[string]interface{}
				if _, err := plist.Unmarshal(data, &v); err == nil {
					version = v["CFBundleVersion"].(string)
				}
			}
		}
		if version == "" {
			v.NewItem("Cannot get the version").SetIcon("Alert").SetSubtitle("Please try again later.").SetAction("")
		} else {
			c.Config.Set("newversion", version)
			if c.Action.Version().Cmp(Version(version)) < 0 {
				dl := fmt.Sprintf("https://github.com/nbjahan/launchbar-pinboard/releases/download/v%s/Pinboard-%s.lbext", version, version)
				v.NewItem(fmt.Sprintf("Download New version v%s (I'm v%s)", version, c.Action.Version())).
					SetIcon("/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/ToolbarDownloadsFolderIcon.icns").
					SetURL(dl).SetSubtitle(dl).SetAction("")
			} else {
				v.NewItem("I'm up to date!").SetAction("").
					SetSubtitle(fmt.Sprintf("latest version is %s (I'm %s)", version, c.Action.Version())).
					SetIcon("at.obdev.LaunchBar:GreenCheckmark")
			}
		}
		v.NewItem("Back").SetIcon("BackTemplate").Run("showView", "config")
		return v.Items
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
