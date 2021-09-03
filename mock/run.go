package mock

import (
	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
)

var _ ots.RunService = (*RunService)(nil)

type RunService struct {
	CreateFn             func(opts *tfe.RunCreateOptions) (*ots.Run, error)
	GetFn                func(id string) (*ots.Run, error)
	ListFn               func(opts ots.RunListOptions) (*ots.RunList, error)
	ApplyFn              func(id string, opts *tfe.RunApplyOptions) error
	DiscardFn            func(id string, opts *tfe.RunDiscardOptions) error
	CancelFn             func(id string, opts *tfe.RunCancelOptions) error
	ForceCancelFn        func(id string, opts *tfe.RunForceCancelOptions) error
	EnqueuePlanFn        func(id string) error
	UpdateStatusFn       func(id string, status tfe.RunStatus) (*ots.Run, error)
	GetPlanJSONFn        func(id string) ([]byte, error)
	GetPlanFileFn        func(id string) ([]byte, error)
	UploadPlanFn         func(id string, plan []byte, json bool) error
	UpdatePlanSummaryFn  func(id string, summary ots.ResourceSummary) error
	UpdateApplySummaryFn func(id string, summary ots.ResourceSummary) error

	StartJobFn      func(id string, opts ots.JobStartOptions) error
	FinishJobFn     func(id string, opts ots.JobFinishOptions) error
	ListJobsFn      func(opts ots.JobListOptions) ([]*ots.Job, error)
	UploadJobLogsFn func(id string, logs []byte) error
	GetJobLogsFn    func(id string, opts ots.JobLogOptions) ([]byte, error)
}

func (s RunService) Create(opts *tfe.RunCreateOptions) (*ots.Run, error) {
	return s.CreateFn(opts)
}

func (s RunService) Get(id string) (*ots.Run, error) {
	return s.GetFn(id)
}

func (s RunService) List(opts ots.RunListOptions) (*ots.RunList, error) {
	return s.ListFn(opts)
}

func (s RunService) Apply(id string, opts *tfe.RunApplyOptions) error {
	return s.ApplyFn(id, opts)
}

func (s RunService) Discard(id string, opts *tfe.RunDiscardOptions) error {
	return s.DiscardFn(id, opts)
}

func (s RunService) Cancel(id string, opts *tfe.RunCancelOptions) error {
	return s.CancelFn(id, opts)
}

func (s RunService) ForceCancel(id string, opts *tfe.RunForceCancelOptions) error {
	return s.ForceCancelFn(id, opts)
}

func (s RunService) EnqueuePlan(id string) error {
	return s.EnqueuePlanFn(id)
}

func (s RunService) UpdateStatus(id string, status tfe.RunStatus) (*ots.Run, error) {
	return s.UpdateStatusFn(id, status)
}

func (s RunService) GetPlanJSON(id string) ([]byte, error) {
	return s.GetPlanJSONFn(id)
}

func (s RunService) GetPlanFile(id string) ([]byte, error) {
	return s.GetPlanFileFn(id)
}

func (s RunService) UploadPlan(id string, plan []byte, json bool) error {
	return s.UploadPlanFn(id, plan, json)
}

func (s RunService) UpdatePlanSummary(id string, summary ots.ResourceSummary) error {
	return s.UpdatePlanSummaryFn(id, summary)
}

func (s RunService) UpdateApplySummary(id string, summary ots.ResourceSummary) error {
	return s.UpdateApplySummaryFn(id, summary)
}

func (s RunService) StartJob(id string, opts ots.JobStartOptions) error {
	return s.StartJobFn(id, opts)
}

func (s RunService) FinishJob(id string, opts ots.JobFinishOptions) error {
	return s.FinishJobFn(id, opts)
}

func (s RunService) ListJobs(opts ots.JobListOptions) ([]*ots.Job, error) {
	return s.ListJobsFn(opts)
}

func (s RunService) UploadJobLogs(id string, logs []byte) error {
	return s.UploadJobLogsFn(id, logs)
}

func (s RunService) GetJobLogs(id string, opts ots.JobLogOptions) ([]byte, error) {
	return s.GetJobLogsFn(id, opts)
}
