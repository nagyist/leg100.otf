package jobs

import (
	"context"

	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/logr"
)

const (
	Created   Status = "created"
	Assigned  Status = "assigned"
	Running   Status = "running"
	Completed Status = "completed"
	Errored   Status = "errored"
)

type (
	Status string

	JobService = Service

	Service interface {
		CreateJob(context.Context, CreateJobOptions) error
		AssignJob(context.Context, AssignJobOptions) error
		GetJob(context.Context, GetJobOptions) (*Job, error)
		UpdateJobStatus(context.Context, Status) error
	}

	CreateJobOptions struct {
		RunID string
		Phase internal.PhaseType
	}

	AssignJobOptions struct {
		RunID   string
		Phase   internal.PhaseType
		AgentID string
	}

	GetJobOptions struct {
		RunID string
		Phase internal.PhaseType
	}

	service struct {
		logr.Logger
	}
)

func (s *service) CreateJob(ctx context.Context, opts CreateJobOptions) error {
	return nil
}
