package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"net/http"
	"strings"

	"github.com/leg100/otf"
)

// authMiddleware authenticates http requests, ensuring they possess a valid api
// token
type authMiddleware struct {
	otf.TokenService
}

func newAuthMiddleware(service otf.TokenService) *authMiddleware {
	return &authMiddleware{TokenService: service}
}

// Middleware can be passed to mux.Use()
func (m *authMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid bearer token format", http.StatusForbidden)
			return
		}

		token := parts[1]

		validTokens, err := m.List(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isValidToken(token, validTokens) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "invalid token", http.StatusForbidden)
		}
	})
}

func isValidToken(token string, validTokens []*otf.Token) bool {
	h := sha256.New()
	if _, err := h.Write([]byte(token)); err != nil {
		panic("producing hash of user token: " + err.Error())
	}

	for _, v := range validTokens {
		if bytes.Equal(v.Hash, h.Sum(nil)) {
			return true
		}
	}

	return false
}
