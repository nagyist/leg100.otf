package http

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

var (
	//go:embed static
	static embed.FS

	cssDir     = "static/css"
	layoutPath = "static/templates/layout.tmpl"
	contentDir = "static/templates/content"

	templates map[string]*template.Template
)

// Parse embedded templates at startup
func init() {
	entries, err := static.ReadDir(contentDir)
	if err != nil {
		panic(fmt.Sprintf("unable to read embedded templates directory: %s", err.Error()))
	}

	templates = make(map[string]*template.Template, len(entries))

	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}

		contentPath := filepath.Join(contentDir, ent.Name())

		templates[ent.Name()] = template.Must(template.ParseFS(static, layoutPath, contentPath))
	}
}

type LayoutTemplateOptions struct {
	Stylesheets   []string
	FlashMessages []string
}

func NewLayoutTemplateOptions(server AssetServer, r *http.Request, w http.ResponseWriter) LayoutTemplateOptions {
	session, _ := store.Get(r, "flash")
	session.Save(r, w)

	return LayoutTemplateOptions{
		FlashMessages: interfaceSliceToStringSlice(session.Flashes()),
	}
}

// AssetServer provides the means to retrieve http assets (templates and static
// files such as CSS).
type AssetServer interface {
	GetTemplate(string) *template.Template
	GetStaticFS() http.FileSystem
}

// EmbeddedAssetServer embeds assets within the go binary (the default).
type EmbeddedAssetServer struct {
	templates map[string]*template.Template
}

func NewEmbeddedAssetServer() *EmbeddedAssetServer {
	return &EmbeddedAssetServer{templates: templates}
}

func (s *EmbeddedAssetServer) GetTemplate(name string) *template.Template {
	return s.templates[name]
}

func (s *EmbeddedAssetServer) GetStaticFS() http.FileSystem {
	return http.FS(static)
}

// DevAssetServer reads assets from developer's machine, permitting use of
// something like livereload to see changes in real-time.
type DevAssetServer struct{}

func (s *DevAssetServer) GetTemplate(name string) *template.Template {
	layoutPath := filepath.Join("http", layoutPath)
	contentPath := filepath.Join("http", contentDir, name)

	return template.Must(template.ParseFiles(layoutPath, contentPath))
}

func (s *DevAssetServer) GetStaticFS() http.FileSystem {
	return http.Dir("http")
}
