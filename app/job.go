package app

import (
	"fmt"

	"github.com/leg100/ots"
)

var _ ots.JobService = (*JobService)(nil)

type JobService struct {
	db ots.JobStore
	es ots.EventService
}

func NewJobService(db ots.JobStore, es ots.EventService) *JobService {
	return &JobService{
		db: db,
		es: es,
	}
}

// Create constructs and persists a new job object to the db and sends a
// notification a job has been created.
func (s JobService) Create(run *ots.Run) (job *ots.Job, err error) {
	job = ots.NewJobFromRun(run)

	job, err = s.db.Create(job)
	if err != nil {
		return nil, err
	}

	s.es.Publish(ots.Event{Type: ots.JobCreated, Payload: job})

	return job, nil
}

func (s JobService) Start(id string, opts ots.JobStartOptions) error {
	_, err := s.db.Update(id, func(job *ots.Job) error {
		return job.Start(opts.AgentID)
	})
	return err
}

func (s JobService) Finish(id string) error {
	_, err := s.db.Update(id, func(job *ots.Job) error {
		job.Status = ots.JobCompleted
		return nil
	})
	return err
}

func (s JobService) UploadLogs(id string, out []byte) error {
	_, err := s.db.Update(id, func(job *ots.Job) error {
		job.Logs = out
		return nil
	})
	return err
}

func (s JobService) GetLogs(id string, opts ots.JobLogOptions) ([]byte, error) {
	job, err := s.db.Get(id)
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

// Get retrieves a job obj with the given ID from the db.
func (s JobService) Get(id string) (*ots.Job, error) {
	return s.db.Get(id)
}

// List retrieves multiple job objs. Use opts to filter and paginate the list.
func (s JobService) List() ([]*ots.Job, error) {
	return s.db.List()
}
