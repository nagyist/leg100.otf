package organization

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/leg100/otf"
)

type Service struct {
	application
	api *api
	web *web
}

func NewService(ctx context.Context, opts Options) *Service {
	app := &app{
		Authorizer: opts.Authorizer,
		Logger:     opts.Logger,
		db:         newDB(opts.DB),
	}
	return &Service{
		application: app,
		api:         &api{app},
	}
}

type Options struct {
	otf.Authorizer
	otf.DB
	otf.Renderer
	logr.Logger
}

func (s *Service) AddHandlers(r *mux.Router) {
	s.api.AddHandlers(r)
	s.web.AddHandlers(r)
}

func (a *Service) CreateOrganization(ctx context.Context, opts otf.OrganizationCreateOptions) (otf.Organization, error) {
	return a.create(ctx, opts)
}

func (a *Service) GetOrganization(ctx context.Context, name string) (otf.Organization, error) {
	return a.get(ctx, name)
}

func (a *Service) DeleteOrganization(ctx context.Context, name string) error {
	return a.delete(ctx, name)
}