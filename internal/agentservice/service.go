package agentservice

import (
	"context"

	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/logr"
	"github.com/leg100/otf/internal/rbac"
)

type (
	AgentService = Service

	// Service for server-side management of agents
	Service interface {
		// Register an agent and return unique ID
		Register(ctx context.Context, opts RegisterOptions) (string, error)
		// UpdateStatus updates the status of an agent with the given ID.
		UpdateStatus(ctx context.Context, id string, status Status) error
		// ListAgents lists agents
		ListAgents(ctx context.Context, opts ListAgentsOptions) ([]*Agent, error)
	}

	ListAgentsOptions struct {
		Organization *string // Optionally filter agents by organization
	}

	RegisterOptions struct {
		IPAddress    string
		Version      string
		Name         *string // optional agent name
		External     bool
		Organization *string
	}

	// Options for constructing new service
	Options struct {
		internal.DB
		logr.Logger
	}

	service struct {
		logr.Logger

		db   *pgdb
		site internal.Authorizer
	}
)

func NewService(opts Options) *service {
	return &service{
		db:   &pgdb{opts.DB},
		site: &internal.SiteAuthorizer{Logger: opts.Logger},
	}
}

func (s *service) Register(ctx context.Context, opts RegisterOptions) (string, error) {
	_, err := s.site.CanAccess(ctx, rbac.RegisterAgentAction, "")
	if err != nil {
		return "", err
	}
	agent, err := newAgent(opts)
	if err != nil {
		s.Error(err, "registering agent")
		return "", err
	}
	if err := s.db.create(ctx, agent); err != nil {
		s.Error(err, "registering agent")
		return "", err
	}
	s.Info("registered agent", "agent", agent)
	return agent.ID, nil
}

func (s *service) UpdateStatus(ctx context.Context, id string, status Status) error {
	_, err := s.site.CanAccess(ctx, rbac.UpdateAgentStatusAction, "")
	if err != nil {
		return err
	}
	if err := s.db.updateStatus(ctx, id, status); err != nil {
		s.Error(err, "updating agent status", "id", id)
		return err
	}
	s.V(4).Info("updated agent status", "id", id)
	return nil
}

func (s *service) ListAgents(ctx context.Context, opts ListAgentsOptions) ([]*Agent, error) {
	_, err := s.site.CanAccess(ctx, rbac.ListAgentsAction, "")
	if err != nil {
		return nil, err
	}
	if err := s.db.updateStatus(ctx, id, status); err != nil {
		s.Error(err, "updating agent status", "id", id)
		return nil, err
	}
	s.V(4).Info("updated agent status", "id", id)
	return nil
}
