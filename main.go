package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xba/html"
	"gopkg.in/russross/blackfriday.v2"
)

const (
	YEAR = "2019"
	NOTE = `Made with <a href="https://golang.org/">Go</a>.`

	SCRIPTS = `
	<script data-no-instant>InstantClick.init();</script>`
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
						footer,
						[]byte(SCRIPTS)))))), []byte{'\x00'}, []byte(""))
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
			// code highlight
			/*html.T("link",
				html.A(
					"rel", "stylesheet",
					"href", "/code.css")),
			html.E("script", html.A("src", "/highlight.pack.js"), nil),
			html.E("script", nil,
				[]byte("hljs.initHighlightingOnLoad();")),*/
			// instant click
			html.E("script", html.A("src", "/ic.min.js"), nil),
			// roboto font
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
					"href", "/favicon.ico"))))
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
		postsAttr = html.A("href", "/posts/", "class", "active")
	} else if active == "/tags/" {
		tagsAttr = html.A("href", "/tags/", "class", "active")
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

func POST(path, title, date string, tags []string, content []byte) []byte {
	// href := "/posts/" +
	// 	strings.ReplaceAll(strings.ToLower(title), " ", "-") + "/"
	href := "/" + path[:len(path)-3] + "/"

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

func ParsePost(path string) (string, string, string, []string, []byte) {
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

	return path,
		title[:len(title)-1],
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

// GetPostTitles takes in a directory path and returns a 2d string slice of
// post titles, dates, and hrefs. Ordered in descending order.
func GetPostTitles(path string) [][]string {
	posts := GetPosts(path)

	titles := make([][]string, 0)
	for _, post := range posts {
		_, title, date, _, _ := ParsePost(post)

		titles = append(titles,
			[]string{title, "/posts/" + post[:len(post)-3] + "/", date})
	}

	// sort in descending order
	sort.Sort(SortablePosts(titles))

	return titles
}

// GetTags takes in a directy path and returns a string slice of tag names and
// the number of posts with that tag.
func GetTags(path string) [][]string {
	// get list of all post paths
	posts := GetPosts(path)

	// will hold all tags
	all := make(map[string]int)

	for _, post := range posts {
		// get tags of post
		_, _, _, tags, _ := ParsePost(post)
		for _, tag := range tags {
			// keep count of how many posts are associated with the tag
			all[tag]++
		}
	}

	// create [][]string from map with tag name and post count string
	out := make([][]string, 0)
	for tag, count := range all {
		out = append(out,
			[]string{tag, "/tags/" + tag + "/", fmt.Sprintf("%d posts", count)})
	}

	// sort in descending order
	sort.Sort(SortableTags(out))

	return out
}

// PutPosts takes a string slice of post paths (to the markdown files) and
// creates a directory for each post, with an index.html containing the
// rendered post content.
func PutPosts(paths []string) {
	postList := make([][]string, 0)

	for _, path := range paths {
		full := "./docs/" + path[:len(path)-3] + "/"

		// create /docs/posts/post-title/ directory
		os.MkdirAll(full, os.ModePerm)

		// parse post md
		_, title, date, tags, content := ParsePost(path)

		// add the title, href, and date to postList
		postList = append(postList,
			[]string{title, "/" + path[:len(path)-3] + "/", date})

		// write index.html
		ioutil.WriteFile(full+"index.html",
			PAGE(
				HEAD(title, ""),
				HEADER(""),
				POST(path, title, date, tags, content),
				FOOTER(YEAR, NOTE)), os.ModePerm)
	}

	// sort post list descending
	sort.Sort(SortablePosts(postList))

	// /posts/
	os.MkdirAll("./docs/posts/", os.ModePerm)

	// /posts/index.html
	ioutil.WriteFile("./docs/posts/index.html",
		PAGE(
			HEAD("Posts", ""),
			HEADER("/posts/"),
			LIST(
				// h1
				"Posts"+fmt.Sprintf(" (%d)", len(paths)),
				// h1 href
				"/posts/",
				// rows
				postList),
			FOOTER(YEAR, NOTE)), os.ModePerm)

}

func PutTags(tags [][]string) {
	// /tags/
	os.MkdirAll("./docs/tags/", os.ModePerm)

	// /tags/index.html
	ioutil.WriteFile("./docs/tags/index.html",
		PAGE(
			HEAD("Tags", ""),
			HEADER("/tags/"),
			LIST("Tags"+fmt.Sprintf(" (%d)", len(tags)), "/tags/", tags),
			FOOTER(YEAR, NOTE)), os.ModePerm)

	for _, tag := range tags {
		full := "./docs/tags/" + tag[0] + "/"

		// rows will contain the post names and dates of each post associated
		// with the current tag
		rows := make([][]string, 0)

		// populate rows
		posts := GetPosts("posts/")
		for _, post := range posts {
			_, postTitle, postDate, postTags, _ := ParsePost(post)

			for _, postTag := range postTags {
				// post is associated with the current tag
				if postTag == tag[0] {
					rows = append(rows,
						[]string{
							postTitle,
							"/" + post[:len(post)-3] + "/",
							postDate})
				}
			}
		}

		// sort tag posts by date
		sort.Sort(SortablePosts(rows))

		// /tags/tag/
		os.MkdirAll(full, os.ModePerm)

		// /tags/tag/index.html
		ioutil.WriteFile(full+"index.html",
			PAGE(
				HEAD(tag[0], ""),
				HEADER(""),
				LIST(
					// h1
					"Tags / "+tag[0]+fmt.Sprintf(" (%d)", len(rows)),
					// h1 href
					"/tags/"+tag[0]+"/",
					// list of posts using the tag
					rows),
				FOOTER(YEAR, NOTE)), os.ModePerm)
	}
}

func PutIndex(path string) {
	// read index md file from disk
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// parse markdown into []byte
	data := blackfriday.Run(file)

	ioutil.WriteFile("./docs/index.html",
		PAGE(
			HEAD("Karl McGuire", ""),
			HEADER("/"),
			html.E("main", nil,
				html.E("div", html.A("class", "wrap"),
					html.C(
						html.E("div", html.A("class", "title"),
							html.E("h1", nil, []byte("Hello,"))),
						html.E("div", html.A("class", "content"), data)))),
			FOOTER(YEAR, NOTE)), os.ModePerm)
}

type Sortable struct {
	data [][]string
	len  func([][]string) int
	less func([][]string, int, int) bool
	swap func([][]string, int, int)
}

func (s *Sortable) Len() int           { return s.len(s.data) }
func (s *Sortable) Less(i, j int) bool { return s.less(s.data, i, j) }
func (s *Sortable) Swap(i, j int)      { s.swap(s.data, i, j) }

func SortableTags(tags [][]string) sort.Interface {
	return &Sortable{
		data: tags,

		len:  func(d [][]string) int { return len(d) },
		swap: func(d [][]string, i, j int) { d[i], d[j] = d[j], d[i] },

		less: func(d [][]string, i, j int) bool {
			// get the int count of i tag
			ii, err := strconv.Atoi(d[i][2][:strings.Index(d[i][2], " ")])
			if err != nil {
				panic(err)
			}

			// get the int count of j tag
			ji, err := strconv.Atoi(d[j][2][:strings.Index(d[j][2], " ")])
			if err != nil {
				panic(err)
			}

			// descending order
			return ii > ji
		},
	}
}

func SortablePosts(posts [][]string) sort.Interface {
	return &Sortable{
		data: posts,

		len:  func(d [][]string) int { return len(d) },
		swap: func(d [][]string, i, j int) { d[i], d[j] = d[j], d[i] },

		less: func(d [][]string, i, j int) bool {
			format := "January 2, 2006"

			// get the date of the i post
			it, err := time.Parse(format, d[i][2])
			if err != nil {
				panic(err)
			}

			// get the date of the j post
			jt, err := time.Parse(format, d[j][2])
			if err != nil {
				panic(err)
			}

			// newest articles up top
			return jt.Before(it)
		},
	}
}

func main() {
	var err error

	if err = os.Remove("docs/index.html"); err != nil {
		panic(err)
	}

	if err = os.RemoveAll("docs/posts/"); err != nil {
		panic(err)
	}

	if err = os.RemoveAll("docs/tags/"); err != nil {
		panic(err)
	}

	posts := GetPosts("posts/")

	// /index.html
	PutIndex("./index.md")
	// /posts/
	PutPosts(posts)
	// /tags/
	PutTags(GetTags("posts/"))
}
