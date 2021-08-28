package ots

import (
	"context"
	"fmt"
	"io"

	tfe "github.com/leg100/go-tfe"
)

const (
	JobPlanOp  JobOperation = "plan"
	JobApplyOp JobOperation = "apply"

	JobPending   JobStatus = "pending"
	JobStarted   JobStatus = "started"
	JobCompleted JobStatus = "completed"
	JobErrored   JobStatus = "errored"
)

type ErrJobAlreadyStarted error

type JobOperation string

type JobStatus string

// Job is a scheduled terraform task (as of writing, a plan or apply)
type Job struct {
	ID string
	// Do does the piece of work Do(ctx context.Context) error
	Status JobStatus

	// Step is the task (or task composed of tasks) the job carries out.
	Step

	// AgentID is the ID of the agent starting the job
	AgentID string

	// Logs are the stdout/stderr log output
	Logs []byte

	// Relations
	Plan                 *Plan
	Apply                *Apply
	Workspace            *Workspace
	ConfigurationVersion *ConfigurationVersion
}

type JobService interface {
	Create(*Run) (*Job, error)
	// Start is called by an agent when it starts the job. ErrJobAlreadyStarted
	// should be returned if another agent has already started it.
	Start(id string, opts JobStartOptions) error

	UploadLogs(id string, out []byte) error

	// Finish is called to signal completion of the job
	Finish(id string) error
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

type JobLogOptions struct {
	// The maximum number of bytes of logs to return to the client
	Limit int `schema:"limit"`

	// The start position in the logs from which to send to the client
	Offset int `schema:"offset"`
}

type Doer interface {
	Do(ctx context.Context, path string, out io.Writer) error
}

// NewJobFromRun constructs a job from a run.
func NewJobFromRun(run *Run) (*Job, error) {
	job := &Job{
		ID:     GenerateID("job"),
		Status: JobPending,
	}

	switch run.Status {
	case tfe.RunPlanQueued:
		job.Step = NewMultiStep([]Step{
			NewDownloadConfigStep(run),
			NewDeleteBackendStep,
			NewDownloadStateStep(run),
			NewInitStep,
			NewPlanStep,
			NewJSONPlanStep,
		})
	}

	return job, nil
}

// Start updates the state of the job to indicate an agent has started it.
func (j *Job) Start(agentID string) error {
	if j.Status == JobStarted {
		return ErrJobAlreadyStarted(fmt.Errorf("job already started by agent %s", j.AgentID))

	}

	j.Status = JobStarted
	j.AgentID = agentID

	return nil
}
