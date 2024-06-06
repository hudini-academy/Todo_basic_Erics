package main

import (
	"todo/pkg/models"
	"html/template"
	"path/filepath"

	"todo/pkg/forms"
)

// templateData is used to pass data to the template.
type templateData struct {
	Snippets []*models.Todo
	Form *forms.Form
	Flash string
	AuthenticatedUser int
	CSRFToken string
}

// newTemplateCache creates a new template cache.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get all the page files.
	pages, errFetchingPages := filepath.Glob(filepath.Join(dir,"*.page.tmpl"))
	if errFetchingPages != nil {
        return nil, errFetchingPages
    }

	// For each page file add partials to it.
	for _, page := range pages {
		name := filepath.Base(page)

		ts, errParsingPage := template.ParseFiles(page)
		if errParsingPage != nil {
            return nil, errParsingPage
        }

		ts, errParsingPage = ts.ParseGlob(filepath.Join(dir,"*.partial.tmpl"))
		if errParsingPage!= nil {
            return nil, errParsingPage
        }

		cache[name] = ts;
	}
	return cache, nil;
}