package agentservice

import (
	"context"
	"errors"

	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/logr"
	"github.com/leg100/otf/internal/organization"
	"github.com/leg100/otf/internal/pubsub"
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
		ListAgents(ctx context.Context) ([]*Agent, error)
		// ListOrganizationAgents lists agents by organization
		ListOrganizationAgents(ctx context.Context, organization string) ([]*Agent, error)
		// GetJob retrieves a job assigned to the given agent. The
		// returned channel only receives at most job before being closed.
		GetJob(ctx context.Context, agentID string) (<-chan *Job, error)
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
		pubsub.Subscriber

		db           *pgdb
		site         internal.Authorizer
		organization internal.Authorizer
	}
)

func NewService(opts Options) *service {
	return &service{
		db:           &pgdb{opts.DB},
		site:         &internal.SiteAuthorizer{Logger: opts.Logger},
		organization: &organization.Authorizer{Logger: opts.Logger},
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

func (s *service) ListAgents(ctx context.Context) ([]*Agent, error) {
	subject, err := s.site.CanAccess(ctx, rbac.ListAgentsAction, "")
	if err != nil {
		return nil, err
	}
	agents, err := s.db.list(ctx)
	if err != nil {
		s.Error(err, "listing agents", "subject", subject)
		return nil, err
	}
	s.V(4).Info("listed agents", "subject", subject)
	return agents, nil
}

func (s *service) ListOrganizationAgents(ctx context.Context, organization string) ([]*Agent, error) {
	subject, err := s.organization.CanAccess(ctx, rbac.ListAgentsAction, organization)
	if err != nil {
		return nil, err
	}
	agents, err := s.db.listByOrganization(ctx, organization)
	if err != nil {
		s.Error(err, "listing agents", "subject", subject, "organization", organization)
		return nil, err
	}
	s.V(9).Info("listed agents", "subject", subject, "organization", organization)
	return agents, nil
}

func (s *service) GetJob(ctx context.Context, agentID string) (<-chan *Job, error) {
	// unsubscribe whenever exiting this routine.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// subscribe to job events
	sub, err := s.Subscribe(ctx, "get-job-")
	if err != nil {
		return nil, err
	}
	// retrieve job assigned job, if any
	job, err := s.db.getAssignedJobByAgentID(ctx, agentID)
	if err != nil && !errors.Is(err, internal.ErrResourceNotFound) {
		// unexpected error
		return nil, err
	}
	ch := make(chan *Job)
	go func() {
		if job == nil {
			// no existing assigned job; wait til there is one
			for event := range sub {
				var ok bool
				job, ok = event.Payload.(*Job)
				if !ok {
					continue // skip anything other than a job
				}
				if job.Status != AssignedJob {
					continue // skip jobs in other states
				}
				if job.AgentID == nil {
					// should never happen
					s.Error(nil, "assigned job missing agent ID", "job_id", job.ID)
					continue
				}
				if *job.AgentID != agentID {
					continue // skip jobs assigned to other agents
				}
				break
			}
		}
		// wait til either context is done or the caller receives the job
		select {
		case <-ctx.Done():
		case ch <- job:
		}
		// either way, close the channel
		close(ch)
	}()
	return ch, nil
}
