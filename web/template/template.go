package template

import "html/template"

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}
