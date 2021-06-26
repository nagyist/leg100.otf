package main

import (
	"context"
	"fmt"
	"os"

	"github.com/leg100/ots/agent"
	cmdutil "github.com/leg100/ots/cmd"
	"github.com/spf13/cobra"
)

const (
	DefaultAddress = "localhost:8080"
)

func main() {
	// Configure ^C to terminate program
	ctx, cancel := context.WithCancel(context.Background())
	cmdutil.CatchCtrlC(cancel)

	if err := Run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, args []string) error {
	var dataDir string

	a := &agent.Agent{}

	cmd := &cobra.Command{
		Use:           "ots-agent",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetArgs(args)

	cmd.Flags().StringVar(&a.ServerAddr, "address", DefaultAddress, "Address of OTS server")
	cmd.Flags().StringVar(&a.ID, "id", agent.DefaultID, "Agent identifier")
	cmd.Flags().StringVar(&dataDir, "data-dir", agent.DefaultDataDir, "Path to directory to store agent related data")

	cmdutil.SetFlagsFromEnvVariables(cmd.Flags())

	if err := cmd.ParseFlags(os.Args[1:]); err != nil {
		return err
	}

	// Validate agent struct's fields
	if err := a.Validate(); err != nil {
		return err
	}

	a.Poller(ctx)

	return nil
}
