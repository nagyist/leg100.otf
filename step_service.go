package ots

import tfe "github.com/leg100/go-tfe"

type StepService interface {
	DownloadConfig(id string) ([]byte, error)
	GetPlanFile(id string) ([]byte, error)
	GetCurrentState(workspaceID string) (*StateVersion, error)
	DownloadState(id string) ([]byte, error)
	CreateStateVersion(workspaceID string, opts tfe.StateVersionCreateOptions) (*StateVersion, error)
	UploadPlanFile(runID string, file []byte, json bool) error
	UpdatePlanSummary(runID string, summary ResourceSummary) error
	UpdateApplySummary(runID string, summary ResourceSummary) error
}

type stepService struct {
	StateVersionService
	ConfigurationVersionService
	RunService
}

func NewStepService(
	rs RunService,
	cvs ConfigurationVersionService,
	svs StateVersionService,
) *stepService {
	return &stepService{
		RunService:                  rs,
		ConfigurationVersionService: cvs,
		StateVersionService:         svs,
	}
}

func (s *stepService) DownloadConfig(id string) ([]byte, error) {
	return s.ConfigurationVersionService.Download(id)
}

func (s *stepService) GetPlanFile(id string) ([]byte, error) {
	return s.RunService.GetPlanFile(id)
}

func (s *stepService) GetCurrentState(workspaceID string) (*StateVersion, error) {
	return s.StateVersionService.Current(workspaceID)
}

func (s *stepService) DownloadState(workspaceID string) ([]byte, error) {
	return s.StateVersionService.Download(workspaceID)
}

func (s *stepService) CreateStateVersion(workspaceID string, opts tfe.StateVersionCreateOptions) (*StateVersion, error) {
	return s.StateVersionService.Current(workspaceID)
}

func (s *stepService) UploadPlanFile(runID string, plan []byte, json bool) error {
	return s.RunService.UploadPlan(runID, plan, json)
}

func (s *stepService) UpdatePlanSummary(runID string, summary ResourceSummary) error {
	return s.RunService.UpdatePlanSummary(runID, summary)
}

func (s *stepService) UpdateApplySummary(runID string, summary ResourceSummary) error {
	return s.RunService.UpdateApplySummary(runID, summary)
}
