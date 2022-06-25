package app

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/leg100/otf"
)

var _ otf.RunService = (*RunService)(nil)

type RunService struct {
	db    otf.DB
	es    otf.EventService
	cache otf.Cache
	logr.Logger
	*otf.RunFactory
}

func NewRunService(db otf.DB, logger logr.Logger, wss otf.WorkspaceService, cvs otf.ConfigurationVersionService, es otf.EventService, cache otf.Cache) *RunService {
	return &RunService{
		db:     db,
		es:     es,
		cache:  cache,
		Logger: logger,
		RunFactory: &otf.RunFactory{
			WorkspaceService:            wss,
			ConfigurationVersionService: cvs,
		},
	}
}

// Create constructs and persists a new run object to the db, before scheduling
// the run.
func (s RunService) Create(ctx context.Context, spec otf.WorkspaceSpec, opts otf.RunCreateOptions) (*otf.Run, error) {
	run, err := s.New(ctx, spec, opts)
	if err != nil {
		s.Error(err, "constructing new run")
		return nil, err
	}

	if err = s.db.CreateRun(ctx, run); err != nil {
		s.Error(err, "creating run", "id", run.ID())
		return nil, err
	}
	s.V(1).Info("created run", "id", run.ID())

	s.es.Publish(otf.Event{Type: otf.EventRunCreated, Payload: run})

	return run, nil
}

// Get retrieves a run obj with the given ID from the db.
func (s RunService) Get(ctx context.Context, id string) (*otf.Run, error) {
	run, err := s.db.GetRun(ctx, otf.RunGetOptions{ID: &id})
	if err != nil {
		s.Error(err, "retrieving run", "id", id)
		return nil, err
	}

	s.V(2).Info("retrieved run", "id", run.ID())

	return run, nil
}

// List retrieves multiple run objs. Use opts to filter and paginate the list.
func (s RunService) List(ctx context.Context, opts otf.RunListOptions) (*otf.RunList, error) {
	rl, err := s.db.ListRuns(ctx, opts)
	if err != nil {
		s.Error(err, "listing runs")
		return nil, err
	}

	s.V(2).Info("listed runs", append(opts.LogFields(), "count", len(rl.Items))...)

	return rl, nil
}

// ListWatch lists runs and then watches for changes to runs. Note: The options
// filter the list but not the watch.
func (s RunService) ListWatch(ctx context.Context, opts otf.RunListOptions) (<-chan *otf.Run, error) {
	// retrieve incomplete runs from db
	existing, err := s.db.ListRuns(ctx, opts)
	if err != nil {
		return nil, err
	}
	spool := make(chan *otf.Run, len(existing.Items))
	for _, r := range existing.Items {
		spool <- r
	}
	sub, err := s.es.Subscribe("run-listwatch")
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				// context cancelled; shutdown spooler
				close(spool)
				return
			case event, ok := <-sub.C():
				if !ok {
					// sender closed channel; shutdown spooler
					close(spool)
					return
				}
				run, ok := event.Payload.(*otf.Run)
				if !ok {
					// skip non-run events
					continue
				}
				spool <- run
			}
		}
	}()
	return spool, nil
}

func (s RunService) Apply(ctx context.Context, id string, opts otf.RunApplyOptions) error {
	run, err := s.db.UpdateStatus(ctx, otf.RunGetOptions{ID: &id}, func(run *otf.Run) error {
		return run.ApplyRun()
	})
	if err != nil {
		s.Error(err, "applying run", "id", id)
		return err
	}

	s.V(0).Info("applied run", "id", id)

	s.es.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})

	return err
}

func (s RunService) UpdateStatus(ctx context.Context, opts otf.RunGetOptions, fn func(*otf.Run) error) (*otf.Run, error) {
	return s.db.UpdateStatus(ctx, opts, fn)
}

func (s RunService) Discard(ctx context.Context, id string, opts otf.RunDiscardOptions) error {
	run, err := s.db.UpdateStatus(ctx, otf.RunGetOptions{ID: &id}, func(run *otf.Run) error {
		return run.Discard()
	})
	if err != nil {
		s.Error(err, "discarding run", "id", id)
		return err
	}

	s.V(0).Info("discarded run", "id", id)

	s.es.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})

	return err
}

// Cancel enqueues a cancel request to cancel a currently queued or active plan
// or apply.
func (s RunService) Cancel(ctx context.Context, id string, opts otf.RunCancelOptions) error {
	_, err := s.db.UpdateStatus(ctx, otf.RunGetOptions{ID: &id}, func(run *otf.Run) error {
		return run.Cancel()
	})
	return err
}

func (s RunService) ForceCancel(ctx context.Context, id string, opts otf.RunForceCancelOptions) error {
	_, err := s.db.UpdateStatus(ctx, otf.RunGetOptions{ID: &id}, func(run *otf.Run) error {
		return run.ForceCancel()

		// TODO: send KILL signal to running terraform process

		// TODO: unlock workspace - could do this by publishing event and having
		// workspace scheduler subscribe to event
	})

	return err
}

func (s RunService) Start(ctx context.Context, runID string) (*otf.Run, error) {
	run, err := s.UpdateStatus(ctx, otf.RunGetOptions{ID: &runID}, func(run *otf.Run) error {
		return run.EnqueuePlan()
	})
	if err != nil {
		s.Error(err, "started run", "id", runID)
		return nil, err
	}
	s.V(0).Info("started run", "id", runID)

	s.es.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})

	return run, nil
}

// GetPlanFile returns the plan file for the run.
func (s RunService) GetPlanFile(ctx context.Context, planID string, format otf.PlanFormat) ([]byte, error) {
	if plan, err := s.cache.Get(format.CacheKey(planID)); err == nil {
		return plan, nil
	}
	// Cache is empty; retrieve from DB
	file, err := s.db.GetPlanFile(ctx, planID, format)
	if err != nil {
		s.Error(err, "retrieving plan file", "id", planID, "format", format)
		return nil, err
	}
	// Cache plan before returning
	if err := s.cache.Set(format.CacheKey(planID), file); err != nil {
		return nil, fmt.Errorf("caching plan: %w", err)
	}
	return file, nil
}

// UploadPlanFile persists a run's plan file. The plan file is expected to have
// been produced using `terraform plan`. If the plan file is JSON serialized
// then its parsed for a summary of planned changes and the Plan object is
// updated accordingly.
func (s RunService) UploadPlanFile(ctx context.Context, planID string, plan []byte, format otf.PlanFormat) error {
	if err := s.db.SetPlanFile(ctx, planID, plan, format); err != nil {
		s.Error(err, "uploading plan file", "plan_id", planID, "format", format)
		return err
	}
	s.V(0).Info("uploaded plan file", "plan_id", planID, "format", format)

	if err := s.cache.Set(format.CacheKey(planID), plan); err != nil {
		return fmt.Errorf("caching plan: %w", err)
	}
	return nil
}

// Delete deletes a terraform run.
func (s RunService) Delete(ctx context.Context, id string) error {
	// get run first so that we can include it in an event below
	run, err := s.db.GetRun(ctx, otf.RunGetOptions{ID: &id})
	if err != nil {
		return err
	}

	if err := s.db.DeleteRun(ctx, id); err != nil {
		s.Error(err, "deleting run", "id", id)
		return err
	}

	s.V(0).Info("deleted run", "id", id)

	s.es.Publish(otf.Event{Type: otf.EventRunDeleted, Payload: run})

	return nil
}
