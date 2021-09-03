package sqlite

import (
	"github.com/leg100/ots"
	"gorm.io/gorm"
)

// Job models a row in a jobs table.
type Job struct {
	gorm.Model

	ExternalID string `gorm:"uniqueIndex"`

	Status ots.JobStatus

	AgentID string

	Logs []byte

	// Job belongs to a run
	RunID uint
}

// JobList is a list of job models
type JobList []Job

func (model *Job) ToDomain(run *Run) *ots.Job {
	domain := ots.Job{
		ID:      model.ExternalID,
		AgentID: model.AgentID,
		Logs:    model.Logs,
		RunID:   run.ExternalID,
	}

	if run.ConfigurationVersion != nil {
		domain.ConfigurationVersionID = run.ConfigurationVersion.ExternalID
	}

	if run.Workspace != nil {
		domain.WorkspaceID = run.Workspace.ExternalID
	}

	return &domain
}

// NewJobFromDomain constructs a model obj from a domain obj
func NewJobFromDomain(domain *ots.Job, run *Run) *Job {
	model := &Job{
		RunID: run.ID,
	}
	model.FromDomain(domain, run)

	return model
}

// FromDomain updates job model fields with a job domain object's fields
func (model *Job) FromDomain(domain *ots.Job, run *Run) {
	model.ExternalID = domain.ID
	model.Status = domain.Status
	model.Logs = domain.Logs
	model.AgentID = domain.AgentID
}

func (l JobList) ToDomain(run *Run) (dl []*ots.Job) {
	for _, i := range l {
		dl = append(dl, i.ToDomain(run))
	}
	return
}
