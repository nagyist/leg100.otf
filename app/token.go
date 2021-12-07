package app

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/leg100/otf"
)

var _ otf.TokenService = (*TokenService)(nil)

type TokenService struct {
	db otf.TokenStore

	logr.Logger
}

func NewTokenService(db otf.TokenStore, logger logr.Logger) *TokenService {
	return &TokenService{
		db:     db,
		Logger: logger,
	}
}

func (s TokenService) Create(ctx context.Context, opts otf.TokenCreateOptions) (*otf.Token, string, error) {
	token, secret, err := otf.NewToken(opts)
	if err != nil {
		s.Error(err, "constructing token")
		return nil, "", err
	}

	if err := s.db.Create(ctx, token); err != nil {
		s.Error(err, "creating token")
		return nil, "", err
	}

	s.V(2).Info("created token", "id", token.ID)

	return &token, secret, nil
}

func (s TokenService) Get(ctx context.Context, id string) (*otf.Token, error) {
	token, err := s.db.Get(ctx, id)
	if err != nil {
		s.Error(err, "retrieving token")
		return nil, err
	}

	s.V(2).Info("retrieved token", "id", token.ID)

	return token, nil
}

func (s TokenService) List(ctx context.Context) ([]*otf.Token, error) {
	tokens, err := s.db.List(ctx)
	if err != nil {
		s.Error(err, "listing tokens")
		return nil, err
	}

	s.V(2).Info("listed tokens", "count", len(tokens))

	return tokens, nil
}

func (s TokenService) Delete(ctx context.Context, id string) error {
	if err := s.db.Delete(ctx, id); err != nil {
		s.Error(err, "deleting token")
		return err
	}

	s.V(2).Info("deleted token", "id", id)

	return nil
}
