package agent

import (
	"context"
	"log"
	"time"

	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
	"github.com/mitchellh/go-homedir"
)

const (
	DefaultDataDir = "~/.ots-agent"
	DefaultID      = "agent-001"
)

// Agent processes jobs
type Agent struct {
	// ID uniquely identifies the agent.
	ID string

	// DataDir stores artefacts relating to runs, i.e. downloaded plugins,
	// modules (?), configuration versions, state, etc.
	DataDir string

	// ServerAddr is the address (<host>:<port>) of the OTS server to connect
	// to.
	ServerAddr string

	ConfigurationService ots.ConfigurationVersionService
	StateVersionService  ots.StateVersionService
	PlanService          ots.PlanService
	RunService           ots.RunService
	TerraformRunner      TerraformRunner
}

// Validate ensures the correctness of the Agent's public fields. Necessary
// because there is no constructor.
func (a *Agent) Validate() error {
	dataDir, err := homedir.Expand(a.DataDir)
	if err != nil {
		return err
	}
	a.DataDir = dataDir

	return nil
}

// Poller polls the daemon for queued runs
func (a *Agent) NewJob(run *ots.Run) (*Job, error) {
	//os.MkdirTemp(filepath.Ea.DataDir
	return &Job{
		Agent: a,
		Run:   run,
	}, nil
}

// Poller polls the daemon for queued runs and launches jobs accordingly.
func (a *Agent) Poller(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		case <-time.After(time.Second):
		}

		runs, err := a.RunService.GetQueuedRuns(ots.RunListOptions{})
		if err != nil {
			log.Printf("unable to poll daemon: %s", err.Error())
		}
		if len(runs.Items) == 0 {
			continue
		}
		run := runs.Items[0]
		job, _ := a.NewJob(run)

		if err := job.Process(ctx); err != nil {
			_, err := a.PlanService.UpdatePlanStatus(run.Plan.ID, tfe.PlanPending)
			if err != nil {
				log.Printf("unable to update status: %s", err.Error())
			}
		}
	}
}
