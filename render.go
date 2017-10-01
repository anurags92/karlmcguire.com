package main

import (
	"os"
	"text/template"
)

func Render(dir, tem string, data interface{}) error {
	os.Mkdir(dir, 0700)

	file, err := os.Create(dir + "index.html")
	if err != nil {
		return err
	}

	template.Must(template.ParseFiles(
		"templates/base.html",
		tem,
	)).ExecuteTemplate(file, "base", data)

	return nil
}
