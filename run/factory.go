package run

import (
	"context"

	"github.com/leg100/otf"
	"github.com/leg100/otf/workspace"
)

// factory constructs runs
type factory struct {
	otf.ConfigurationVersionService
	workspace.WorkspaceService
}

// NewRun constructs a new run at the beginning of its lifecycle using the
// provided options.
func (f *factory) NewRun(ctx context.Context, workspaceID string, opts RunCreateOptions) (*Run, error) {
	ws, err := f.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	var cv *otf.ConfigurationVersion
	if opts.ConfigurationVersionID != nil {
		cv, err = f.GetConfigurationVersion(ctx, *opts.ConfigurationVersionID)
	} else {
		cv, err = f.GetLatestConfigurationVersion(ctx, workspaceID)
	}
	if err != nil {
		return nil, err
	}

	return NewRun(cv, ws, opts), nil
}