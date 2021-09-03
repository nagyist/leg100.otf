package app

import (
	"fmt"

	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
)

var _ ots.RunService = (*RunService)(nil)

type RunService struct {
	db ots.RunStore
	bs ots.BlobStore
	es ots.EventService

	*ots.RunFactory
}

func NewRunService(db ots.RunStore, wss ots.WorkspaceService, cvs ots.ConfigurationVersionService, bs ots.BlobStore, es ots.EventService) *RunService {
	return &RunService{
		bs: bs,
		db: db,
		es: es,
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

// ListJobs retrieves a list of jobs. Use opts to filter the list.
func (s RunService) ListJobs(opts ots.JobLogOptions) ([]*ots.Job, error) {
	var jobs []*ots.Job

	runs, err := s.db.List(ots.RunListOptions{})
	if err != nil {
		return nil, err
	}

	for _, run := range runs.Items {
		if run.ApplyJob != nil {
			jobs = append(jobs, run.ApplyJob)
		}
		if run.PlanJob != nil {
			jobs = append(jobs, run.PlanJob)
		}
	}

	return jobs, nil
}

func (s RunService) Apply(id string, opts *tfe.RunApplyOptions) error {
	var job *ots.Job

	_, err := s.db.Update(id, func(run *ots.Run) (err error) {
		job, err = run.EnqueueApply()
		return err
	})
	if err != nil {
		return err
	}

	s.es.Publish(ots.Event{Type: ots.JobCreatedEvent, Payload: job})

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

// Cancel cancels a run. Triggers a job canceled event if the run has an active
// job.
func (s RunService) Cancel(id string, opts *tfe.RunCancelOptions) error {
	var job *ots.Job

	_, err := s.db.Update(id, func(run *ots.Run) (err error) {
		job, err = run.Cancel()
		return err
	})

	s.es.Publish(ots.Event{Type: ots.JobCanceledEvent, Payload: job})

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
	var job *ots.Job

	_, err := s.db.Update(id, func(run *ots.Run) (err error) {
		job, err = run.EnqueuePlan()
		return err
	})
	if err != nil {
		return err
	}

	s.es.Publish(ots.Event{Type: ots.JobCreatedEvent, Payload: job})

	return err
}

// StartJob starts one of a run's jobs.
func (s RunService) StartJob(jobID string, opts ots.JobStartOptions) error {
	run, err := s.db.Get(ots.RunGetOptions{JobID: &jobID})
	if err != nil {
		return err
	}

	_, err = s.db.Update(run.ID, func(run *ots.Run) error {
		return run.StartJob(jobID, opts)
	})
	if err != nil {
		return err
	}

	return err
}

// FinishJob finishes one of a run's jobs.
func (s RunService) FinishJob(jobID string, opts ots.JobFinishOptions) error {
	run, err := s.db.Get(ots.RunGetOptions{JobID: &jobID})
	if err != nil {
		return err
	}

	// If finished job is a plan then a new apply job might be created (i.e. if
	// AutoApply is enabled)
	var applyJob *ots.Job

	_, err = s.db.Update(run.ID, func(run *ots.Run) (err error) {
		applyJob, err = run.FinishJob(jobID, opts)
		return err
	})
	if err != nil {
		return err
	}

	if applyJob != nil {
		s.es.Publish(ots.Event{Type: ots.JobCreatedEvent, Payload: applyJob})
	}

	return err
}

// UploadJobLogs persists a run job's logs.
//
// TODO: avoid DB queries - job constructor should instead create a GUID/BlobID
// and then client uses this to upload logs and server can then use that to
// persist directly to blob service
func (s RunService) UploadJobLogs(jobID string, out []byte) error {
	run, err := s.db.Get(ots.RunGetOptions{JobID: &jobID})
	if err != nil {
		return nil
	}

	_, err = s.db.Update(run.ID, func(run *ots.Run) error {
		job, err := run.GetJob(jobID)
		if err != nil {
			return nil
		}
		job.Logs = out
		return nil
	})
	return err
}

func (s RunService) GetJobLogs(jobID string, opts ots.JobLogOptions) ([]byte, error) {
	run, err := s.db.Get(ots.RunGetOptions{JobID: &jobID})
	if err != nil {
		return nil, err
	}
	job, err := run.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	logs := job.Logs

	// Add start marker
	logs = append([]byte{byte(2)}, logs...)

	// Add end marker
	logs = append(logs, byte(3))

	if opts.Offset > len(logs) {
		return nil, fmt.Errorf("offset cannot be bigger than total logs")
	}

	if opts.Limit > ots.MaxPlanLogsLimit {
		opts.Limit = ots.MaxPlanLogsLimit
	}

	// Ensure specified chunk does not exceed slice length
	if (opts.Offset + opts.Limit) > len(logs) {
		opts.Limit = len(logs) - opts.Offset
	}

	resp := logs[opts.Offset:(opts.Offset + opts.Limit)]

	return resp, nil
}

func (s RunService) UpdateStatus(id string, status tfe.RunStatus) (*ots.Run, error) {
	return s.db.Update(id, func(run *ots.Run) error {
		return run.UpdateStatus(status)
	})
}

// UploadPlan persists a run's plan file. The plan file is expected to have been
// produced using `terraform plan`. If the plan file is JSON serialized then set
// json to true.
func (s RunService) UploadPlan(id string, plan []byte, json bool) error {
	blobID, err := s.bs.Put(plan)
	if err != nil {
		return err
	}

	_, err = s.db.Update(id, func(run *ots.Run) error {
		if json {
			run.Plan.PlanJSONBlobID = blobID
		} else {
			run.Plan.PlanFileBlobID = blobID
		}

		return nil
	})
	return err
}

// UpdatePlanSummary updates the resource summary for a run's plan.
func (s RunService) UpdatePlanSummary(id string, summary ots.ResourceSummary) error {
	_, err := s.db.Update(id, func(run *ots.Run) error {
		run.Plan.ResourceSummary = summary

		return nil
	})
	return err
}

// UpdateApplySummary updates the resource summary for a run's plan.
func (s RunService) UpdateApplySummary(id string, summary ots.ResourceSummary) error {
	_, err := s.db.Update(id, func(run *ots.Run) error {
		run.Apply.ResourceSummary = summary

		return nil
	})
	return err
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
