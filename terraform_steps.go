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

	"github.com/go-logr/logr"
	"github.com/leg100/go-tfe"
)

const (
	LocalStateFilename  = "terraform.tfstate"
	PlanFilename        = "plan.out"
	JSONPlanFilename    = "plan.out.json"
	ApplyOutputFilename = "apply.out"
)

var (
	NewDeleteBackendStep = NewFuncStep(deleteBackendConfigFromDirectory)
	NewInitStep          = NewCommandStep("terraform", "init", "-no-color")
	NewPlanStep          = NewCommandStep("terraform", "plan", "-no-color", fmt.Sprintf("-out=%s", PlanFilename))
	NewJSONPlanStep      = NewCommandStep("sh", "-c", fmt.Sprintf("terraform show -json %s > %s", PlanFilename, JSONPlanFilename))
	ApplyStep            = NewCommandStep("sh", "-c", fmt.Sprintf("terraform apply -no-color %s | tee %s", PlanFilename, ApplyOutputFilename))
)

func NewDownloadConfigStep(run *Run) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
		// Download config
		cv, err := svc.DownloadConfig(run.ConfigurationVersion.ID)
		if err != nil {
			return fmt.Errorf("unable to download config: %w", err)
		}

		// Decompress and untar config
		if err := Unpack(bytes.NewBuffer(cv), path); err != nil {
			return fmt.Errorf("unable to unpack config: %w", err)
		}

		return nil
	})
}

func NewUpdatePlanStatusStep(run *Run, status tfe.PlanStatus) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
		_, err := svc.UpdatePlanStatus(run.ID, tfe.PlanRunning)
		return err
	})
}

func UpdateApplyStatusStep(run *Run, status tfe.ApplyStatus) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
		_, err := svc.UpdateApplyStatus(run.ID, tfe.ApplyRunning)
		return err
	})
}

func NewFinishPlanStep(run *Run, rs RunService, logger logr.Logger) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
		file, err := os.ReadFile(filepath.Join(path, PlanFilename))
		if err != nil {
			return err
		}
		jsonFile, err := os.ReadFile(filepath.Join(path, JSONPlanFilename))
		if err != nil {
			return err
		}

		planFile := PlanFile{}
		if err := json.Unmarshal(jsonFile, &planFile); err != nil {
			return err
		}

		// Parse plan output
		adds, updates, deletes := planFile.Changes()

		// Update status
		_, err = rs.FinishPlan(run.ID, PlanFinishOptions{
			ResourceAdditions:    adds,
			ResourceChanges:      updates,
			ResourceDestructions: deletes,
			Plan:                 file,
			PlanJSON:             jsonFile,
		})
		if err != nil {
			return fmt.Errorf("unable to finish plan: %w", err)
		}

		logger.Info("job completed", "run", run.ID,
			"additions", adds,
			"changes", updates,
			"deletions", deletes,
		)

		return nil
	})
}

func FinishApplyStep(run *Run, rs RunService, logger logr.Logger) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
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
		_, err = rs.FinishApply(run.ID, ApplyFinishOptions{
			ResourceAdditions:    info.adds,
			ResourceChanges:      info.changes,
			ResourceDestructions: info.deletions,
		})
		if err != nil {
			return fmt.Errorf("unable to finish apply: %w", err)
		}

		logger.Info("job completed", "run", run.ID,
			"additions", info.adds,
			"changes", info.changes,
			"deletions", info.deletions)

		return nil
	})
}

// NewDownloadStateStep downloads current state to disk. If there is no state yet
// nothing will be downloaded and no error will be reported.
func NewDownloadStateStep(run *Run) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
		state, err := svc.GetCurrentState(run.Workspace.ID)
		if IsNotFound(err) {
			return nil
		} else if err != nil {
			return err
		}

		statefile, err := svc.DownloadState(state.ID)
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(path, LocalStateFilename), statefile, 0644); err != nil {
			return err
		}

		return nil
	})
}

func DownloadPlanFileStep(run *Run, rs RunService) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
		plan, err := rs.GetPlanFile(run.ID)
		if err != nil {
			return err
		}

		return os.WriteFile(filepath.Join(path, PlanFilename), plan, 0644)
	})
}

// UploadStateStep reads, parses, and uploads state
func UploadStateStep(run *Run, svs StateVersionService) *FuncStep {
	return NewFuncStep(func(ctx context.Context, path string, svc StepService) error {
		stateFile, err := os.ReadFile(filepath.Join(path, LocalStateFilename))
		if err != nil {
			return err
		}

		state, err := Parse(stateFile)
		if err != nil {
			return err
		}

		_, err = svs.Create(run.Workspace.ID, tfe.StateVersionCreateOptions{
			State:   String(base64.StdEncoding.EncodeToString(stateFile)),
			MD5:     String(fmt.Sprintf("%x", md5.Sum(stateFile))),
			Lineage: &state.Lineage,
			Serial:  Int64(state.Serial),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
