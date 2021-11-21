package http

import (
	"fmt"
	"net/http"

	"github.com/leg100/otf"
)

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
	Stylesheets []string
	Tokens      []*otf.Token
	Flashes     []string
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

	token, err := s.TokenService.Create(r.Context(), otf.TokenCreateOptions{Description: desc})
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return
	}

	session, _ := store.Get(r, "fmessages")
	session.AddFlash(fmt.Sprintf("created token: %s", token))
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	session, _ := store.Get(r, "flash")

	output := ListTokenOutput{
		Tokens:  tokens,
		Flashes: interfaceSliceToStringSlice(session.Flashes()),
	}

	session.Save(r, w)

	if err := s.GetTemplate("tokens_list.tmpl").Execute(w, output); err != nil {
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

	session, _ := store.Get(r, "flash")
	session.AddFlash("deleted token")
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/app/settings/tokens", http.StatusMovedPermanently)
}

func interfaceSliceToStringSlice(is []interface{}) (ss []string) {
	for _, i := range is {
		ss = append(ss, i.(string))
	}
	return ss
}
