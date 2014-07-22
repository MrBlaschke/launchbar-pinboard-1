package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	json "github.com/bitly/go-simplejson"
	. "github.com/nbjahan/go-launchbar"
)

func search(query string) chan [2]string {
	// FIXME: sort search results
	uri := "https://pinboard.in/search/?all=Search+All&query=%s&start=%d"
	jobs := 4
	done := make(chan bool, jobs)
	posts := make(chan [2]string, jobs)
	go func() {
		for i := 0; i < jobs; i++ {
			go func(i int) {
				pb.Logger.Println(fmt.Sprintf(uri, query, i*20))
				d, _ := goquery.NewDocument(fmt.Sprintf(uri, query, i*20))
				d.Find("a.bookmark_title").Each(func(i int, s *goquery.Selection) {
					t := strings.TrimSpace(s.Text())
					h, _ := s.Attr("href")
					posts <- [2]string{t, h}
				})
				pb.Logger.Println(fmt.Sprintf(uri, query, i*20), "done")
				done <- true
			}(i)
		}
		for i := 0; i < jobs; i++ {
			<-done
		}
		close(posts)
	}()
	return posts
}

func init() {
	var i *Item
	v := pb.NewView("main")

	i = v.NewItem("Pinboard: Search All").SetIcon("Pinboard").SetOrder(4).SetSubtitle("Enter a query")
	i.SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		if c.Input.IsEmpty() {
			return
		}
		c.Self.SetOrder(-1)
		c.Self.SetSubtitle("Query: " + c.Input.String())
	})
	i.SetRun(func(c *Context) Items {
		if c.Input.IsEmpty() {
			return nil
		}

		m := make(map[string]string)
		for item := range search(url.QueryEscape(c.Input.String())) {
			m[item[0]] = item[1]
		}

		var items = &Items{}
		for t,url := range m {
			// TODO: parse posts as a pin
			items.Add(NewItem(t).SetURL(url))
		}
		return *items

		// d, _ := goquery.NewDocument(uri)
		// d.Find("a.bookmark_title").Each(func(i int, s *goquery.Selection) {
		// 	t := strings.TrimSpace(s.Text())
		// 	h, _ := s.Attr("href")
		// 	items.Add(NewItem(t).SetURL(h))
		// })
		// return *items
	})

	i = v.NewItem("Pinboard: Preferences").SetOrder(9998)
	i.SetIcon("DebugTemplate").SetRun(ShowViewFunc("config"))

	v = pb.GetView("main")
	i = v.NewItem("Pinboard: Popular").SetIcon("Pinboard")

	i.SetSubtitle("Hold ⌃(CTRL) to discard the cache.")
	i.SetActionRunsInBackground(false)
	i.SetMatch(func(c *Context) bool {
		return c.Input.String() == ""
	})
	i.SetRun(func(c *Context) Items {
		cache := c.Cache
		var items *Items
		if c.Action.IsControlKey() {
			items = nil
		} else {
			items = cache.GetItems("popular")
		}
		if items == nil {
			items = &Items{}

			url := "http://feeds.pinboard.in/json/popular/"
			res, err := http.Get(url)
			if err != nil {
				c.Logger.Println(err)
			}
			defer res.Body.Close()
			j, err := json.NewFromReader(res.Body)
			if err != nil {
				c.Logger.Println(err)
			}
			for _, post := range j.MustArray() {
				p := post.(map[string]interface{})
				title := p["u"].(string)
				if d, ok := p["d"].(string); ok && d != "" {
					title = d
				} else if n, ok := p["n"].(string); ok && n != "" {
					title = n
				}
				title = strings.Replace(title, "\n", " ", -1)
				url := p["u"].(string)
				items.Add(NewItem(title).SetURL(url))
			}
			cache.SetItems("popular", items, pb.Config.GetTimeDuration("cachelife"))
		}
		return *items
	})

	i = v.NewItem("Pinboard: Recent").SetIcon("Pinboard")

	i.SetSubtitle("Hold ⌃(CTRL) to discard the cache.")
	i.SetActionRunsInBackground(false)
	i.SetMatch(func(c *Context) bool {
		return c.Input.String() == ""
	})
	i.SetRun(func(c *Context) Items {
		cache := c.Cache
		var items *Items
		if c.Action.IsControlKey() {
			items = nil
		} else {
			items = cache.GetItems("recent")
		}
		if items == nil {
			items = &Items{}

			url := "http://feeds.pinboard.in/json/recent/"
			res, err := http.Get(url)
			if err != nil {
				c.Logger.Println(err)
			}
			defer res.Body.Close()
			j, err := json.NewFromReader(res.Body)
			if err != nil {
				c.Logger.Println(err)
			}
			for _, post := range j.MustArray() {
				p := post.(map[string]interface{})
				title := p["u"].(string)
				if d, ok := p["d"].(string); ok && d != "" {
					title = d
				} else if n, ok := p["n"].(string); ok && n != "" {
					title = n
				}
				title = strings.Replace(title, "\n", " ", -1)
				url := p["u"].(string)
				items.Add(NewItem(title).SetURL(url))
			}
			cache.SetItems("recent", items, pb.Config.GetTimeDuration("cachelife"))
		}
		return *items
	})

	v = pb.GetView("main")
	i = v.NewItem("Pinboard: Browse User").SetIcon("Pinboard")

	i.SetSubtitle("Enter a username")
	i.SetActionRunsInBackground(false)
	i.SetRender(func(c *Context) {
		if c.Input.String() != "" {
			c.Self.SetSubtitle("Query: " + c.Input.String())
		}
	})
	i.SetRun(func(c *Context) Items {
		var out Items
		c.Logger.Println(c.Input.String())
		if c.Input.String() == "" {
			c.Action.ShowView("main")
			return out
		}
		url := "http://feeds.pinboard.in/json/u:" + c.Input.String()
		res, err := http.Get(url)
		if err != nil {
			c.Logger.Println(err)
			c.Action.ShowView("main")
			return out
		}
		defer res.Body.Close()
		j, err := json.NewFromReader(res.Body)
		if err != nil {
			c.Logger.Println(err)
			c.Action.ShowView("main")
			return out
		}

		for _, post := range j.MustArray() {
			p := post.(map[string]interface{})
			title := p["u"].(string)
			if d, ok := p["d"].(string); ok && d != "" {
				title = d
			} else if n, ok := p["n"].(string); ok && n != "" {
				title = n
			}
			title = strings.Replace(title, "\n", " ", -1)
			url := p["u"].(string)
			(&out).Add(NewItem(title).SetURL(url))
		}
		return out
	})
}
