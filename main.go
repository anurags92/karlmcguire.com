package main

import (
	"io/ioutil"
	"os"
)

func main() {
	// grab style.css and reset docs/
	style, err := ioutil.ReadFile("docs/style.css")
	if err != nil {
		panic(err)
	}
	os.RemoveAll("docs/")
	os.Mkdir("docs/", 0700)
	err = ioutil.WriteFile("docs/style.css", style, 0700)
	if err != nil {
		panic(err)
	}

	// read posts and parse tags
	posts, err := ReadPosts("posts/")
	if err != nil {
		panic(err)
	}
	posts.Render("docs/")
	tags := posts.Tags()
	tags.Render("docs/")

	// render index.html and about/
	if err = Render("docs/", "templates/index.html",
		&struct {
			Selected string
			Tags     Tags
			Posts    Posts
		}{"index", tags, posts}); err != nil {
		panic(err)
	}
	if err = Render("docs/about/", "templates/about.html",
		&struct {
			Selected string
		}{"about"}); err != nil {
		panic(err)
	}
}
