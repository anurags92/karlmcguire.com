package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/xba/html"
	"gopkg.in/russross/blackfriday.v2"
)

const (
	YEAR = "2019"
	NOTE = "~"
)

func PAGE(head, header, content, footer []byte) []byte {
	return bytes.ReplaceAll(html.C(
		html.T("!doctype html", nil),
		html.E("html", html.A("lang", "en"),
			html.C(
				head,
				html.E("body", nil,
					html.C(
						header,
						content,
						footer))))), []byte{'\x00'}, []byte(""))
}

func HEAD(title, desc string) []byte {
	return html.E("head", nil,
		html.C(
			html.T("meta",
				html.A(
					"charset", "utf-8")),
			html.T("meta",
				html.A(
					"http-equiv", "X-UA-Compatible",
					"content", "IE=edge")),
			html.T("meta",
				html.A(
					"name", "viewport",
					"content", "width=device-width, initial-scale=1")),
			html.T("meta",
				html.A(
					"name", "description",
					"content", desc)),
			html.E("title", nil, []byte(title)),
			html.T("link",
				html.A(
					"href",
					"https://fonts.googleapis.com/css?family=Roboto+Mono",
					"rel", "stylesheet")),
			html.T("link",
				html.A(
					"rel", "stylesheet",
					"href", "/index.css")),
			html.T("link",
				html.A(
					"rel", "shortcut icon",
					"href", "favicon.ico"))))
}

func HEADER(active string) []byte {
	var (
		homeAttr  = html.A("href", "/")
		postsAttr = html.A("href", "/posts/")
		tagsAttr  = html.A("href", "/tags/")
	)

	if active == "/" {
		homeAttr = html.A("href", "/", "class", "active")
	} else if active == "/posts/" {
		homeAttr = html.A("href", "/posts/", "class", "active")
	} else if active == "/tags/" {
		homeAttr = html.A("href", "/tags/", "class", "active")
	}

	return html.E("header", nil,
		html.E("nav", html.A("class", "wrap"),
			html.C(
				html.E("div", nil,
					html.E("a", homeAttr, []byte("Karl McGuire"))),
				html.E("div", nil,
					html.C(
						html.E("a", postsAttr, []byte("Posts")),
						html.E("a", tagsAttr, []byte("Tags")))))))
}

func POST(title, date string, tags []string, content []byte) []byte {
	href := "/posts/" +
		strings.ReplaceAll(strings.ToLower(title), " ", "-") + "/"

	for i := range tags {
		tags[i] = string(
			html.E("a", html.A("href", "/tags/"+tags[i]+"/"),
				[]byte(tags[i])))
	}

	return html.E("main", nil,
		html.E("div", html.A("class", "wrap"),
			html.C(
				html.E("a", html.A("class", "title", "href", href),
					html.E("h1", nil, []byte(title))),
				html.E("div", html.A("class", "meta"),
					html.C(
						html.E("span", nil, []byte(date)),
						html.E("div", html.A("class", "meta__tags"),
							[]byte("("+strings.Join(tags, ", ")+")")))),
				html.E("div", html.A("class", "content"), content))))
}

func LIST(title, href string, rows [][]string) []byte {
	var data bytes.Buffer

	for _, row := range rows {
		// row[0] - "programming"
		// row[1] - "/programming/"
		// row[2] - "5 posts"
		data.Write(
			html.E("li", nil,
				html.C(
					html.E("a", html.A("href", row[1]), []byte(row[0])),
					html.E("span", nil, []byte(row[2])))))
	}

	return html.E("main", nil,
		html.E("div", html.A("class", "wrap"),
			html.C(
				html.E("a", html.A("class", "title", "href", href),
					html.E("h1", nil, []byte(title))),
				html.E("div", html.A("class", "content"),
					html.E("ul", html.A("class", "tags"), data.Bytes())))))
}

func FOOTER(year, note string) []byte {
	return html.E("footer", html.A("class", "wrap"),
		html.C(
			html.E("span", nil, []byte("&copy; "+year+" Karl McGuire")),
			html.E("span", nil, []byte(note))))
}

func ParsePost(path string) (string, string, []string, []byte) {
	var (
		data []byte
		err  error
	)

	// load post from disk
	if data, err = ioutil.ReadFile(path); err != nil {
		panic(err)
	}

	var (
		buff  = bytes.NewBuffer(data)
		title string
		date  string
		tags  string
	)

	// ignore ``` line
	if _, err = buff.ReadString('\n'); err != nil {
		panic(err)
	}

	// get title line
	if title, err = buff.ReadString('\n'); err != nil {
		panic(err)
	}

	// get date line
	if date, err = buff.ReadString('\n'); err != nil {
		panic(err)
	}

	// get tags line
	if tags, err = buff.ReadString('\n'); err != nil {
		panic(err)
	}

	// ignore ``` line
	if _, err = buff.ReadString('\n'); err != nil {
		panic(err)
	}

	// ignore previous line
	if _, err = buff.ReadString('\n'); err != nil {
		panic(err)
	}

	return title[:len(title)-1],
		date[:len(date)-1],
		strings.Split(tags[:len(tags)-1], " "),
		blackfriday.Run(buff.Bytes())
}

// GetPosts takes in a directory path and returns a string slice of path for
// each file contained in the directory path.
func GetPosts(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	posts := make([]string, len(files))
	for i, file := range files {
		posts[i] = path + file.Name()
	}

	return posts
}

// PutPosts takes a string slice of post paths (to the markdown files) and
// creates a directory for each post, with an index.html containing the
// rendered post content.
func PutPosts(paths []string) {
	for _, path := range paths {
		full := "./docs/" + path[:len(path)-3] + "/"

		// create /docs/posts/post-title/ directory
		os.MkdirAll(full, os.ModePerm)

		// parse post md
		title, date, tags, content := ParsePost(path)

		// write index.html
		ioutil.WriteFile(full+"index.html",
			PAGE(
				HEAD(title, ""),
				HEADER(""),
				POST(title, date, tags, content),
				FOOTER(YEAR, NOTE)), os.ModePerm)
	}
}

func GetTags(path string) [][]string {
	// get list of all post paths
	posts := GetPosts(path)

	// will hold all tags
	all := make(map[string]int)

	for _, post := range posts {
		// get tags of post
		_, _, tags, _ := ParsePost(post)

		for _, tag := range tags {
			all[tag]++
		}
	}

	// create [][]string from map with tag name and post count string
	out := make([][]string, 0)
	for tag, count := range all {
		out = append(out, []string{tag, fmt.Sprintf("%d posts", count)})
	}

	return out
}

func main() {
	/*
		os.Stdout.Write(
			LIST(
				"title",
				"href",
				[][]string{
					{"programming", "/programming/", "5 posts"},
					{"personal", "/personal/", "2 posts"}}))
	*/

	fmt.Println(GetTags("posts/"))
}
