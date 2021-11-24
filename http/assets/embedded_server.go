package assets

import (
	"html/template"
	"net/http"
)

// EmbeddedServer embeds assets within the go binary (the default).
type EmbeddedServer struct {
	templates map[string]*template.Template
}

func NewEmbeddedServer() *EmbeddedServer {
	return &EmbeddedServer{templates: templates}
}

func (s *EmbeddedServer) GetTemplate(name string) *template.Template {
	return s.templates[name]
}

func (s *EmbeddedServer) GetStaticFS() http.FileSystem {
	return http.FS(static)
}
