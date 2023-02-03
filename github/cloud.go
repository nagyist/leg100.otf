package github

import (
	"context"
	"net/http"

	"github.com/leg100/otf/cloud"
)

type Cloud struct{}

func (g *Cloud) NewClient(ctx context.Context, opts cloud.ClientOptions) (cloud.Client, error) {
	return NewClient(ctx, opts)
}

func (Cloud) NewHandler(opts cloud.HandlerOptions) http.Handler {
	return &Handler{opts}
}
