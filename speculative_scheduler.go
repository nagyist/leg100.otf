package otf

import (
	"context"

	"github.com/go-logr/logr"
)

type SpeculativeScheduler struct {
	RunService
	incoming <-chan *Run
	logr.Logger
}

func NewSpeculativeScheduler(ctx context.Context, logger logr.Logger, rs RunService) (*SpeculativeScheduler, error) {
	lw, err := rs.ListWatch(ctx, RunListOptions{Statuses: IncompleteRun})
	if err != nil {
		return nil, err
	}
	return &SpeculativeScheduler{
		RunService: rs,
		Logger:     logger,
		incoming:   lw,
	}, nil
}

func (s *SpeculativeScheduler) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case run := <-s.incoming:
			if err := s.handleRun(ctx, run); err != nil {
				s.Error(err, "scheduling speculative run")
			}
		}
	}
	return
}

func (s *SpeculativeScheduler) handleRun(ctx context.Context, run *Run) error {
	if !run.Speculative() {
		return nil
	}
	switch run.Status() {
	case RunPending:
		// immediately enqueue plan for speculative runs
		_, err := s.RunService.Start(ctx, run.ID())
		if err != nil {
			return err
		}
		// publish event
	case RunPlanned:
		// finish run after plan phase
		_, err := s.RunService.UpdateStatus(ctx, RunGetOptions{ID: String(run.ID())}, func(run *Run) error {
			run.updateStatus(RunPlannedAndFinished)
			return nil
		})
		if err != nil {
			return err
		}
		// publish event
	}
	return nil
}
