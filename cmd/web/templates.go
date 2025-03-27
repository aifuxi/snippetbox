package main

import (
	"github.com/aifuxi/snippetbox/internal/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{"./ui/html/base.tmpl", "./ui/html/partials/nav.tmpl", page}

		// 将一组template解析成1个template
		ts, err2 := template.ParseFiles(files...)
		if err2 != nil {
			return nil, err2
		}

		// 将解析好后的template放入缓存
		cache[name] = ts
	}

	return cache, nil
}
