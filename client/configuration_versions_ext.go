package client

import (
	"context"
)

var _ ConfigurationVersionsExt = (*configurationVersionsExt)(nil)

type ConfigurationVersionsExt interface {
	Download(ctx context.Context, id string) ([]byte, error)
}

// configurationVersionsExt implements ConfigurationVersionsExt.
type configurationVersionsExt struct {
	client Client
}

func (s *configurationVersionsExt) Download(ctx context.Context, id string) ([]byte, error) {
	return nil, nil
}
