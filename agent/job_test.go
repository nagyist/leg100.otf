package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockConfigurationVersionService struct {
	ots.ConfigurationVersionService
}

func (s *MockConfigurationVersionService) DownloadConfigurationVersion(id string) ([]byte, error) {
	return os.ReadFile("testdata/unpack.tar.gz")
}

type MockStateVersionService struct {
	ots.StateVersionService
}

func (s *MockStateVersionService) CurrentStateVersion(workspaceID string) (*ots.StateVersion, error) {
	return &ots.StateVersion{ID: "sv-123"}, nil
}

func (s *MockStateVersionService) DownloadStateVersion(id string) ([]byte, error) {
	return []byte("state"), nil
}

type MockPlanService struct {
	ots.PlanService
}

func (s *MockPlanService) UploadPlanLogs(id string, logs []byte) error { return nil }

func (s *MockPlanService) UpdatePlanStatus(id string, status tfe.PlanStatus) (*ots.Plan, error) {
	return nil, nil
}

func (s *MockPlanService) FinishPlan(id string, opts *ots.PlanFinishOptions) (*ots.Plan, error) {
	return nil, nil
}

type MockRunner struct{}

func (r *MockRunner) Plan(ctx context.Context) ([]byte, error) {
	return os.ReadFile("testdata/init.log")
}

func TestJob(t *testing.T) {
	path := t.TempDir()

	job := Job{
		ConfigurationService: &MockConfigurationVersionService{},
		StateVersionService:  &MockStateVersionService{},
		PlanService:          &MockPlanService{},
		TerraformRunner:      &MockRunner{},
		Path:                 path,
		Run: ots.Run{
			Plan: &ots.Plan{
				ID: "plan-123",
			},
			ConfigurationVersion: &ots.ConfigurationVersion{
				ID: "cv-123",
			},
			Workspace: &ots.Workspace{
				ID: "ws-123",
			},
		},
	}

	require.NoError(t, job.Process(context.Background()))

	var got []string
	filepath.Walk(path, func(lpath string, info os.FileInfo, err error) error {
		lpath, err = filepath.Rel(path, lpath)
		require.NoError(t, err)
		got = append(got, lpath)
		return nil
	})
	assert.Equal(t, []string{
		".",
		"dir",
		"dir/file",
		"dir/symlink",
		"file",
		"terraform.tfstate",
	}, got)
}
