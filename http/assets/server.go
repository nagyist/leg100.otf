package assets

import (
	"html/template"
	"net/http"
)

// Server provides the means to retrieve http assets (templates and static
// files such as CSS).
type Server interface {
	GetTemplate(string) *template.Template
	GetStaticFS() http.FileSystem
}

type ServerOptions struct {
	// LayoutPath is the relative path to the layout template
	LayoutPath string

	// ContentPath is the relative path to the directory containing content
	// templates
	ContentPath string

	// StylesheetPath is the relative path to the directory containing CSS files
	StylesheetPath string
}
