package assets

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// DevServer reads assets from developer's machine, permitting use of
// something like livereload to see changes in real-time.
type DevServer struct{}

func (s *DevServer) GetTemplate(name string) *template.Template {
	layoutPath := filepath.Join("http", layoutPath)
	contentPath := filepath.Join("http", contentDir, name)

	return template.Must(template.ParseFiles(layoutPath, contentPath))
}

func (s *DevServer) GetStaticFS() http.FileSystem {
	return http.Dir("http")
}
