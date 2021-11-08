package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	cmdutil "github.com/leg100/otf/cmd"
)

func main() {
	handleExitCode(run())
}

func run() error {
	// Configure ^C to terminate program
	ctx, cancel := context.WithCancel(context.Background())
	cmdutil.CatchCtrlC(cancel)

	cmd := exec.CommandContext(ctx, "htop")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func handleExitCode(err error) {
	if err == nil {
		os.Exit(0)
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}

	fmt.Println("error: ", err.Error())
	os.Exit(1)
}
