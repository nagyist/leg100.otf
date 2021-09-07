package agent

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/leg100/ots"
	"github.com/leg100/ots/agent/mocks"
	"github.com/leg100/ots/mock"
	"github.com/stretchr/testify/assert"
)

// TestSupervisor_Start tests starting up the daemon and tests it handling a
// single job
func TestSupervisor_Start(t *testing.T) {
	t.SkipNow()
	want := &ots.Job{ID: "job-123", Status: ots.JobPending}

	// Capture the run ID that is passed to the job processor
	got := make(chan string)

	supervisor := &Supervisor{
		Logger: logr.Discard(),
		RunService: &mock.RunService{
			UploadJobLogsFn: func(id string, _ []byte) error {
				got <- id
				return nil
			},
		},
		Spooler:     mocks.NewSpooler(want),
		concurrency: 1,
	}

	go supervisor.Start(context.Background())

	assert.Equal(t, "run-123", <-got)
}

// TestSupervisor_StartError tests starting up the agent daemon and tests it handling
// it a single job that errors
func TestSupervisor_StartError(t *testing.T) {
	t.SkipNow()
	// Mock run service and capture the run status it receives
	got := make(chan ots.JobStatus)
	runService := &mock.RunService{
		UploadJobLogsFn: func(id string, _ []byte) error { return nil },
		StartJobFn: func(id string, opts ots.JobStartOptions) error {
			return nil
		},
		FinishJobFn: func(id string, opts ots.JobFinishOptions) error {
			got <- opts.Status
			return nil
		},
	}

	supervisor := &Supervisor{
		Logger:     logr.Discard(),
		RunService: runService,
		Spooler: mocks.NewSpooler(&ots.Job{
			ID:     "job-123",
			Status: ots.JobPending,
		}),
		concurrency: 1,
	}

	go supervisor.Start(context.Background())

	// assert agent correctly propagates an errored status update
	assert.Equal(t, ots.JobErrored, <-got)
}
