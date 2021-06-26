package agent

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
)

const (
	LocalStateFilename = "terraform.tfstate"
)

type Job struct {
	Run *ots.Run

	Path string

	*Agent
}

// Process processes a job
func (j *Job) Process(ctx context.Context) error {
	// Download config
	if err := j.downloadConfig(ctx); err != nil {
		return err
	}

	// Zap backend config
	if err := j.deleteBackendConfig(ctx); err != nil {
		return err
	}

	// Download state
	if err := j.downloadState(ctx); err != nil {
		return err
	}

	// Update status
	_, err := j.PlanService.UpdatePlanStatus(j.Run.Plan.ID, tfe.PlanRunning)
	if err != nil {
		return err
	}

	// Run terraform init then plan
	out, err := j.TerraformRunner.Plan(ctx)
	if err != nil {
		return err
	}

	// Upload logs
	if err := j.PlanService.UploadPlanLogs(j.Run.Plan.ID, out); err != nil {
		return err
	}

	// Update status
	_, err = j.PlanService.FinishPlan(j.Run.Plan.ID, &ots.PlanFinishOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) downloadConfig(ctx context.Context) error {
	// Download config
	cv, err := j.ConfigurationService.DownloadConfigurationVersion(j.Run.ConfigurationVersion.ID)
	if err != nil {
		return err
	}

	// Decompress and untar config
	if err := Unpack(bytes.NewBuffer(cv), j.Path); err != nil {
		return err
	}

	return nil
}

func (j *Job) deleteBackendConfig(ctx context.Context) error {
	filepath.Walk(j.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		if filepath.Ext(path) != "tf" {
			return nil
		}

		in, err := os.ReadFile(filepath.Join(j.Path, path))
		if err != nil {
			return nil
		}

		deleted, out, err := deleteBackendConfig(in)
		if err != nil {
			return nil
		}

		if deleted {
			if err := os.WriteFile(filepath.Join(j.Path, path), out, 0644); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func (j *Job) downloadState(ctx context.Context) error {
	// Download state
	state, err := j.StateVersionService.CurrentStateVersion(j.Run.Workspace.ID)
	if err != nil {
		return err
	}

	statefile, err := j.StateVersionService.DownloadStateVersion(state.ID)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(j.Path, LocalStateFilename), statefile, 0644); err != nil {
		return err
	}

	return nil
}
