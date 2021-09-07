package ots

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leg100/go-tfe"
)

const (
	LocalStateFilename  = "terraform.tfstate"
	PlanFilename        = "plan.out"
	JSONPlanFilename    = "plan.out.json"
	ApplyOutputFilename = "apply.out"
)

var (
	DeleteBackendStep = NewFuncStep(deleteBackendConfigFromDirectory)
	InitStep          = NewCommandStep("terraform", "init", "-no-color")
	PlanStep          = NewCommandStep("terraform", "plan", "-no-color", fmt.Sprintf("-out=%s", PlanFilename))
	JSONPlanStep      = NewCommandStep("sh", "-c", fmt.Sprintf("terraform show -json %s > %s", PlanFilename, JSONPlanFilename))
	ApplyStep         = NewCommandStep("sh", "-c", fmt.Sprintf("terraform apply -no-color %s | tee %s", PlanFilename, ApplyOutputFilename))
)

type PlanStepsProvider struct{}

type ApplyStepsProvider struct{}

func (*PlanStepsProvider) String() string { return "plan" }

func (*PlanStepsProvider) Steps(rs RunService, cvs ConfigurationVersionService, svs StateVersionService) []Step {
	ts := &TerraformSteps{
		ConfigurationVersionService: cvs,
		RunService:                  rs,
		StateVersionService:         svs,
	}

	return []Step{
		NewFuncStep(ts.DownloadConfigStep),
		DeleteBackendStep,
		NewFuncStep(ts.DownloadStateStep),
		InitStep,
		PlanStep,
		JSONPlanStep,
		NewFuncStep(ts.UploadPlanStep),
		NewFuncStep(ts.UploadJSONPlanStep),
		NewFuncStep(ts.SummarizePlanStep),
	}
}

func (*ApplyStepsProvider) Steps(rs RunService, cvs ConfigurationVersionService, svs StateVersionService) []Step {
	ts := &TerraformSteps{
		ConfigurationVersionService: cvs,
		RunService:                  rs,
		StateVersionService:         svs,
	}

	return []Step{
		NewFuncStep(ts.DownloadConfigStep),
		DeleteBackendStep,
		NewFuncStep(ts.DownloadPlanFileStep),
		NewFuncStep(ts.DownloadStateStep),
		InitStep,
		NewFuncStep(ts.UploadStateStep),
		NewFuncStep(ts.SummarizeApplyStep),
	}
}

func (*ApplyStepsProvider) String() string { return "apply" }

type TerraformSteps struct {
	ConfigurationVersionService
	RunService
	StateVersionService
}

func (ts *TerraformSteps) DownloadConfigStep(ctx context.Context, path string, job *Job) error {
	// Download config
	cv, err := ts.ConfigurationVersionService.Download(job.ConfigurationVersionID)
	if err != nil {
		return fmt.Errorf("unable to download config: %w", err)
	}

	// Decompress and untar config
	if err := Unpack(bytes.NewBuffer(cv), path); err != nil {
		return fmt.Errorf("unable to unpack config: %w", err)
	}

	return nil
}

func (ts *TerraformSteps) UploadPlanStep(ctx context.Context, path string, job *Job) error {
	file, err := os.ReadFile(filepath.Join(path, PlanFilename))
	if err != nil {
		return err
	}

	return ts.RunService.UploadPlan(job.RunID, file, false)
}

func (ts *TerraformSteps) UploadJSONPlanStep(ctx context.Context, path string, job *Job) error {
	file, err := os.ReadFile(filepath.Join(path, JSONPlanFilename))
	if err != nil {
		return err
	}

	return ts.RunService.UploadPlan(job.RunID, file, true)
}

// DownloadStateStep downloads current state to disk. If there is no state yet
// nothing will be downloaded and no error will be reported.
func (ts *TerraformSteps) DownloadStateStep(ctx context.Context, path string, job *Job) error {
	state, err := ts.StateVersionService.Current(job.WorkspaceID)
	if IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	statefile, err := ts.StateVersionService.Download(state.ID)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(path, LocalStateFilename), statefile, 0644); err != nil {
		return err
	}

	return nil
}

func (ts *TerraformSteps) DownloadPlanFileStep(ctx context.Context, path string, job *Job) error {
	plan, err := ts.GetPlanFile(job.RunID)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path, PlanFilename), plan, 0644)
}

// UploadStateStep reads, parses, and uploads state
func (ts *TerraformSteps) UploadStateStep(ctx context.Context, path string, job *Job) error {
	stateFile, err := os.ReadFile(filepath.Join(path, LocalStateFilename))
	if err != nil {
		return err
	}

	state, err := Parse(stateFile)
	if err != nil {
		return err
	}

	_, err = ts.StateVersionService.Create(job.WorkspaceID, tfe.StateVersionCreateOptions{
		State:   String(base64.StdEncoding.EncodeToString(stateFile)),
		MD5:     String(fmt.Sprintf("%x", md5.Sum(stateFile))),
		Lineage: &state.Lineage,
		Serial:  Int64(state.Serial),
	})
	return err
}

func (ts *TerraformSteps) SummarizePlanStep(ctx context.Context, path string, job *Job) error {
	jsonPlan, err := os.ReadFile(filepath.Join(path, JSONPlanFilename))
	if err != nil {
		return err
	}

	// Parse plan file
	planFile := &PlanFile{}
	if err := json.Unmarshal(jsonPlan, planFile); err != nil {
		return err
	}
	adds, changes, deletions := planFile.Changes()

	// Update status
	return ts.RunService.UpdatePlanSummary(job.RunID, ResourceSummary{
		ResourceAdditions:    adds,
		ResourceChanges:      changes,
		ResourceDestructions: deletions,
	})
}

func (ts *TerraformSteps) SummarizeApplyStep(ctx context.Context, path string, job *Job) error {
	out, err := os.ReadFile(filepath.Join(path, ApplyOutputFilename))
	if err != nil {
		return err
	}

	// Parse apply output
	info, err := parseApplyOutput(string(out))
	if err != nil {
		return fmt.Errorf("unable to parse apply output: %w", err)
	}

	// Update status
	return ts.RunService.UpdateApplySummary(job.RunID, ResourceSummary{
		ResourceAdditions:    info.adds,
		ResourceChanges:      info.changes,
		ResourceDestructions: info.deletions,
	})
}
