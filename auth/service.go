package auth

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/leg100/otf"
	"github.com/leg100/otf/cloud"
	"github.com/leg100/otf/organization"
)

type Service struct {
	app

	api *api
	web *web
}

func NewService(ctx context.Context, opts Options) (*Service, error) {
	db := newDB(opts.DB, opts.Logger)
	app := &Application{
		Authorizer: opts.Authorizer,
		Logger:     opts.Logger,
		db:         db,
	}
	app.synchroniser = &synchroniser{opts.Logger, opts.Service, app}

	authenticators, err := newAuthenticators(authenticatorOptions{
		Logger:          opts.Logger,
		HostnameService: opts.HostnameService,
		app:             app,
		configs:         opts.Configs,
	})
	if err != nil {
		return nil, err
	}

	api := &api{
		app: app,
	}
	web := &web{
		Renderer:       opts.Renderer,
		app:            app,
		authenticators: authenticators,
		siteToken:      opts.SiteToken,
	}

	// purge expired sessions
	go db.startExpirer(ctx, defaultExpiry)

	return &Service{
		app: app,
		api: api,
		web: web,
	}, nil
}

type Options struct {
	Configs   []*cloud.CloudOAuthConfig
	SiteToken string

	organization.Service
	otf.Authorizer
	otf.DB
	otf.Renderer
	otf.HostnameService
	logr.Logger
}

func (a *Service) AddHandlers(r *mux.Router) {
	a.api.addHandlers(r)
	a.web.addHandlers(r)
}