package otf

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/leg100/otf/rbac"
)

const (
	PlanFormatBinary = "bin"  // plan file in binary format
	PlanFormatJSON   = "json" // plan file in json format

	// defaultRefresh specifies that the state be refreshed prior to running a
	// plan
	defaultRefresh = true
)

var (
	ErrRunDiscardNotAllowed      = errors.New("run was not paused for confirmation or priority; discard not allowed")
	ErrRunCancelNotAllowed       = errors.New("run was not planning or applying; cancel not allowed")
	ErrRunForceCancelNotAllowed  = errors.New("run was not planning or applying, has not been canceled non-forcefully, or the cool-off period has not yet passed")
	ErrInvalidRunGetOptions      = errors.New("invalid run get options")
	ErrInvalidRunStateTransition = errors.New("invalid run state transition")
)

// PlanFormat is the format of the plan file
type (
	PlanFormat string

	// Run is a terraform run.
	Run struct {
		ID                     string
		CreatedAt              time.Time
		IsDestroy              bool
		ForceCancelAvailableAt *time.Time
		Message                string
		Organization           string
		Refresh                bool
		RefreshOnly            bool
		ReplaceAddrs           []string
		PositionInQueue        int
		TargetAddrs            []string
		AutoApply              bool
		Speculative            bool
		Status                 RunStatus
		StatusTimestamps       []RunStatusTimestamp
		WorkspaceID            string
		ConfigurationVersionID string
		ExecutionMode          ExecutionMode
		Plan                   Phase
		Apply                  Phase

		Latest bool    // is latest run for workspace
		Commit *string // commit sha that triggered this run
	}

	// RunCreateOptions represents the options for creating a new run. See
	// dto.RunCreateOptions for further detail.
	RunCreateOptions struct {
		IsDestroy              *bool
		Refresh                *bool
		RefreshOnly            *bool
		Message                *string
		ConfigurationVersionID *string
		TargetAddrs            []string
		ReplaceAddrs           []string
		AutoApply              *bool
	}
)

// NewRun creates a new run with defaults.
func NewRun(cv *ConfigurationVersion, ws *Workspace, opts RunCreateOptions) *Run {
	run := Run{
		ID:                     NewID("run"),
		CreatedAt:              CurrentTimestamp(),
		Refresh:                defaultRefresh,
		Organization:           ws.Organization,
		ConfigurationVersionID: cv.ID,
		WorkspaceID:            ws.ID,
		Speculative:            cv.Speculative,
		ReplaceAddrs:           opts.ReplaceAddrs,
		TargetAddrs:            opts.TargetAddrs,
		ExecutionMode:          ws.ExecutionMode,
		AutoApply:              ws.AutoApply,
	}
	run.Plan = NewPhase(run.ID, PlanPhase)
	run.Apply = NewPhase(run.ID, ApplyPhase)
	run.updateStatus(RunPending)

	if opts.IsDestroy != nil {
		run.IsDestroy = *opts.IsDestroy
	}
	if opts.Message != nil {
		run.Message = *opts.Message
	}
	if opts.Refresh != nil {
		run.Refresh = *opts.Refresh
	}
	if opts.AutoApply != nil {
		run.AutoApply = *opts.AutoApply
	}
	if cv.IngressAttributes != nil {
		run.Commit = &cv.IngressAttributes.CommitSHA
	}
	return &run
}

func (r *Run) Queued() bool {
	return r.Status == RunPlanQueued || r.Status == RunApplyQueued
}

func (r *Run) HasChanges() bool {
	return r.Plan.HasChanges()
}

func (r *Run) PlanOnly() bool {
	return r.Status == RunPlannedAndFinished
}

// HasApply determines whether the run has started applying yet.
func (r *Run) HasApply() bool {
	_, err := r.Apply.StatusTimestamp(PhaseRunning)
	return err == nil
}

// Phase returns the current phase.
func (r *Run) Phase() PhaseType {
	switch r.Status {
	case RunPending:
		return PendingPhase
	case RunPlanQueued, RunPlanning, RunPlanned:
		return PlanPhase
	case RunApplyQueued, RunApplying, RunApplied:
		return ApplyPhase
	default:
		return UnknownPhase
	}
}

// Discard updates the state of a run to reflect it having been discarded.
func (r *Run) Discard() error {
	if !r.Discardable() {
		return ErrRunDiscardNotAllowed
	}
	r.updateStatus(RunDiscarded)

	if r.Status == RunPending {
		r.Plan.UpdateStatus(PhaseUnreachable)
	}
	r.Apply.UpdateStatus(PhaseUnreachable)

	return nil
}

// Cancel run. Returns a boolean indicating whether a cancel request should be
// enqueued (for an agent to kill an in progress process)
func (r *Run) Cancel() (enqueue bool, err error) {
	if !r.Cancelable() {
		return false, ErrRunCancelNotAllowed
	}
	// permit run to be force canceled after a cool off period of 10 seconds has
	// elapsed.
	tenSecondsFromNow := CurrentTimestamp().Add(10 * time.Second)
	r.ForceCancelAvailableAt = &tenSecondsFromNow

	switch r.Status {
	case RunPending:
		r.Plan.UpdateStatus(PhaseUnreachable)
		r.Apply.UpdateStatus(PhaseUnreachable)
	case RunPlanQueued, RunPlanning:
		r.Plan.UpdateStatus(PhaseCanceled)
		r.Apply.UpdateStatus(PhaseUnreachable)
	case RunApplyQueued, RunApplying:
		r.Apply.UpdateStatus(PhaseCanceled)
	}

	if r.Status == RunPlanning || r.Status == RunApplying {
		enqueue = true
	}

	r.updateStatus(RunCanceled)

	return enqueue, nil
}

// ForceCancel force cancels a run. A cool-off period of 10 seconds must have
// elapsed following a cancelation request before a run can be force canceled.
func (r *Run) ForceCancel() error {
	if r.ForceCancelAvailableAt != nil && time.Now().After(*r.ForceCancelAvailableAt) {
		r.updateStatus(RunCanceled)
		return nil
	}
	return ErrRunForceCancelNotAllowed
}

// Done determines whether run has reached an end state, e.g. applied,
// discarded, etc.
func (r *Run) Done() bool {
	switch r.Status {
	case RunApplied, RunPlannedAndFinished, RunDiscarded, RunCanceled, RunErrored:
		return true
	default:
		return false
	}
}

// EnqueuePlan enqueues a plan for the run. It also sets the run as the latest
// run for its workspace (speculative runs are ignored).
func (r *Run) EnqueuePlan() error {
	if r.Status != RunPending {
		return fmt.Errorf("cannot enqueue run with status %s", r.Status)
	}
	r.updateStatus(RunPlanQueued)
	r.Plan.UpdateStatus(PhaseQueued)

	return nil
}

func (*Run) CanAccessSite(action rbac.Action) bool {
	// run cannot carry out site-level actions
	return false
}

func (r *Run) CanAccessOrganization(action rbac.Action, name string) bool {
	// run cannot access organization-level resources
	return false
}

func (r *Run) CanAccessWorkspace(action rbac.Action, policy *WorkspacePolicy) bool {
	// run can access anything within its workspace
	return r.WorkspaceID == policy.WorkspaceID
}

func (r *Run) EnqueueApply() error {
	if r.Status != RunPlanned {
		return fmt.Errorf("cannot apply run")
	}
	r.updateStatus(RunApplyQueued)
	r.Apply.UpdateStatus(PhaseQueued)
	return nil
}

func (r *Run) StatusTimestamp(status RunStatus) (time.Time, error) {
	for _, rst := range r.StatusTimestamps {
		if rst.Status == status {
			return rst.Timestamp, nil
		}
	}
	return time.Time{}, ErrStatusTimestampNotFound
}

// Start a run phase
func (r *Run) Start(phase PhaseType) error {
	switch r.Status {
	case RunPlanQueued:
		return r.startPlan()
	case RunApplyQueued:
		return r.startApply()
	case RunPlanning, RunApplying:
		return ErrPhaseAlreadyStarted
	default:
		return ErrInvalidRunStateTransition
	}
}

// Finish updates the run to reflect its plan or apply phase having finished.
func (r *Run) Finish(phase PhaseType, opts PhaseFinishOptions) error {
	if r.Status == RunCanceled {
		// run was canceled before the phase finished so nothing more to do.
		return nil
	}
	switch phase {
	case PlanPhase:
		return r.finishPlan(opts)
	case ApplyPhase:
		return r.finishApply(opts)
	default:
		return fmt.Errorf("unknown phase")
	}
}

func (r *Run) startPlan() error {
	if r.Status != RunPlanQueued {
		return ErrInvalidRunStateTransition
	}
	r.updateStatus(RunPlanning)
	r.Plan.UpdateStatus(PhaseRunning)
	return nil
}

func (r *Run) startApply() error {
	if r.Status != RunApplyQueued {
		return ErrInvalidRunStateTransition
	}
	r.updateStatus(RunApplying)
	r.Apply.UpdateStatus(PhaseRunning)
	return nil
}

func (r *Run) finishPlan(opts PhaseFinishOptions) error {
	if r.Status != RunPlanning {
		return ErrInvalidRunStateTransition
	}
	if opts.Errored {
		r.updateStatus(RunErrored)
		r.Plan.UpdateStatus(PhaseErrored)
		r.Apply.UpdateStatus(PhaseUnreachable)
		return nil
	}

	r.updateStatus(RunPlanned)
	r.Plan.UpdateStatus(PhaseFinished)

	if !r.HasChanges() || r.Speculative {
		r.updateStatus(RunPlannedAndFinished)
		r.Apply.UpdateStatus(PhaseUnreachable)
	} else if r.AutoApply {
		return r.EnqueueApply()
	}
	return nil
}

func (r *Run) finishApply(opts PhaseFinishOptions) error {
	if r.Status != RunApplying {
		return ErrInvalidRunStateTransition
	}
	if opts.Errored {
		r.updateStatus(RunErrored)
		r.Apply.UpdateStatus(PhaseErrored)
	} else {
		r.updateStatus(RunApplied)
		r.Apply.UpdateStatus(PhaseFinished)
	}
	return nil
}

func (r *Run) updateStatus(status RunStatus) {
	r.Status = status
	r.StatusTimestamps = append(r.StatusTimestamps, RunStatusTimestamp{
		Status:    status,
		Timestamp: CurrentTimestamp(),
	})
}

// Discardable determines whether run can be discarded.
func (r *Run) Discardable() bool {
	switch r.Status {
	case RunPending, RunPlanned:
		return true
	default:
		return false
	}
}

// Cancelable determines whether run can be cancelled.
func (r *Run) Cancelable() bool {
	switch r.Status {
	case RunPending, RunPlanQueued, RunPlanning, RunPlanned, RunApplyQueued, RunApplying:
		return true
	default:
		return false
	}
}

// Confirmable determines whether run can be confirmed.
func (r *Run) Confirmable() bool {
	switch r.Status {
	case RunPlanned:
		return true
	default:
		return false
	}
}

// RunList represents a list of runs.
type RunList struct {
	*Pagination
	Items []*Run
}

// RunListOptions are options for paginating and filtering a list of runs
type RunListOptions struct {
	ListOptions
	// Filter by run statuses (with an implicit OR condition)
	Statuses []RunStatus `schema:"statuses,omitempty"`
	// Filter by workspace ID
	WorkspaceID *string `schema:"workspace_id,omitempty"`
	// Filter by organization name
	Organization *string `schema:"organization_name,omitempty"`
	// Filter by workspace name
	WorkspaceName *string `schema:"workspace_name,omitempty"`
	// Filter by speculative or non-speculative
	Speculative *bool `schema:"-"`
	// A list of relations to include. See available resources:
	// https://www.terraform.io/docs/cloud/api/run.html#available-related-resources
	Include *string `schema:"include,omitempty"`
}

// type (
//
//	Plan interface {
//		ResourceReport() *ResourceReport
//	}
//	Apply interface {
//		ResourceReport() *ResourceReport
//	}
//
// )
type RunDB interface {
	GetRun(context.Context, string) (Run, error)
}

// RunService implementations allow interactions with runs
type RunService interface {
	// Create a new run with the given options.
	// CreateRun(ctx context.Context, workspaceID string, opts RunCreateOptions) (*Run, error)
	// Get retrieves a run with the given ID.
	GetRun(ctx context.Context, id string) (Run, error)
	// List lists runs according to the given options.
	//ListRuns(ctx context.Context, opts RunListOptions) (*RunList, error)
	// Delete deletes a run with the given ID.
	//DeleteRun(ctx context.Context, id string) error
	// EnqueuePlan enqueues a plan
	//EnqueuePlan(ctx context.Context, id string) (*Run, error)
	// Apply a run with the given ID.
	//
	// TODO: return run
	//ApplyRun(ctx context.Context, id string, opts RunApplyOptions) error
	// Discard discards a run with the given ID.
	//
	// TODO: return run
	//DiscardRun(ctx context.Context, id string, opts RunDiscardOptions) error
	// Cancel run.
	//
	// TODO: return run
	//CancelRun(ctx context.Context, id string, opts RunCancelOptions) error
	// Forcefully cancel a run.
	//
	// TODO: return run
	//ForceCancelRun(ctx context.Context, id string, opts RunForceCancelOptions) error
	// Start a run phase.
	//StartPhase(ctx context.Context, id string, phase PhaseType, opts PhaseStartOptions) (*Run, error)
	// Finish a run phase.
	//FinishPhase(ctx context.Context, id string, phase PhaseType, opts PhaseFinishOptions) (*Run, error)
	// GetPlanFile retrieves a run's plan file with the requested format.
	//GetPlanFile(ctx context.Context, id string, format PlanFormat) ([]byte, error)
	// UploadPlanFile saves a run's plan file with the requested format.
	//UploadPlanFile(ctx context.Context, id string, plan []byte, format PlanFormat) error
	// GetLockFile retrieves a run's lock file (.terraform.lock.hcl)
	//GetLockFile(ctx context.Context, id string) ([]byte, error)
	// UploadLockFile saves a run's lock file (.terraform.lock.hcl)
	//UploadLockFile(ctx context.Context, id string, lockFile []byte) error
	// StartRun creates and starts a run.
	//StartRun(ctx context.Context, workspaceID string, opts ConfigurationVersionCreateOptions) (*Run, error)
}

type RunStatusTimestamp struct {
	Status    RunStatus
	Timestamp time.Time
}
