package jobs

import "github.com/leg100/otf/internal"

type (
	Job struct {
		ID     string
		RunID  string
		Phase  internal.PhaseType
		Status Status
	}
)
