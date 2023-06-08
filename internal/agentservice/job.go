package agentservice

import "github.com/leg100/otf/internal"

const (
	UnassignedJob JobStatus = "unassigned"
	AssignedJob   JobStatus = "assigned"
	RunningJob    JobStatus = "running"
	CompletedJob  JobStatus = "completed"
	ErroredJob    JobStatus = "errored"
)

type (
	Job struct {
		ID      string
		RunID   string
		Phase   internal.PhaseType
		AgentID *string
		Status  JobStatus
	}

	JobStatus string
)
