package main

type TagPost struct {
	Title string
	Date  string
	Path  string
}

type Tag struct {
	Name  string
	Posts []*TagPost
}

func (t *Tag) Render(dir string) error {
	return Render(dir+t.Name+"/", "templates/tag.html",
		&struct {
			Title    string
			Selected string
			Tag      *Tag
		}{"Tag - " + t.Name, " ", t},
	)
}
