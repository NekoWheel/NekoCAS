package templates

import (
	"embed"
)

//go:embed layouts mail manage *.tmpl
var FS embed.FS
