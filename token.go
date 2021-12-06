package otf

import "context"

type TokenService interface {
	Create(ctx context.Context, opts TokenCreateOptions) (string, error)
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

// Token is the metadata for a user's API token (but not the token itself)
type Token struct {
	Timestamps

	ID          string `db:"token_id"`
	Description string
}

// TokenCreateOptions represents the options for creating a new token.
type TokenCreateOptions struct {
	// Type is a public field utilized by JSON:API to set the resource type via
	// the field tag.  It is not a user-defined value and does not need to be
	// set.  https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,authentication-tokens"`

	Description string `jsonapi:"attr,description"`
}

// TokenListOptions are the options for listing tokens.
type TokenListOptions struct {
	// Optional. Order the listing by a specific field.
	OrderBy *string
}

func NewToken(opts TokenCreateOptions) (string, Token) {
	token := Token{
		ID:          NewID("at"),
		Timestamps:  NewTimestamps(),
		Description: opts.Description,
	}

	return GenerateRandomString(32), token
}
