package otf

import (
	"context"
	"crypto/sha256"
	"fmt"
)

type TokenService interface {
	Create(ctx context.Context, opts TokenCreateOptions) (*Token, string, error)
	Get(ctx context.Context, id string) (*Token, error)
	List(ctx context.Context) ([]*Token, error)
	Delete(ctx context.Context, id string) error
}

// TokenStore is a store storing API tokens
type TokenStore interface {
	Create(ctx context.Context, token Token) error
	Get(ctx context.Context, id string) (*Token, error)
	List(ctx context.Context) ([]*Token, error)
	Delete(ctx context.Context, id string) error
}

// Token is the metadata for a user's API token
type Token struct {
	Timestamps

	ID          string `db:"token_id"`
	Description string

	// SHA-256 hash of token string
	Hash []byte
}

// TokenCreateOptions represents the options for creating a new token.
type TokenCreateOptions struct {
	// Type is a public field utilized by JSON:API to set the resource type via
	// the field tag.  It is not a user-defined value and does not need to be
	// set.  https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,authentication-tokens"`

	Description string `jsonapi:"attr,description"`

	// Optional. Provide hardcoded token. Skips generation of a random secure
	// string.
	Token *string
}

// TokenListOptions are the options for listing tokens.
type TokenListOptions struct {
	// Optional. Order the listing by a specific field.
	OrderBy *string
}

// NewToken constructs a new Token object. Returns the Token object, along with
// the secure string itself.
func NewToken(opts TokenCreateOptions) (Token, string, error) {
	var secureStr string

	if opts.Token != nil {
		secureStr = *opts.Token
	} else {
		secureStr = GenerateRandomString(32)
	}

	h := sha256.New()
	if _, err := h.Write([]byte(secureStr)); err != nil {
		return Token{}, "", fmt.Errorf("producing hash of secure string: %w", err)
	}

	token := Token{
		ID:          NewID("at"),
		Timestamps:  NewTimestamps(),
		Description: opts.Description,
		Hash:        h.Sum(nil),
	}

	return token, secureStr, nil
}
