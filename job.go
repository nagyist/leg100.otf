package ots

import (
	"fmt"

	tfe "github.com/leg100/go-tfe"
)

const (
	JobPending   JobStatus = "pending"
	JobStarted   JobStatus = "started"
	JobCompleted JobStatus = "completed"
	JobErrored   JobStatus = "errored"
	JobCanceled  JobStatus = "canceled"
)

type ErrJobAlreadyStarted error

type JobStatus string

// Job is the specification and status of a scheduled terraform task
type Job struct {
	ID string

	// Operation is the particular terraform task the job is carrying out: a
	// plan or an apply
	Operation Step

	// Status is the current status of the job
	Status JobStatus

	// AgentID is the ID of the agent handling the job
	AgentID string

	// RunID is the ID of the run the job belongs to.
	RunID string

	// WorkspaceID is the ID of the workspace the job belongs to.
	WorkspaceID string

	// ConfigurationVersionID is the ID of the configuration version pertaining
	// to the job.
	ConfigurationVersionID string

	// StateVersionID is the ID of the state version pertaining to the job.
	StateVersionID string

	// Logs are the stdout/stderr log output
	Logs []byte
}

type JobService interface {
	// Start is called by an agent when it starts the job. ErrJobAlreadyStarted
	// should be returned if another agent has already started it.
	Start(id string, opts JobStartOptions) error

	// Finish is called to signal completion of the job
	Finish(id string, opts JobFinishOptions) error

	UploadLogs(id string, out []byte) error
}

// JobStore implementations persist Job objects.
type JobStore interface {
	Create(*Job) (*Job, error)
	Get(id string) (*Job, error)
	List() ([]*Job, error)
	// TODO: add support for a special error type that tells update to skip
	// updates - useful when fn checks current fields and decides not to update
	Update(id string, fn func(*Job) error) (*Job, error)
}

type JobStartOptions struct {
	// AgentID is the ID of the agent starting the job
	AgentID string
}

type JobFinishOptions struct {
	Status JobStatus
}

type JobLogOptions struct {
	// The maximum number of bytes of logs to return to the client
	Limit int `schema:"limit"`

	// The start position in the logs from which to send to the client
	Offset int `schema:"offset"`
}

// NewJobFromRun constructs a job from a run.
func NewJobFromRun(run *Run) (*Job, error) {
	job := &Job{
		ID:                     GenerateID("job"),
		Status:                 JobPending,
		RunID:                  run.ID,
		WorkspaceID:            run.Workspace.ID,
		ConfigurationVersionID: run.ConfigurationVersion.ID,
	}

	switch run.Status {
	case tfe.RunPlanQueued:
		job.Operation = NewPlanOperation()
	case tfe.RunApplyQueued:
		job.Operation = NewApplyOperation()
	default:
		return nil, fmt.Errorf("invalid run status for new job: %s", run.Status)
	}

	return job, nil
}

func NewErrJobAlreadyStarted(agentID string) ErrJobAlreadyStarted {
	return ErrJobAlreadyStarted(fmt.Errorf("job already started by agent %s", agentID))
}

// Start updates the state of the job to indicate an agent has started it.
func (j *Job) Start(opts JobStartOptions) error {
	if j.Status == JobStarted {
		return NewErrJobAlreadyStarted(j.AgentID)
	}

	j.Status = JobStarted
	j.AgentID = opts.AgentID

	return nil
}

// Finish updates the state of the job to indicate an agent has finished it.
func (j *Job) Finish(opts JobFinishOptions) error {
	j.Status = opts.Status

	return nil
}
