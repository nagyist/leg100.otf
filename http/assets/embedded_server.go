package assets

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

// EmbeddedServer provides access to assets embedded in the go binary.
type EmbeddedServer struct {
	// templates maps template names to parsed contents
	templates map[string]*template.Template
}

func NewEmbeddedServer(filesystem embed.FS, opts ServerOptions) (*EmbeddedServer, error) {
	pattern := fmt.Sprintf("%s/*.tmpl", opts.ContentPath)
	paths, err := fs.Glob(filesystem, pattern)
	if err != nil {
		return nil, fmt.Errorf("unable to read embedded templates directory: %w", err)
	}

	server := EmbeddedServer{
		templates: make(map[string]*template.Template, len(paths)),
	}

	for _, p := range paths {
		contentPath := filepath.Join(contentDir, p.Name())

		templates[p.Name()] = template.Must(template.ParseFS(static, layoutPath, contentPath))
	}
	return &EmbeddedServer{templates: templates}
}

func (s *EmbeddedServer) GetTemplate(name string) *template.Template {
	return s.templates[name]
}

func (s *EmbeddedServer) GetStaticFS() http.FileSystem {
	return http.FS(static)
}
