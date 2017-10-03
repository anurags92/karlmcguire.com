package main

import (
	"io/ioutil"
	"sort"
)

type Posts []*Post

func (p Posts) Len() int {
	return len(p)
}

func (p Posts) Less(i, j int) bool {
	return p[j].Time.Before(p[i].Time)
}

func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Posts) Sort() {
	sort.Sort(p)
}

func (p Posts) Tags() Tags {
	var (
		tagm = make(map[string]*Tag)
		tags = make(Tags, 0)
	)
	for _, post := range p {
		for _, tag := range post.Tags {
			if _, ok := tagm[tag]; !ok {
				tagm[tag] = &Tag{
					Name:  tag,
					Posts: make([]*TagPost, 0),
				}
			}
			tagm[tag].Posts = append(tagm[tag].Posts, &TagPost{
				Title: post.Title,
				Date:  post.Date,
				Path:  post.Path,
			})
		}
	}
	for _, tag := range tagm {
		tags = append(tags, tag)
	}
	return tags
}

func (p Posts) Render(dir string) error {
	p.Sort()

	for _, post := range p {
		err := post.Render(dir)
		if err != nil {
			return err
		}
	}

	return Render(dir+"posts/", "templates/posts.html",
		&struct {
			Title       string
			Selected    string
			Posts       Posts
			Description string
		}{"Posts", "posts", p, "All of Karl McGuire's posts."},
	)
}

func ReadPosts(d string) (Posts, error) {
	ds, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, err
	}

	posts := make(Posts, 0)

	for _, sd := range ds {
		post, err := ReadPost(d + sd.Name() + "/")
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
