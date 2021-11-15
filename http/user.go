package http

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/leg100/otf"
)

var (
	//go:embed static
	static embed.FS

	templatesGlob = "static/templates/*.tmpl"

	tmpl = template.Must(template.ParseFS(static, templatesGlob))
)

func getTemplate() *template.Template {
	if os.Getenv("OTF_DEBUG") == "true" {
		return template.Must(template.ParseGlob(filepath.Join("http", templatesGlob)))
	}

	return tmpl
}

func cssFilesystem() http.FileSystem {
	if os.Getenv("OTF_DEBUG") == "true" {
		return http.Dir(http.StripPrefix"http/static/css")
	}

	return http.FS(static)
}

// User represents a Terraform Enterprise user.
type User struct {
	ID               string     `jsonapi:"primary,users"`
	AvatarURL        string     `jsonapi:"attr,avatar-url"`
	Email            string     `jsonapi:"attr,email"`
	IsServiceAccount bool       `jsonapi:"attr,is-service-account"`
	TwoFactor        *TwoFactor `jsonapi:"attr,two-factor"`
	UnconfirmedEmail string     `jsonapi:"attr,unconfirmed-email"`
	Username         string     `jsonapi:"attr,username"`
	V2Only           bool       `jsonapi:"attr,v2-only"`

	// Relations
	// AuthenticationTokens *AuthenticationTokens `jsonapi:"relation,authentication-tokens"`
}

type ListTokenOutput struct {
	Tokens []*otf.Token
}

// TwoFactor represents the organization permissions.
type TwoFactor struct {
	Enabled  bool `jsonapi:"attr,enabled"`
	Verified bool `jsonapi:"attr,verified"`
}

func (s *Server) CreateToken(w http.ResponseWriter, r *http.Request) {
	desc := r.FormValue("description")
	if desc == "" {
		WriteError(w, http.StatusPreconditionRequired, fmt.Errorf("missing description form input"))
		return
	}

	_, err := s.TokenService.Create(r.Context(), otf.TokenCreateOptions{Description: desc})
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	http.Redirect(w, r, "/app/settings/tokens", http.StatusMovedPermanently)
}

func (s *Server) ListTokens(w http.ResponseWriter, r *http.Request) {
	tokens, err := s.TokenService.List(r.Context())
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	if err := getTemplate().Execute(w, ListTokenOutput{Tokens: tokens}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) DeleteToken(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		WriteError(w, http.StatusPreconditionRequired, fmt.Errorf("missing id form input"))
		return
	}

	if err := s.TokenService.Delete(r.Context(), id); err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	http.Redirect(w, r, "/app/settings/tokens", http.StatusMovedPermanently)
}
