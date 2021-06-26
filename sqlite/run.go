package sqlite

import (
	"time"

	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ots.RunService = (*RunService)(nil)

type RunModel struct {
	gorm.Model

	ExternalID string
	Refresh    bool
	Message    string
	Status     string

	ots.RunActions          `gorm:"embedded;embeddedPrefix:action_"`
	ots.RunPermissions      `gorm:"embedded;embeddedPrefix:permission_"`
	ots.RunStatusTimestamps `gorm:"embedded;embeddedPrefix:timestamp_"`

	WorkspaceID uint
	Workspace   WorkspaceModel

	ConfigurationVersionID uint
	ConfigurationVersion   ConfigurationVersionModel

	PlanID uint
	Plan   PlanModel

	ApplyID uint
	Apply   ApplyModel
}

type RunService struct {
	*gorm.DB
}

func NewRunService(db *gorm.DB) *RunService {
	db.AutoMigrate(&RunModel{})

	return &RunService{
		DB: db,
	}
}

func NewRunFromModel(model *RunModel) *ots.Run {
	return &ots.Run{
		ID:               model.ExternalID,
		Refresh:          model.Refresh,
		Message:          model.Message,
		Status:           tfe.RunStatus(model.Status),
		Actions:          &model.RunActions,
		Permissions:      &model.RunPermissions,
		StatusTimestamps: &model.RunStatusTimestamps,
		Plan:             NewPlanFromModel(&model.Plan),
		Apply:            NewApplyFromModel(&model.Apply),
	}
}

func (RunModel) TableName() string {
	return "runs"
}

func (s RunService) CreateRun(opts *tfe.RunCreateOptions) (*ots.Run, error) {
	ws, err := getWorkspaceByID(s.DB, opts.Workspace.ID)
	if err != nil {
		return nil, err
	}

	// If CV ID not provided then get workspace's latest CV
	var cv *ConfigurationVersionModel
	if opts.ConfigurationVersion != nil {
		cv, err = getConfigurationVersionByID(s.DB, opts.ConfigurationVersion.ID)
		if err != nil {
			return nil, err
		}
	} else {
		cv, err = getMostRecentConfigurationVersion(s.DB, ws.ID)
		if err != nil {
			return nil, err
		}
	}

	// TODO: wrap in TX
	plan, err := createPlan(s.DB)
	if err != nil {
		return nil, err
	}

	// TODO: wrap in TX
	apply, err := createApply(s.DB)
	if err != nil {
		return nil, err
	}

	model := RunModel{
		ExternalID: ots.NewRunID(),
		Refresh:    ots.DefaultRefresh,
		RunPermissions: ots.RunPermissions{
			CanApply: true,
		},
		Status: string(tfe.RunPlanQueued),
		RunStatusTimestamps: ots.RunStatusTimestamps{
			QueuedAt: time.Now(),
		},
		ConfigurationVersionID: cv.ID,
		WorkspaceID:            ws.ID,
		Plan:                   *plan,
		Apply:                  *apply,
	}

	if opts.Message != nil {
		model.Message = *opts.Message
	}

	if opts.Refresh != nil {
		model.Refresh = *opts.Refresh
	}

	if result := s.DB.Create(&model); result.Error != nil {
		return nil, result.Error
	}

	return NewRunFromModel(&model), nil
}

func (s RunService) ApplyRun(id string, opts *tfe.RunApplyOptions) error {
	return nil
}

func (s RunService) ListRuns(workspaceID string, opts ots.RunListOptions) (*ots.RunList, error) {
	ws, err := getWorkspaceByID(s.DB, workspaceID)
	if err != nil {
		return nil, err
	}

	var models []RunModel

	if result := s.DB.Where("workspace_id = ?", ws.ID).Limit(opts.PageSize).Offset((opts.PageNumber - 1) * opts.PageSize).Find(&models); result.Error != nil {
		return nil, result.Error
	}

	runs := &ots.RunList{
		RunListOptions: ots.RunListOptions{
			ListOptions: opts.ListOptions,
		},
	}
	for _, m := range models {
		runs.Items = append(runs.Items, NewRunFromModel(&m))
	}

	return runs, nil
}

func (s RunService) GetRun(id string) (*ots.Run, error) {
	model, err := getRunByID(s.DB, id)
	if err != nil {
		return nil, err
	}
	return NewRunFromModel(model), nil
}

func (s RunService) GetQueuedRuns(opts ots.RunListOptions) (*ots.RunList, error) {
	var models []RunModel

	if result := s.DB.Where("status = ?", tfe.RunPlanQueued).Limit(opts.PageSize).Offset((opts.PageNumber - 1) * opts.PageSize).Find(&models); result.Error != nil {
		return nil, result.Error
	}

	runs := &ots.RunList{
		RunListOptions: ots.RunListOptions{
			ListOptions: opts.ListOptions,
		},
	}
	for _, m := range models {
		runs.Items = append(runs.Items, NewRunFromModel(&m))
	}

	return runs, nil
}

func (s RunService) DiscardRun(id string, opts *tfe.RunDiscardOptions) error {
	return nil
}

func (s RunService) CancelRun(id string, opts *tfe.RunCancelOptions) error {
	return nil
}

func (s RunService) ForceCancelRun(id string, opts *tfe.RunForceCancelOptions) error {
	return nil
}

func getRunByID(db *gorm.DB, id string) (*RunModel, error) {
	var model RunModel

	if result := db.Preload(clause.Associations).Where("external_id = ?", id).First(&model); result.Error != nil {
		return nil, result.Error
	}

	return &model, nil
}
