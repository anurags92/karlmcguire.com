package main

import (
	"encoding/json"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

type Post struct {
	Title    string
	Time     time.Time
	Date     string
	Tags     []string
	Markdown string
	Html     string
	Path     string
}

func (p *Post) Render(dir string) error {
	sort.Strings(p.Tags)

	return Render(dir+p.Path+"/", "templates/post.html",
		&struct {
			Title    string
			Selected string
			Post     *Post
			Max      int
		}{p.Title, " ", p, len(p.Tags) - 1},
	)
}

func ReadPost(d string) (*Post, error) {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, err
	}

	var (
		t    time.Time
		post struct {
			Title    string   `json:"title"`
			Tags     []string `json:"tags"`
			Date     string   `json:"date"`
			Markdown []byte
		}
	)

	for _, file := range files {
		if strings.Contains(file.Name(), "json") && file.Name()[len(file.Name())-4:len(file.Name())] == "json" {
			data, err := ioutil.ReadFile(d + file.Name())
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, &post)
			if err != nil {
				return nil, err
			}
			t, err = time.Parse("010206", post.Date)
			if err != nil {
				return nil, err
			}
		} else if strings.Contains(file.Name(), "md") && file.Name()[len(file.Name())-2:len(file.Name())] == "md" {
			post.Markdown, err = ioutil.ReadFile(d + file.Name())
			if err != nil {
				return nil, err
			}
		}
	}

	// remove any duplicate tags
	tagm := make(map[string]struct{})
	for _, tag := range post.Tags {
		tagm[tag] = struct{}{}
	}
	tags := make([]string, 0)
	for tag, _ := range tagm {
		tags = append(tags, tag)
	}

	//Path:     strings.Replace(strings.ToLower(post.Title), " ", "-", -1),

	return &Post{
		Title:    post.Title,
		Time:     t,
		Date:     t.Format("January 2, 2006"),
		Tags:     tags,
		Path:     strings.Split(d, "/")[1],
		Markdown: string(post.Markdown),
		Html:     string(blackfriday.MarkdownCommon(post.Markdown)),
	}, nil
}
