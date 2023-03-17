package auth

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/leg100/otf"
	"github.com/leg100/otf/cloud"
	"github.com/leg100/otf/organization"
)

type (
	AuthService interface {
		AgentTokenService
		RegistrySessionService
		sessionService
		teamService
		tokenService
		UserService
	}

	service struct {
		logr.Logger
		TokenMiddleware, SessionMiddleware mux.MiddlewareFunc

		*synchroniser

		api          *api
		db           *pgdb
		organization otf.Authorizer
		web          *webHandlers
	}

	Options struct {
		Configs   []*cloud.CloudOAuthConfig
		SiteToken string

		OrganizationService
		otf.DB
		otf.Renderer
		otf.HostnameService
		logr.Logger
	}

	OrganizationService organization.Service
)

func NewService(ctx context.Context, opts Options) (*service, error) {
	svc := service{Logger: opts.Logger}
	svc.TokenMiddleware = AuthenticateToken(&svc)
	svc.SessionMiddleware = AuthenticateSession(&svc)

	authenticators, err := newAuthenticators(authenticatorOptions{
		Logger:          opts.Logger,
		HostnameService: opts.HostnameService,
		AuthService:     &svc,
		configs:         opts.Configs,
	})
	if err != nil {
		return nil, err
	}

	db := newDB(opts.DB, opts.Logger)
	// purge expired sessions
	go db.startExpirer(ctx, defaultExpiry)

	svc.synchroniser = &synchroniser{opts.Logger, opts.OrganizationService, &svc}
	svc.api = &api{svc: &svc}
	svc.db = db
	svc.organization = &organization.Authorizer{opts.Logger}
	svc.web = &webHandlers{
		Renderer:       opts.Renderer,
		svc:            &svc,
		authenticators: authenticators,
		siteToken:      opts.SiteToken,
	}

	return &svc, nil
}

func (a *service) AddHandlers(r *mux.Router) {
	a.api.addHandlers(r)
	a.web.addHandlers(r)
}
