package app

import (
	"fmt"

	"github.com/leg100/ots"
)

var _ ots.JobService = (*JobService)(nil)

type JobService struct {
	db ots.JobStore
	rs ots.RunService
}

func NewJobService(db ots.JobStore, rs ots.RunService) *JobService {
	return &JobService{
		db: db,
		rs: rs,
	}
}

func (s JobService) Start(id string, opts ots.JobStartOptions) error {
	return s.rs.StartJob(id, opts)
}

func (s JobService) Finish(id string, opts ots.JobFinishOptions) error {
	return s.rs.FinishJob(id, opts)
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

// List retrieves multiple jobs.
func (s JobService) List() ([]*ots.Job, error) {
	return s.db.List()
}
