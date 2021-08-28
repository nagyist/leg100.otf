package app

import (
	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
)

var _ ots.RunService = (*RunService)(nil)

type RunService struct {
	db ots.RunStore
	bs ots.BlobStore
	es ots.EventService
	js ots.JobService

	*ots.RunFactory
}

func NewRunService(db ots.RunStore, wss ots.WorkspaceService, cvs ots.ConfigurationVersionService, bs ots.BlobStore, es ots.EventService, js ots.JobService) *RunService {
	return &RunService{
		bs: bs,
		db: db,
		es: es,
		js: js,
		RunFactory: &ots.RunFactory{
			WorkspaceService:            wss,
			ConfigurationVersionService: cvs,
		},
	}
}

// Create constructs and persists a new run object to the db, before scheduling
// the run.
func (s RunService) Create(opts *tfe.RunCreateOptions) (*ots.Run, error) {
	run, err := s.NewRun(opts)
	if err != nil {
		return nil, err
	}

	run, err = s.db.Create(run)
	if err != nil {
		return nil, err
	}

	s.es.Publish(ots.Event{Type: ots.RunCreated, Payload: run})

	return run, nil
}

// Get retrieves a run obj with the given ID from the db.
func (s RunService) Get(id string) (*ots.Run, error) {
	return s.db.Get(ots.RunGetOptions{ID: &id})
}

// List retrieves multiple run objs. Use opts to filter and paginate the list.
func (s RunService) List(opts ots.RunListOptions) (*ots.RunList, error) {
	return s.db.List(opts)
}

func (s RunService) Apply(id string, opts *tfe.RunApplyOptions) error {
	run, err := s.db.Update(id, func(run *ots.Run) error {
		run.UpdateApplyStatus(tfe.ApplyQueued)

		return nil
	})
	if err != nil {
		return err
	}

	s.es.Publish(ots.Event{Type: ots.ApplyQueued, Payload: run})

	return err
}

func (s RunService) Discard(id string, opts *tfe.RunDiscardOptions) error {
	run, err := s.db.Update(id, func(run *ots.Run) error {
		return run.Discard()
	})
	if err != nil {
		return err
	}

	s.es.Publish(ots.Event{Type: ots.RunCompleted, Payload: run})

	return err
}

// Cancel enqueues a cancel request to cancel a currently queued or active plan
// or apply.
func (s RunService) Cancel(id string, opts *tfe.RunCancelOptions) error {
	_, err := s.db.Update(id, func(run *ots.Run) error {
		if err := run.IssueCancel(); err != nil {
			return err
		}

		// Immediately mark pending/queued runs as cancelled
		switch run.Status {
		case tfe.RunPending, tfe.RunPlanQueued, tfe.RunApplyQueued:
			run.Status = tfe.RunCanceled
			run.StatusTimestamps.CanceledAt = ots.TimeNow()
		}

		return nil
	})
	return err
}

func (s RunService) ForceCancel(id string, opts *tfe.RunForceCancelOptions) error {
	_, err := s.db.Update(id, func(run *ots.Run) error {
		if err := run.ForceCancel(); err != nil {
			return err
		}

		// TODO: send KILL signal to running terraform process

		// TODO: unlock workspace

		return nil
	})

	return err
}

// EnqueuePlan enqueues a run's plan by creating a plan job
func (s RunService) EnqueuePlan(id string) error {
	run, err := s.db.Get(ots.RunGetOptions{ID: &id})
	if err != nil {
		return err
	}

	job, err := s.js.Create(run)
	if err != nil {
		return err
	}

	run, err = s.db.Update(id, func(run *ots.Run) error {
		if err := run.UpdateStatus(tfe.RunPlanQueued); err != nil {
			return err
		}
		run.PlanJob = job
		return nil
	})
	if err != nil {
		return err
	}

	s.es.Publish(ots.Event{Type: ots.PlanQueued, Payload: run})

	return err
}

func (s RunService) UpdateStatus(id string, status tfe.RunStatus) (*ots.Run, error) {
	return s.db.Update(id, func(run *ots.Run) error {
		return run.UpdateStatus(status)
	})
}

func (s RunService) UpdatePlanStatus(id string, status tfe.PlanStatus) (*ots.Run, error) {
	run, err := s.db.Update(id, func(run *ots.Run) error {
		run.UpdatePlanStatus(status)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return run, nil
}

func (s RunService) UpdateApplyStatus(id string, status tfe.ApplyStatus) (*ots.Run, error) {
	run, err := s.db.Update(id, func(run *ots.Run) error {
		run.UpdateApplyStatus(status)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return run, nil
}

func (s RunService) FinishPlan(id string, opts ots.PlanFinishOptions) (*ots.Run, error) {
	planFileBlobID, err := s.bs.Put(opts.Plan)
	if err != nil {
		return nil, err
	}

	planJSONBlobID, err := s.bs.Put(opts.PlanJSON)
	if err != nil {
		return nil, err
	}

	run, err := s.db.Update(id, func(run *ots.Run) error {
		run.FinishPlan(opts)
		run.Plan.PlanFileBlobID = planFileBlobID
		run.Plan.PlanJSONBlobID = planJSONBlobID

		return nil
	})
	if err != nil {
		return nil, err
	}
	return run, nil
}

func (s RunService) FinishApply(id string, opts ots.ApplyFinishOptions) (*ots.Run, error) {
	run, err := s.db.Update(id, func(run *ots.Run) error {
		run.FinishApply(opts)

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.es.Publish(ots.Event{Type: ots.RunCompleted, Payload: run})

	return run, nil
}

// GetPlanJSON returns the JSON formatted plan file for the run.
func (s RunService) GetPlanJSON(id string) ([]byte, error) {
	run, err := s.db.Get(ots.RunGetOptions{ID: &id})
	if err != nil {
		return nil, err
	}
	return s.bs.Get(run.Plan.PlanJSONBlobID)
}

// GetPlanFile returns the binary plan file for the run.
func (s RunService) GetPlanFile(id string) ([]byte, error) {
	run, err := s.db.Get(ots.RunGetOptions{ID: &id})
	if err != nil {
		return nil, err
	}
	return s.bs.Get(run.Plan.PlanFileBlobID)
}
