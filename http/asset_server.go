package http

import (
	"embed"
	"net/http"
	"path/filepath"
	"text/template"
)

var (
	//go:embed static
	static embed.FS

	templatesGlob = "static/templates/*.tmpl"

	templates = template.Must(template.ParseFS(static, templatesGlob))
)

// AssetServer provides the means to retrieve http assets (templates and static
// files such as CSS).
type AssetServer interface {
	GetTemplates() *template.Template
	GetStaticFS() http.FileSystem
}

// EmbeddedAssetServer embeds assets within the go binary (the default).
type EmbeddedAssetServer struct{}

func (s *EmbeddedAssetServer) GetTemplates() *template.Template {
	return templates
}

func (s *EmbeddedAssetServer) GetStaticFS() http.FileSystem {
	return http.FS(static)
}

// DevAssetServer reads assets from developer's machine, permitting use of
// something like livereload to see changes in real-time.
type DevAssetServer struct{}

func (s *DevAssetServer) GetTemplates() *template.Template {
	return template.Must(template.ParseGlob(filepath.Join("http", templatesGlob)))
}

func (s *DevAssetServer) GetStaticFS() http.FileSystem {
	return http.Dir("http")
}
