package ots

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
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
	NewDeleteBackendStep = NewFuncStep(deleteBackendConfigFromDirectory)
	NewInitStep          = NewCommandStep("terraform", "init", "-no-color")
	NewPlanStep          = NewCommandStep("terraform", "plan", "-no-color", fmt.Sprintf("-out=%s", PlanFilename))
	NewJSONPlanStep      = NewCommandStep("sh", "-c", fmt.Sprintf("terraform show -json %s > %s", PlanFilename, JSONPlanFilename))
	ApplyStep            = NewCommandStep("sh", "-c", fmt.Sprintf("terraform apply -no-color %s | tee %s", PlanFilename, ApplyOutputFilename))
)

func DownloadConfigStep(ctx context.Context, path string, job *Job, svc StepService) error {
	// Download config
	cv, err := svc.DownloadConfig(job.ConfigurationVersion.ID)
	if err != nil {
		return fmt.Errorf("unable to download config: %w", err)
	}

	// Decompress and untar config
	if err := Unpack(bytes.NewBuffer(cv), path); err != nil {
		return fmt.Errorf("unable to unpack config: %w", err)
	}

	return nil
}

func UploadPlanStep(ctx context.Context, path string, job *Job, svc StepService) error {
	file, err := os.ReadFile(filepath.Join(path, PlanFilename))
	if err != nil {
		return err
	}

	return svc.UploadPlanFile(job.Run.ID, file, false)
}

func UploadJSONPlanStep(ctx context.Context, path string, job *Job, svc StepService) error {
	file, err := os.ReadFile(filepath.Join(path, JSONPlanFilename))
	if err != nil {
		return err
	}

	return svc.UploadPlanFile(job.Run.ID, file, true)
}

// DownloadStateStep downloads current state to disk. If there is no state yet
// nothing will be downloaded and no error will be reported.
func DownloadStateStep(ctx context.Context, path string, job *Job, svc StepService) error {
	state, err := svc.GetCurrentState(job.Workspace.ID)
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
}

func DownloadPlanFileStep(ctx context.Context, path string, job *Job, svc StepService) error {
	plan, err := svc.GetPlanFile(job.Run.ID)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path, PlanFilename), plan, 0644)
}

// UploadStateStep reads, parses, and uploads state
func UploadStateStep(ctx context.Context, path string, job *Job, svc StepService) error {
	stateFile, err := os.ReadFile(filepath.Join(path, LocalStateFilename))
	if err != nil {
		return err
	}

	state, err := Parse(stateFile)
	if err != nil {
		return err
	}

	_, err = svc.CreateStateVersion(job.Workspace.ID, tfe.StateVersionCreateOptions{
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
