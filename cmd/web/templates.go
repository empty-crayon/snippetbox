package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/empty-crayon/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// custom tempalate functions can accept as many articles as they need to but must return one value only (although returning error as second value is allowed)
var functions = template.FuncMap{
	"humanDate" : humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// initialising a new map to act as the cache
	cache := map[string]*template.Template{}

	// gives a slice of pages which all match the given wildcard pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract file name like base.tmpl, etc
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the tempalate set before we 
		// call parsefiles()
		// hence, we have to use template.new() to create an empty template set
		// use the funcs() method to register the funcmap and then parse the files as we do
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// adding template set to the map, using the name of the page
		// like using "home.tmpl" as the key
		cache[name] = ts
	}
	return cache, nil
}
