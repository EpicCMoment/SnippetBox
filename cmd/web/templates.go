package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.ariffil.com/internal/models"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
	Form        any
	Flash       string
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		pageName := filepath.Base(page)

		ts, err := template.New(pageName).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)

		if err != nil {
			return nil, err
		}

		cache[pageName] = ts

	}

	return cache, nil

}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
