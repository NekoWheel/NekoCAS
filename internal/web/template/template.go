package template

import (
	"html/template"
	"time"
)

func FuncMap() []template.FuncMap {
	return []template.FuncMap{map[string]interface{}{
		"Year": func() int {
			return time.Now().Year()
		},
		"Safe": func(raw string) template.HTML {
			return Safe(raw)
		},
	}}
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}
