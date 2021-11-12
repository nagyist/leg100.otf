package otf

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type UserService interface {
	ListTokens(ctx context.Context) (*TokenList, error)
}

// RunStore implementations persist Run objects.
type RunStore interface {
	Create(run *Run) (*Run, error)
	Get(opts RunGetOptions) (*Run, error)
	List(opts RunListOptions) (*RunList, error)
	// TODO: add support for a special error type that tells update to skip
	// updates - useful when fn checks current fields and decides not to update
	Update(opts RunGetOptions, fn func(*Run) error) (*Run, error)
}

type TokenList struct {
	*Pagination
	Items []*Token
}

// RunGetOptions are options for retrieving a single Run. Either ID or ApplyID
// or PlanID must be specfiied.
type RunGetOptions struct {
	// ID of run to retrieve
	ID *string

	// Get run via apply ID
	ApplyID *string

	// Get run via plan ID
	PlanID *string

	// IncludePlanFile toggles including the plan file in the retrieved run.
	IncludePlanFile bool

	// IncludePlanFile toggles including the plan file, in JSON format, in the
	// retrieved run.
	IncludePlanJSON bool
}

// RunListOptions are options for paginating and filtering a list of runs
type RunListOptions struct {
	ListOptions

	// A list of relations to include. See available resources:
	// https://www.terraform.io/docs/cloud/api/run.html#available-related-resources
	Include *string `schema:"include"`

	// Filter by run statuses (with an implicit OR condition)
	Statuses []RunStatus

	// Filter by workspace ID
	WorkspaceID *string
}

func (r *Run) GetID() string  { return r.ID }
func (r *Run) String() string { return r.ID }

func (o RunCreateOptions) Valid() error {
	if o.Workspace == nil {
		return errors.New("workspace is required")
	}
	return nil
}

func (r *Run) GetStatus() string {
	return string(r.Status)
}

// Discard updates the state of a run to reflect it having been discarded.
func (r *Run) Discard() error {
	if !r.IsDiscardable() {
		return ErrRunDiscardNotAllowed
	}

	r.UpdateStatus(RunDiscarded)

	return nil
}

// Cancel run.
func (r *Run) Cancel() error {
	if !r.IsCancelable() {
		return ErrRunCancelNotAllowed
	}

	r.UpdateStatus(RunCanceled)

	return nil
}

func (r *Run) ForceCancelAvailableAt() time.Time {
	canceledAt, ok := r.StatusTimestamps[string(RunCanceled)]
	if !ok {
		return time.Time{}
	}

	// Run can be forcefully cancelled after a cool-off period of ten seconds
	return canceledAt.Add(10 * time.Second)
}

// ForceCancel updates the state of a run to reflect it having been forcefully
// cancelled.
func (r *Run) ForceCancel() error {
	if !r.IsForceCancelable() {
		return ErrRunForceCancelNotAllowed
	}

	r.setTimestamp(RunForceCanceled)

	return nil
}

// IsCancelable determines whether run can be cancelled.
func (r *Run) IsCancelable() bool {
	switch r.Status {
	case RunPending, RunPlanQueued, RunPlanning, RunApplyQueued, RunApplying:
		return true
	default:
		return false
	}
}

// IsConfirmable determines whether run can be confirmed.
func (r *Run) IsConfirmable() bool {
	switch r.Status {
	case RunPlanned:
		return true
	default:
		return false
	}
}

// IsDiscardable determines whether run can be discarded.
func (r *Run) IsDiscardable() bool {
	switch r.Status {
	case RunPending, RunPlanned:
		return true
	default:
		return false
	}
}

// IsForceCancelable determines whether a run can be forcibly cancelled.
func (r *Run) IsForceCancelable() bool {
	availAt := r.ForceCancelAvailableAt()

	if availAt.IsZero() {
		return false
	}

	return time.Now().After(availAt)
}

// IsActive determines whether run is currently the active run on a workspace,
// i.e. it is neither finished nor pending
func (r *Run) IsActive() bool {
	if r.IsDone() || r.Status == RunPending {
		return false
	}
	return true
}

// IsDone determines whether run has reached an end state, e.g. applied,
// discarded, etc.
func (r *Run) IsDone() bool {
	switch r.Status {
	case RunApplied, RunPlannedAndFinished, RunDiscarded, RunCanceled, RunErrored:
		return true
	default:
		return false
	}
}

func (r *Run) IsSpeculative() bool {
	return r.ConfigurationVersion.Speculative
}

// UpdateStatus updates the status of the run as well as its plan and apply
func (r *Run) UpdateStatus(status RunStatus) {
	switch status {
	case RunPending:
		r.Plan.UpdateStatus(PlanPending)
	case RunPlanQueued:
		r.Plan.UpdateStatus(PlanQueued)
	case RunPlanning:
		r.Plan.UpdateStatus(PlanRunning)
	case RunPlanned, RunPlannedAndFinished:
		r.Plan.UpdateStatus(PlanFinished)
	case RunApplyQueued:
		r.Apply.UpdateStatus(ApplyQueued)
	case RunApplying:
		r.Apply.UpdateStatus(ApplyRunning)
	case RunApplied:
		r.Apply.UpdateStatus(ApplyFinished)
	case RunErrored:
		switch r.Status {
		case RunPlanning:
			r.Plan.UpdateStatus(PlanErrored)
		case RunApplying:
			r.Apply.UpdateStatus(ApplyErrored)
		}
	case RunCanceled:
		switch r.Status {
		case RunPlanQueued, RunPlanning:
			r.Plan.UpdateStatus(PlanCanceled)
		case RunApplyQueued, RunApplying:
			r.Apply.UpdateStatus(ApplyCanceled)
		}
	}

	r.Status = status
	r.setTimestamp(status)

	// TODO: determine when ApplyUnreachable is applicable and set
	// accordingly
}

func (r *Run) setTimestamp(status RunStatus) {
	r.StatusTimestamps[string(status)] = time.Now()
}

// Do invokes the necessary steps before a plan or apply can proceed.
func (r *Run) Do(env Environment) error {
	if err := env.RunFunc(r.downloadConfig); err != nil {
		return err
	}

	err := env.RunFunc(func(ctx context.Context, env Environment) error {
		return deleteBackendConfigFromDirectory(ctx, env.GetPath())
	})
	if err != nil {
		return err
	}

	if err := env.RunFunc(r.downloadState); err != nil {
		return err
	}

	if err := env.RunCLI("terraform", "init"); err != nil {
		return fmt.Errorf("running terraform init: %w", err)
	}

	return nil
}

func (r *Run) downloadConfig(ctx context.Context, env Environment) error {
	// Download config
	cv, err := env.GetConfigurationVersionService().Download(r.ConfigurationVersion.ID)
	if err != nil {
		return fmt.Errorf("unable to download config: %w", err)
	}

	// Decompress and untar config
	if err := Unpack(bytes.NewBuffer(cv), env.GetPath()); err != nil {
		return fmt.Errorf("unable to unpack config: %w", err)
	}

	return nil
}

// downloadState downloads current state to disk. If there is no state yet
// nothing will be downloaded and no error will be reported.
func (r *Run) downloadState(ctx context.Context, env Environment) error {
	state, err := env.GetStateVersionService().Current(r.Workspace.ID)
	if errors.Is(err, ErrResourceNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("retrieving current state version: %w", err)
	}

	statefile, err := env.GetStateVersionService().Download(state.ID)
	if err != nil {
		return fmt.Errorf("downloading state version: %w", err)
	}

	if err := os.WriteFile(filepath.Join(env.GetPath(), LocalStateFilename), statefile, 0644); err != nil {
		return fmt.Errorf("saving state to local disk: %w", err)
	}

	return nil
}

func (r *Run) uploadPlan(ctx context.Context, env Environment) error {
	file, err := os.ReadFile(filepath.Join(env.GetPath(), PlanFilename))
	if err != nil {
		return err
	}

	opts := PlanFileOptions{Format: PlanBinaryFormat}

	if err := env.GetRunService().UploadPlanFile(ctx, r.ID, file, opts); err != nil {
		return fmt.Errorf("unable to upload plan: %w", err)
	}

	return nil
}

func (r *Run) uploadJSONPlan(ctx context.Context, env Environment) error {
	jsonFile, err := os.ReadFile(filepath.Join(env.GetPath(), JSONPlanFilename))
	if err != nil {
		return err
	}

	opts := PlanFileOptions{Format: PlanJSONFormat}

	if err := env.GetRunService().UploadPlanFile(ctx, r.ID, jsonFile, opts); err != nil {
		return fmt.Errorf("unable to upload JSON plan: %w", err)
	}

	return nil
}

func (r *Run) downloadPlanFile(ctx context.Context, env Environment) error {
	opts := PlanFileOptions{Format: PlanBinaryFormat}

	plan, err := env.GetRunService().GetPlanFile(ctx, r.ID, opts)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(env.GetPath(), PlanFilename), plan, 0644)
}

// uploadState reads, parses, and uploads state
func (r *Run) uploadState(ctx context.Context, env Environment) error {
	stateFile, err := os.ReadFile(filepath.Join(env.GetPath(), LocalStateFilename))
	if err != nil {
		return err
	}

	state, err := Parse(stateFile)
	if err != nil {
		return err
	}

	_, err = env.GetStateVersionService().Create(r.Workspace.ID, StateVersionCreateOptions{
		State:   String(base64.StdEncoding.EncodeToString(stateFile)),
		MD5:     String(fmt.Sprintf("%x", md5.Sum(stateFile))),
		Lineage: &state.Lineage,
		Serial:  Int64(state.Serial),
	})
	if err != nil {
		return err
	}

	return nil
}

// NewRun constructs a run object.
func (f *RunFactory) NewRun(opts RunCreateOptions) (*Run, error) {
	if opts.Workspace == nil {
		return nil, errors.New("workspace is required")
	}

	id := NewID("run")
	run := Run{
		ID:               id,
		Timestamps:       NewTimestamps(),
		Refresh:          DefaultRefresh,
		ReplaceAddrs:     opts.ReplaceAddrs,
		TargetAddrs:      opts.TargetAddrs,
		StatusTimestamps: make(TimestampMap),
		Plan:             newPlan(id),
		Apply:            newApply(id),
	}

	run.UpdateStatus(RunPending)

	ws, err := f.WorkspaceService.Get(context.Background(), WorkspaceSpecifier{ID: &opts.Workspace.ID})
	if err != nil {
		return nil, err
	}
	run.Workspace = ws

	cv, err := f.getConfigurationVersion(opts)
	if err != nil {
		return nil, err
	}
	run.ConfigurationVersion = cv

	if opts.IsDestroy != nil {
		run.IsDestroy = *opts.IsDestroy
	}

	if opts.Message != nil {
		run.Message = *opts.Message
	}

	if opts.Refresh != nil {
		run.Refresh = *opts.Refresh
	}

	return &run, nil
}

func (f *RunFactory) getConfigurationVersion(opts RunCreateOptions) (*ConfigurationVersion, error) {
	// Unless CV ID provided, get workspace's latest CV
	if opts.ConfigurationVersion != nil {
		return f.ConfigurationVersionService.Get(opts.ConfigurationVersion.ID)
	}
	return f.ConfigurationVersionService.GetLatest(opts.Workspace.ID)
}
