package assets

import (
	"html/template"
	"net/http"
)

var (
	// Paths to assets in relative to the package director
	layoutTemplatePath  = "templates/layout.tmpl"
	contentTemplatesDir = "templates/content"
	stylesheetDir       = "css"
)

// Server provides the means to retrieve http assets (templates and static files
// such as CSS).
type Server interface {
	GetTemplate(string) *template.Template
	GetStaticFS() http.FileSystem
	Links() []string
}

type LayoutTemplateOptions struct {
	Title         string
	Stylesheets   []string
	FlashMessages []template.HTML
}
