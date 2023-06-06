package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/pubsub"
	"github.com/leg100/otf/internal/run"
)

type (
	// Creator creates jobs
	Creator struct {
		Service
		run.RunService
		pubsub.Subscriber
	}
)

func (s *Creator) Start(ctx context.Context) error {
	// Unsubscribe whenever exiting this routine.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// subscribe to run events
	sub, err := s.Subscribe(ctx, "job-creator-")
	if err != nil {
		return err
	}

	// retrieve existing runs in a queued state
	existing, err := run.ListAll(ctx, s.RunService, run.RunListOptions{
		Statuses: []internal.RunStatus{
			internal.RunPlanQueued,
			internal.RunApplyQueued,
		},
	})
	if err != nil {
		return err
	}

	// spool existing runs in reverse order; runs.ListAll returns runs newest first,
	// whereas we want oldest first.
	for i := len(existing) - 1; i >= 0; i-- {
		if err := s.handle(ctx, existing[i]); err != nil {
			return fmt.Errorf("spooling existing run: %w", err)
		}
	}
	// then spool run events
	for event := range sub {
		run, ok := event.Payload.(*run.Run)
		if !ok {
			continue
		}
		if err := s.handle(ctx, run); err != nil {
			return fmt.Errorf("spooling run event: %w", err)
		}
	}

	return nil
}

func (s *Creator) handle(ctx context.Context, run *run.Run) error {
	// map run status to phase
	var phase internal.PhaseType
	switch run.Status {
	case internal.RunPlanQueued:
		phase = internal.PlanPhase
	case internal.RunApplyQueued:
		phase = internal.ApplyPhase
	default:
		return fmt.Errorf("unexpected run status: %s", run.Status)
	}
	// check whether job has already been created for phase
	_, err := s.GetJob(ctx, GetJobOptions{RunID: run.ID, Phase: phase})
	if err == nil {
		// job already created
		return nil
	} else if !errors.Is(err, internal.ErrResourceNotFound) {
		// unexpected error
		return err
	}
	return s.CreateJob(ctx, CreateJobOptions{
		RunID: run.ID,
		Phase: internal.PlanPhase,
	})
}
