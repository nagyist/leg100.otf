package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

const (
	secret = "Qi5leXDQ/0pwIHBrL+QHJNwhHxpsQ267kl4V9uqv6uQ="

	loginPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`
)

var (
	secureCookie      = securecookie.New(secret, nil)
	usernamePasswords = map[string]string{
		"louis": "password",
	}
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

// TwoFactor represents the organization permissions.
type TwoFactor struct {
	Enabled  bool `jsonapi:"attr,enabled"`
	Verified bool `jsonapi:"attr,verified"`
}

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, loginPage)
}

func (s *Server) UserLoginSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("pass")

	pass, ok := usernamePasswords[name]

}
