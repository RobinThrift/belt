//go:build dev
// +build dev

package app

import (
	_ "embed"
	"html/template"
)

//go:embed templates/root_dev.tmpl.html
var rootTemplateRaw string

var rootTemplate = template.Must(template.New("root.html").Parse(rootTemplateRaw))

//go:embed templates/error.tmpl.html
var errorPageTemplateRaw string

var errorPageTemplate = template.Must(template.New("error.html").Parse(errorPageTemplateRaw))

//go:embed templates/manifest.tmpl.json
var manifestJSONTemplateRaw string

//nolint:gochecknoglobals
var manifestJSONTemplate = template.Must(template.New("manifest.json").Parse(manifestJSONTemplateRaw))
