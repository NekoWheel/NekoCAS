package template

import (
	"html/template"
	"time"

	"gopkg.in/macaron.v1"
)

func RenderOptions() macaron.RenderOptions {
	return macaron.RenderOptions{
		Directory:  "templates",
		IndentJSON: macaron.Env != macaron.PROD,
		Funcs: []template.FuncMap{map[string]interface{}{
			"Year": func() int {
				return time.Now().Year()
			},
			"Safe": func(raw string) template.HTML {
				return Safe(raw)
			},
		}},
	}
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}
