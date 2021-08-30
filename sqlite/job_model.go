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

	// Job belongs to a run
	RunID uint
	Run   *Run

	// Job belongs to a workspace
	WorkspaceID uint
	Workspace   *Workspace

	// Job belongs to a configuration version
	ConfigurationVersionID uint
	ConfigurationVersion   *ConfigurationVersion
}

// JobList is a list of job models
type JobList []Job

// Update updates the model with the supplied fn. The fn operates on the domain
// obj, so Update handles converting to and from a domain obj.
func (model *Job) Update(fn func(*ots.Job) error) error {
	// model -> domain
	domain := model.ToDomain()

	// invoke user fn
	if err := fn(domain); err != nil {
		return err
	}

	// domain -> model
	model.FromDomain(domain)

	return nil
}

func (model *Job) ToDomain() *ots.Job {
	domain := ots.Job{
		ID: model.ExternalID,
	}

	if model.Run != nil {
		domain.Run = model.Run.ToDomain()
	}

	if model.ConfigurationVersion != nil {
		domain.ConfigurationVersion = model.ConfigurationVersion.ToDomain()
	}

	if model.Workspace != nil {
		domain.Workspace = model.Workspace.ToDomain()
	}

	return &domain
}

// NewJobFromDomain constructs a model obj from a domain obj
func NewJobFromDomain(domain *ots.Job) *Job {
	model := &Job{
		ConfigurationVersion: &ConfigurationVersion{},
		Workspace:            &Workspace{},
	}
	model.FromDomain(domain)

	return model
}

// FromDomain updates run model fields with a run domain object's fields
func (model *Job) FromDomain(domain *ots.Job) {
	model.ExternalID = domain.ID
	model.Status = domain.Status

	model.Run.FromDomain(domain.Run)

	model.Workspace.FromDomain(domain.Workspace)
	model.WorkspaceID = domain.Workspace.Model.ID

	model.ConfigurationVersion.FromDomain(domain.ConfigurationVersion)
	model.ConfigurationVersionID = domain.ConfigurationVersion.Model.ID
}

func (l JobList) ToDomain() (dl []*ots.Job) {
	for _, i := range l {
		dl = append(dl, i.ToDomain())
	}
	return
}
