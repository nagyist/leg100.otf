package ots

import (
	"context"
	"io"
)

var _ Step = (*MultiStep)(nil)

// MultiStep is a series of steps run sequentially.
type MultiStep struct {
	steps    []Step
	current  int
	canceled bool
}

func NewMultiStep(steps []Step) *MultiStep {
	return &MultiStep{
		steps: steps,
	}
}

func (r *MultiStep) Cancel(force bool) {
	r.canceled = true

	if len(r.steps) > 0 {
		r.steps[r.current].Cancel(force)
	}
}

func (r *MultiStep) Run(ctx context.Context, path string, out io.Writer, job *Job) error {
	for i, s := range r.steps {
		if r.canceled {
			return nil
		}

		r.current = i

		if err := s.Run(ctx, path, out, job); err != nil {
			return err
		}
	}

	return nil
}
