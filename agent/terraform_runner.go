package agent

import (
	"bytes"
	"context"
	"os/exec"
)

type TerraformRunner interface {
	Plan(context.Context) ([]byte, error)
}

type runner struct {
	path string
}

func (r *runner) Plan(ctx context.Context) ([]byte, error) {
	initOut, err := r.run(ctx, "init")
	if err != nil {
		return nil, err
	}

	planOut, err := r.run(ctx, "plan")
	if err != nil {
		return nil, err
	}

	return append(initOut, planOut...), nil
}

func (r *runner) run(ctx context.Context, command string) ([]byte, error) {
	buf := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "terraform", command)
	cmd.Dir = r.path
	cmd.Stdout = buf
	cmd.Stderr = buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
