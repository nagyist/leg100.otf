package run

import (
	"context"
	"testing"

	"github.com/leg100/otf"
	"github.com/leg100/otf/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	ctx := context.Background()

	t.Run("defaults", func(t *testing.T) {
		f := testFactory(
			&workspace.Workspace{},
			&otf.ConfigurationVersion{},
		)

		got, err := f.NewRun(ctx, "", RunCreateOptions{})
		require.NoError(t, err)

		assert.Equal(t, otf.RunPending, got.Status)
		assert.NotZero(t, got.CreatedAt)
		assert.False(t, got.Speculative)
		assert.True(t, got.Refresh)
		assert.False(t, got.AutoApply)
	})

	t.Run("speculative run", func(t *testing.T) {
		f := testFactory(
			&workspace.Workspace{},
			&otf.ConfigurationVersion{Speculative: true},
		)

		got, err := f.NewRun(ctx, "", RunCreateOptions{})
		require.NoError(t, err)

		assert.True(t, got.Speculative)
	})

	t.Run("workspace auto-apply", func(t *testing.T) {
		f := testFactory(
			&workspace.Workspace{AutoApply: true},
			&otf.ConfigurationVersion{},
		)

		got, err := f.NewRun(ctx, "", RunCreateOptions{})
		require.NoError(t, err)

		assert.True(t, got.AutoApply)
	})

	t.Run("run auto-apply", func(t *testing.T) {
		f := testFactory(
			&workspace.Workspace{},
			&otf.ConfigurationVersion{},
		)

		got, err := f.NewRun(ctx, "", RunCreateOptions{
			AutoApply: otf.Bool(true),
		})
		require.NoError(t, err)

		assert.True(t, got.AutoApply)
	})
}

func testFactory(ws *workspace.Workspace, cv *otf.ConfigurationVersion) *factory {
	return &factory{
		WorkspaceService:            &fakeFactoryWorkspaceService{ws: ws},
		ConfigurationVersionService: &fakeFactoryConfigurationVersionService{cv: cv},
	}
}

type fakeFactoryWorkspaceService struct {
	ws *workspace.Workspace
	workspace.WorkspaceService
}

func (f *fakeFactoryWorkspaceService) GetWorkspace(context.Context, string) (*workspace.Workspace, error) {
	return f.ws, nil
}

type fakeFactoryConfigurationVersionService struct {
	cv *otf.ConfigurationVersion
	otf.ConfigurationVersionService
}

func (f *fakeFactoryConfigurationVersionService) GetConfigurationVersion(context.Context, string) (*otf.ConfigurationVersion, error) {
	return f.cv, nil
}

func (f *fakeFactoryConfigurationVersionService) GetLatestConfigurationVersion(context.Context, string) (*otf.ConfigurationVersion, error) {
	return f.cv, nil
}