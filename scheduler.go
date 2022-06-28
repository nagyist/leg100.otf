package otf

import (
	"context"

	"github.com/go-logr/logr"
)

// Scheduler schedules runs and updates workspace queues.
type Scheduler struct {
	// RunService for scheduling runs
	RunService
	// WorkspaceService for interacting with workspace queues and locking
	// workspaces
	WorkspaceService
	// Subscription for subscribing to workspace unlock events
	Subscription
	// Logger for logging messages for humans
	logr.Logger
	// incoming is a stream of changes to runs
	incoming <-chan *Run
}

// NewScheduler constructs and initialises the scheduler.
func NewScheduler(ctx context.Context, logger logr.Logger, app Application) (*Scheduler, error) {
	lw, err := app.RunService().ListWatch(ctx, RunListOptions{Statuses: IncompleteRun})
	if err != nil {
		return nil, err
	}
	sub, err := app.EventService().Subscribe("scheduler")
	if err != nil {
		return nil, err
	}
	return &Scheduler{
		RunService:       app.RunService(),
		WorkspaceService: app.WorkspaceService(),
		Logger:           logger,
		incoming:         lw,
		Subscription:     sub,
	}, nil
}

// Start starts the scheduler daemon. Should be invoked in a go routine.
func (s *Scheduler) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case run := <-s.incoming:
			if err := s.handleRun(ctx, run); err != nil {
				s.Error(err, "scheduling run", "run", run.ID())
			}
		case event := <-s.C():
			ws, ok := event.Payload.(*Workspace)
			if !ok {
				// skip non-workspace events
				continue
			}
			// the scheduler watches for workspace unlock events because it may
			// be waiting for it to be unlocked (say by a user) before it can
			// lock the workspace.
			if event.Type == EventWorkspaceUnlocked {
				if err := s.checkQueue(ctx, ws.ID()); err != nil {
					s.Error(err, "checking workspace queue", "workspace", ws.ID())
				}
			}
		}
	}
}

func (s *Scheduler) handleRun(ctx context.Context, run *Run) error {
	if run.Speculative() {
		if run.Status() == RunPending {
			// immediately enqueue plan for pending speculative runs
			_, err := s.RunService.EnqueuePlan(ctx, run.ID())
			if err != nil {
				return err
			}
		}
		// speculative runs need no further scheduling
		return nil
	}
	//
	// Below onwards we deal with non-speculative runs, i.e. runs that are
	// queued and lock a workspace.
	//

	// If a run is finished then unlock the workspace.
	if run.Done() {
		_, err := s.WorkspaceService.Unlock(ctx, WorkspaceSpec{ID: &run.workspaceID}, WorkspaceUnlockOptions{
			Requestor: run,
		})
		if err != nil {
			return err
		}
	}
	// Update workspace queue with updated run
	if err := s.WorkspaceService.UpdateQueue(run); err != nil {
		return err
	}
	// Now check queue to see if run needs its plan enqueued.
	if err := s.checkQueue(ctx, run.workspaceID); err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) checkQueue(ctx context.Context, workspaceID string) error {
	queue, err := s.WorkspaceService.GetQueue(workspaceID)
	if err != nil {
		return err
	}
	if len(queue) == 0 {
		return nil
	}
	if queue[0].Status() != RunPending {
		return nil
	}
	_, err = s.RunService.EnqueuePlan(ctx, queue[0].ID())
	if err != nil {
		return err
	}
	return nil
}
