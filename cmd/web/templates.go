package main

import (
	"html/template"
	"path/filepath"

	"github.com/empty-crayon/snippetbox/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
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

		ts, err := template.ParseFiles("./ui/html/base.tmpl")

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
