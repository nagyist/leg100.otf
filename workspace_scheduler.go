package otf

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
)

const (
	// SchedulerSubscriptionID is the ID the scheduler uses to identify itself
	// when subscribing to the events service
	SchedulerSubscriptionID = "scheduler"
)

// Scheduler schedules workspace runs: each workspace can only have one active
// run at a time, blocking all other runs. The runs are processed in order.
// Note: this does not apply to speculative runs, which can run in parallel.
type Scheduler struct {
	// WorkspaceService for enumerating and locking/unlocking workspaces.
	WorkspaceService

	// RunService for enumerating and activating runs.
	RunService

	// EventService for subscribing to workspace and run events
	EventService

	// Queues is a mapping of a workspace's ID to its queue of runs.
	Queues map[string]WorkspaceQueue

	logr.Logger
}

// NewScheduler constructs the scheduler, populating its workspace queues with
// existing runs.
func NewScheduler(ws WorkspaceService, rs RunService, es EventService, logger logr.Logger) (*Scheduler, error) {
	queues := make(map[string]WorkspaceQueue)

	workspaces, err := ws.List(context.Background(), WorkspaceListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing workspaces: %s", err)
	}
	for _, ws := range workspaces.Items {
		// TODO: exclude completed runs
		runs, err := rs.List(context.Background(), RunListOptions{WorkspaceID: &ws.ID})
		if err != nil {
			return nil, err
		}
		queues[ws.ID], err = NewWorkspaceQueue(runs.Items...)
		if err != nil {
			return nil, err
		}
	}

	s := &Scheduler{
		Queues:       queues,
		RunService:   rs,
		EventService: es,
		Logger:       logger.WithValues("component", "scheduler"),
	}

	return s, nil
}

// Start starts the scheduler event loop
func (s *Scheduler) Start(ctx context.Context) error {
	sub, err := s.Subscribe(SchedulerSubscriptionID)
	if err != nil {
		return err
	}

	defer sub.Close()

	for {
		select {
		case event, ok := <-sub.C():
			// If sub closed then exit.
			if !ok {
				return nil
			}
			s.handleEvent(event)
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Scheduler) handleEvent(ev Event) {
	s.Info("event received", "event", ev.Type)

	switch obj := ev.Payload.(type) {
	case *Workspace:
		switch ev.Type {
		case EventWorkspaceCreated:
			s.Queues[obj.ID] = WorkspaceQueue{}
		case EventWorkspaceDeleted:
			delete(s.Queues, obj.ID)
		}
	case *Run:
		queue, ok := s.Queues[obj.Workspace.ID]
		if !ok {
			s.Error(fmt.Errorf("received run event for non-existent workspace"), "", "workspace", obj.Workspace.ID)
			return
		}

		s.handleRun(queue, obj)
	}
}

func (s *Scheduler) handleRun(queue WorkspaceQueue, run *Run) {
	// Scheduler does not handle speculative runs
	if run.IsSpeculative() {
		return
	}

	// No action to be taken on runs that are still active
	if run.IsActive() {
		return
	}

	// Add pending runs to back of pending queue.
	if obj.Status == RunPending {
		if !queue.inPendingQueue(run)
			queue.Pending = append(queue.Pending, obj)
		}
	}

	// Unset active run if it has completed.
	if run.IsDone() && run.ID == queue.Active.ID {
		queue.Active = nil
	}

	// A run is still active, nothing to be done.
	if queue.Active != nil {
		return
	}

	// Workspace queue is empty, ensure workspace is unlocked
	if len(queue.Pending) == 0 {
}


func (s *Scheduler) makeActive(queue WorkspaceQueue, run *Run) {
	for {
		_, err := s.Lock(ctx, WorkspaceSpec{ID: &obj.Workspace.ID}, WorkspaceLockOptions{})
		if err != nil {
			s.Error(err, "unable to activate run", "run", run.ID)
			continue
		}

		return
	}
}

// Update takes an updated run
func (s *Scheduler) Update(run *Run) error {
	// Enqueue speculative runs onto (global) queue but don't make them active
	// because they do not block pending runs
	if run.IsSpeculative() {
		return q.EnqueuePlan(context.Background(), run.ID)
	}

	// No run is current active, so make this run active
	if q.Active == nil {
		if err := q.EnqueuePlan(context.Background(), run.ID); err != nil {
			return err
		}

		q.Active = run
		return nil
	}

	// Other add run to pending queue
	q.Pending = append(q.Pending, run)

	return nil
}

// WorkspaceQueue maintains the state of the workspace queue in memory.
type WorkspaceQueue struct {
	// Active is the currently active run.
	Active *Run

	// Pending is the list of pending runs waiting for the active run to
	// complete.
	Pending []*Run
}

// NewWorkspaceQueue constructs a workspace queue, populated with existing runs.
func NewWorkspaceQueue(runs ...*Run) (q WorkspaceQueue, err error) {
	// TODO: ensure runs are sorted by created_at, ascending.

	for _, run := range runs {
		switch {
		case run.IsSpeculative():
			continue
		case run.Status == RunPending:
			q.Pending = append(q.Pending, run)
		case run.IsActive():
			if q.Active != nil {
				return WorkspaceQueue{}, fmt.Errorf("more than one active run found: %s and %s", q.Active, run)
			}
			q.Active = run
		}
	}
	return
}

func (q WorkspaceQueue) inPendingQueue(run *Run) bool {
	for _, pending := range q.Pending {
		if run.ID == pending.ID {
			return true
		}
	}
	return false
}

// Add adds a run to the workspace queue.
func (q *WorkspaceQueue) Add(run *Run) error {
	// Enqueue speculative runs onto (global) queue but don't make them active
	// because they do not block pending runs
	if run.IsSpeculative() {
		return q.EnqueuePlan(context.Background(), run.ID)
	}

	// No run is current active, so make this run active
	if q.Active == nil {
		if err := q.EnqueuePlan(context.Background(), run.ID); err != nil {
			return err
		}

		q.Active = run
		return nil
	}

	// Other add run to pending queue
	q.Pending = append(q.Pending, run)

	return nil
}

// Remove removes a run from the queue.
func (q *WorkspaceQueue) Remove(run *Run) error {
	// Speculative runs are never added to the queue in the first place so they
	// do not need to be removed
	if run.IsSpeculative() {
		return nil
	}

	// Remove active run and make the first pending run the active run
	if q.Active.ID == run.ID {
		q.Active = nil
		if len(q.Pending) > 0 {
			if err := q.EnqueuePlan(context.Background(), q.Pending[0].ID); err != nil {
				return err
			}

			q.Active = q.Pending[0]
			q.Pending = q.Pending[1:]
		}
		return nil
	}

	// Remove run from pending queue
	for idx, p := range q.Pending {
		if p.ID == run.ID {
			q.Pending = append(q.Pending[:idx], q.Pending[idx+1:]...)
			return nil
		}
	}

	return nil
}
