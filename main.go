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
	cname, err := ioutil.ReadFile("docs/CNAME")
	if err != nil {
		panic(err)
	}
	favicon, err := ioutil.ReadFile("docs/favicon.ico")
	if err != nil {
		panic(err)
	}
	os.RemoveAll("docs/")
	os.Mkdir("docs/", 0700)
	err = ioutil.WriteFile("docs/style.css", style, 0700)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("docs/CNAME", cname, 0700)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("docs/favicon.ico", favicon, 0700)
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

	// render index.html
	if err = Render("docs/", "templates/index.html",
		&struct {
			Title       string
			Selected    string
			Tags        Tags
			Posts       Posts
			Description string
		}{"Karl McGuire", "index", tags, posts, "Home description"}); err != nil {
		panic(err)
	}
}
