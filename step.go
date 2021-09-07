package ots

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type StepsProvider interface {
	Steps(RunService, ConfigurationVersionService, StateVersionService) []Step
	String() string
}

// Step is a cancelable task that forms part of a larger task (see MultiStep).
type Step interface {
	// Cancel cancels the task started by run. Force toggles whether the task is
	// canceled gracefully or not.
	Cancel(force bool)
	// Run invokes the task. Path can be used to share artefacts with other
	// steps. Informational output is expected to be written to out.
	Run(ctx context.Context, path string, out io.Writer, job *Job) error
}

// CommandStep is a cancelable executable CLI task.
type CommandStep struct {
	cmd  string
	args []string
	proc *os.Process
}

func NewCommandStep(cmd string, args ...string) *CommandStep {
	return &CommandStep{
		cmd:  cmd,
		args: args,
	}
}

func (s *CommandStep) Cancel(force bool) {
	if s.proc == nil {
		return
	}

	if force {
		s.proc.Signal(os.Kill)
	} else {
		s.proc.Signal(os.Interrupt)
	}
}

func (s *CommandStep) Run(ctx context.Context, path string, out io.Writer, _ *Job) error {
	cmd := exec.Command(s.cmd, s.args...)
	cmd.Dir = path
	cmd.Stdout = out
	cmd.Stderr = out

	s.proc = cmd.Process

	return cmd.Run()
}

// FuncStep is a cancelable go func task
type FuncStep struct {
	cancel context.CancelFunc
	fn     func(context.Context, string, *Job) error
}

func NewFuncStep(fn func(context.Context, string, *Job) error) *FuncStep {
	return &FuncStep{
		fn: fn,
	}
}

func (s *FuncStep) Cancel(force bool) {
	if !force {
		return
	}
	if s.cancel == nil {
		return
	}
	s.cancel()
}

// Run invokes the func, setting the working dir to the given path
func (s *FuncStep) Run(ctx context.Context, path string, out io.Writer, job *Job) error {
	ctx, s.cancel = context.WithCancel(ctx)
	return s.fn(ctx, path, job)
}
