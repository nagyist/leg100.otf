package jobs

import (
	"context"
	"errors"

	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/logr"
	"github.com/leg100/otf/internal/pubsub"
)

type (
	Status string

	JobService = Service

	Service interface {
		CreateJob(context.Context, CreateJobOptions) error
		// AssignJob assigns a job to an agent
		AssignJob(context.Context, AssignJobOptions) error
		// GetAssignedJob retrieves a job assigned to the given agent. The
		// returned channel only receives at most job before being closed.
		GetAssignedJob(ctx context.Context, agentID string) (<-chan *Job, error)
		// GetJobByPhase retrieves a job by run ID and phase.
		GetJobByPhase(ctx context.Context, runID string, phase internal.PhaseType) (*Job, error)
		UpdateJobStatus(context.Context, Status) error
		ListJobs(context.Context, ListJobsOptions) ([]*Job, error)
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

	ListJobsOptions struct {
		Status Status // filter jobs by status
	}

	service struct {
		logr.Logger
		pubsub.Subscriber

		db *pgdb
	}
)

func (s *service) CreateJob(ctx context.Context, opts CreateJobOptions) error {
	return nil
}

func (s *service) GetAssignedJob(ctx context.Context, agentID string) (<-chan *Job, error) {
	// unsubscribe whenever exiting this routine.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// subscribe to job events
	sub, err := s.Subscribe(ctx, "get-job-")
	if err != nil {
		return nil, err
	}
	// retrieve existing assigned job, if any
	existing, err := s.db.getAssignedJobByAgentID(ctx, agentID)
	if err != nil && !errors.Is(err, internal.ErrResourceNotFound) {
		// unexpected error
		return nil, err
	}
	assigned := make(chan *Job)
	go func() {
		// send existing assigned job and return early
		if existing != nil {
			assigned <- existing
			close(assigned)
			return
		}
		// ...otherwise block and wait for a job to be assigned
		for event := range sub {
			job, ok := event.Payload.(*Job)
			if !ok {
				continue // skip anything other than a job
			}
			if job.Status != Assigned {
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
			assigned <- job
		}
		close(assigned)
	}()
	return assigned, nil
}
