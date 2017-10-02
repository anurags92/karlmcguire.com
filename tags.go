package main

import (
	"sort"
)

type Tags []*Tag

func (t Tags) Len() int {
	return len(t)
}

func (t Tags) Less(i, j int) bool {
	return len(t[i].Posts) > len(t[j].Posts)
}

func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Tags) Sort() {
	sort.Sort(t)
}

func (t Tags) Render(dir string) error {
	t.Sort()

	err := Render(dir+"tags/", "templates/tags.html",
		&struct {
			Title    string
			Selected string
			Tags     Tags
		}{"Tags", "tags", t},
	)
	if err != nil {
		return err
	}

	for _, tag := range t {
		err := tag.Render(dir + "tags/")
		if err != nil {
			return err
		}
	}

	return nil
}
