package main

import (
	"github.com/gin-contrib/multitemplate"
	"log"
	"path/filepath"
)

func (cas *cas) loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	bases, err := filepath.Glob(templatesDir + "/*.tmpl")
	if err != nil {
		log.Fatalln(err.Error())
	}

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, base := range bases {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append([]string{base}, layoutCopy...)
		r.AddFromFiles(filepath.Base(base), files...)
	}
	return r
}
