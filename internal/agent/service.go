package agent

import "context"

const (
	Busy    Status = "busy"
	Idle    Status = "idle"
	Exited  Status = "exited"
	Errored Status = "errored"
	Unknown Status = "unknown"
)

type (
	// Service for server-side management of agents
	Service interface {
		// Register an agent and return unique ID
		Register(ctx context.Context, opts RegisterOptions) (string, error)
		// UpdateStatus updates the status of an agent with the given ID.
		UpdateStatus(ctx context.Context, id string, status Status) error
	}

	RegisterOptions struct {
		Name *string // optional agent name
	}

	Status string

	service struct{}
)

func (s *service) Register(ctx context.Context, opts RegisterOptions) (string, error) {
	return "", nil
}

func (s *service) UpdateStatus(ctx context.Context, id string, status Status) error {
	return nil
}
