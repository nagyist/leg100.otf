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
)

type ErrJobAlreadyStarted error

type JobStatus string

// Job is the specification and status of a scheduled terraform task
type Job struct {
	ID string

	// Operation is the particular terraform task the job is carrying out: a
	// plan or an apply
	Operation Operation

	// Status is the current status of the job
	Status JobStatus

	// AgentID is the ID of the agent handling the job
	AgentID string

	// Logs are the stdout/stderr log output
	Logs []byte

	// Relations
	Run                  *Run
	Workspace            *Workspace
	ConfigurationVersion *ConfigurationVersion
}

type JobService interface {
	Create(*Run) (*Job, error)
	// Start is called by an agent when it starts the job. ErrJobAlreadyStarted
	// should be returned if another agent has already started it.
	Start(id string, opts JobStartOptions) error

	// Finish is called to signal completion of the job
	Finish(id string, opts JobFinishOptions) error

	// Cancel cancels a job using its run ID
	Cancel(runID string) error

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
		ID:                   GenerateID("job"),
		Status:               JobPending,
		Run:                  run,
		Workspace:            run.Workspace,
		ConfigurationVersion: run.ConfigurationVersion,
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
func (j *Job) Start(agentID string) error {
	if j.Status == JobStarted {
		return NewErrJobAlreadyStarted(j.AgentID)
	}

	j.Status = JobStarted
	j.AgentID = agentID

	j.Operation.Start(j)

	return nil
}

// Finish updates the state of the job to indicate an agent has finished it.
func (j *Job) Finish(opts JobFinishOptions) error {
	j.Status = opts.Status

	j.Operation.Finish(j)

	return nil
}
