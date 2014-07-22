package main

import (
	"sort"

	"github.com/nbjahan/go-pinboard"
	"github.com/nbjahan/gofuzz"
)

var scores []float64

type byScore []*pinboard.Post

func (b byScore) Len() int           { return len(b) }
func (b byScore) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byScore) Less(i, j int) bool { return scores[i] < scores[j] }

func searchPosts(posts []*pinboard.Post, query string) []*pinboard.Post {
	var out []*pinboard.Post
	scores = make([]float64, 0)
	f := gofuzz.New()
	f.Threshold = 0.4
	if len(query) > f.MaxBits {
		query = query[:f.MaxBits]
	}

	for _, post := range posts {
		totalScore := float64(1)
		var score float64
		_, score = f.Search(post.Title, query, 0)
		if score < totalScore {
			totalScore = score
		}
		_, score = f.Search(post.URL, query, 0)
		if score < totalScore {
			totalScore = score
		}
		_, score = f.Search(post.Tag, query, 0)
		if score < totalScore {
			totalScore = score
		}
		_, score = f.Search(post.Description, query, 0)
		if score < totalScore {
			totalScore = score
		}
		if totalScore < 1 {
			out = append(out, post)
			scores = append(scores, totalScore)
		}
	}

	sort.Sort(byScore(out))
	return out
}
