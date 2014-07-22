package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	sjson "github.com/bitly/go-simplejson"
	. "github.com/nbjahan/go-launchbar"
	"github.com/nbjahan/go-pinboard"
)

const (
	refreshAll int = 1 << iota
	refreshTags
	refreshRecent
)

var start = time.Now()
var pb *Action
var funcs = map[string]Func{
	"showTags": func(c *Context) Items {
		items := &Items{}
		posts, err := getAllPosts(false)
		if err != nil {
			c.Logger.Println("err", err)
			handleError(err)
			return nil
		}
		for _, post := range posts {
			if strings.Contains(post.Tag, c.Input.String()) {
				items.Add(newPinItem(post))
			}
		}
		return *items
	},
	"showView": func(c *Context) {
		c.Action.ShowView(c.Input.String())
	},
}

func init() {
	pb = NewAction("Pinboard", ConfigValues{
		"actionDefaultScript": "pinboard",
		"cachelife":           60 * time.Minute,
		"debug":               false,
		"recent-count":        100,
		"unread-count":        100,
	})
}

func newPinItem(post *pinboard.Post) *Item {
	i := NewItem(post.Title)
	i.SetURL(post.URL)
	i.Item().Data["post"] = post
	return i
}

func getAllTags(recache bool) (pinboard.TagCloud, error) {
	var tags pinboard.TagCloud
	pin := pinboard.New(pb.Config.GetString("token"))

	expiry, err := pb.Cache.Get("my-tags", &tags)
	if recache {
		if err == ErrCacheIsExpired {
			if t, err := pin.GetUpdatedTime(); err == nil {
				if t.After(*expiry) {
					tags = nil
				}
			}
		} else if err != nil || pb.Config.GetInt("refresh")&refreshTags == refreshTags {
			tags = nil
		}
	} else {
		if err != nil && err != ErrCacheIsExpired {
			tags = nil
		}
	}
	if tags == nil {
		v := pb.Config.GetString("view")
		pb.Config.Set("view", "working")
		tags, err = pin.GetAllTags()
		if err != nil {
			return nil, err
		}
		pb.Cache.Set("my-tags", tags, pb.Config.GetTimeDuration("cachelife"))
		pb.Config.Set("refresh", pb.Config.GetInt("refresh")&^refreshTags)
		pb.Config.Set("view", v)
	}
	pb.Config.Set("totaltags", len(tags))
	return tags, nil
}

func getAllRecent(recache bool) ([]*pinboard.Post, error) {
	var posts []*pinboard.Post
	pin := pinboard.New(pb.Config.GetString("token"))

	expiry, err := pb.Cache.Get("my-recent", &posts)
	if recache {
		if err == ErrCacheIsExpired {
			if t, err := pin.GetUpdatedTime(); err == nil {
				if t.After(*expiry) {
					posts = nil
				}
			}
		} else if err != nil || pb.Config.GetInt("refresh")&refreshRecent == refreshRecent {
			posts = nil
		}
	} else {
		if err != nil && err != ErrCacheIsExpired {
			posts = nil
		}
	}
	if posts == nil {
		v := pb.Config.GetString("view")
		pb.Config.Set("view", "working")
		posts, err = pin.GetRecentPostsWith().Count(pb.Config.GetInt("recent-count")).Do()
		if err != nil {
			return nil, err
		}
		pb.Cache.Set("my-recent", posts, pb.Config.GetTimeDuration("cachelife"))
		pb.Config.Set("refresh", pb.Config.GetInt("refresh")&^refreshRecent)
		pb.Config.Set("view", v)
	}
	return posts, nil
}

func getAllPosts(recache bool) ([]*pinboard.Post, error) {
	var posts []*pinboard.Post
	pin := pinboard.New(pb.Config.GetString("token"))

	expiry, err := pb.Cache.Get("my-posts", &posts)
	if recache {
		if err == ErrCacheIsExpired {
			if t, err := pin.GetUpdatedTime(); err == nil {
				if t.After(*expiry) {
					posts = nil
				}
			}
		} else if err != nil || pb.Config.GetInt("refresh")&refreshAll == refreshAll {
			posts = nil
		}
	} else {
		if err != nil && err != ErrCacheIsExpired {
			posts = nil
		}
	}
	if posts == nil {
		v := pb.Config.GetString("view")
		pb.Config.Set("view", "working")
		posts, err = pin.GetAllPosts()
		if err != nil {
			return nil, err
		}
		pb.Cache.Set("my-posts", posts, pb.Config.GetTimeDuration("cachelife"))
		pb.Config.Set("refresh", pb.Config.GetInt("refresh")&^refreshAll)
		pb.Config.Set("view", v)
	}
	pb.Config.Set("totalposts", len(posts))
	return posts, nil
}

func handleSuccess(msg string) {
	pb.Config.Set("success", msg)
	pb.ShowView("success")
}

func handleError(err error) {
	if err == pinboard.ErrForbidden {
		pb.Config.Set("error", err.Error())
		pb.Config.Set("token", "")
		pb.Config.Set("loggedin", false)
		pb.ShowView("main")
	} else {
		pb.Config.Set("error", err.Error())
		pb.ShowView("error")
	}
}

func isURL(in *Input) (string, bool) {
	link := in.String()
	if in.IsObject() {
		link = in.Item.Item().URL
	}
	if link != "" {
		_, err := url.Parse(link)
		if err == nil {
			if err := pinboard.ValidateURL(link); err == nil {
				return link, true
			}
		}
	}
	return "", false
}

func main() {
	view := pb.Config.GetString("view")
	if !pb.Config.GetBool("loggedin") && view == "main" || view == "" {
		pb.Config.Set("view", "login")
	}

	pb.Init(funcs)

	if pb.Config.GetBool("indev") {
		pb.Logger.Printf("in:\n%s\n", pb.Input.Raw())
	}

	if view != "login" {
		if url, ok := isURL(pb.Input); ok {
			pb.Config.Set("in-url", url)
			pb.Config.Set("in-title", pb.Input.Title())
			if i := pb.Input.Data("post"); i != nil {
				post, err := pinboard.ParsePost(i)
				if err == nil {
					pb.Config.Set("in-post", post)
				}
			} else {
				pb.Config.Delete("in-post")
				lines := strings.Split(pb.Input.String(), "\n")
				if len(lines) > 1 {
					// TODO: implement
					pb.Config.Set("error", "Not Implemented!")
					pb.Config.Set("error-desc", "Multiple URLs detected!")
					pb.ShowView("error")
					return
				}
			}
			pb.ShowView("modify")
			return
		}
	}
	out := pb.Run()

	if pb.Config.GetBool("indev") {
		nice := out
		js, err := sjson.NewJson([]byte(out))
		if err == nil {
			b, err := js.EncodePretty()
			if err == nil {
				nice = string(b)
			}
		}
		pb.Logger.Println("out:", string(nice))
	}

	fmt.Println(out)
}
