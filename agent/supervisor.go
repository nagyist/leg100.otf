package agent

import (
	"bytes"
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/leg100/ots"
)

const (
	DefaultConcurrency = 5
)

// Supervisor supervises jobs
type Supervisor struct {
	// concurrency is the max number of concurrent jobs
	concurrency int

	logr.Logger

	ots.RunService
	ots.ConfigurationVersionService
	ots.StateVersionService

	Spooler
}

func NewSupervisor(spooler Spooler, cvs ots.ConfigurationVersionService, svs ots.StateVersionService, rs ots.RunService, logger logr.Logger, concurrency int) *Supervisor {
	return &Supervisor{
		Spooler:                     spooler,
		RunService:                  rs,
		StateVersionService:         svs,
		ConfigurationVersionService: cvs,
		Logger:                      logger,
		concurrency:                 concurrency,
	}
}

// Start starts the supervisor daemon and workers
func (s *Supervisor) Start(ctx context.Context) {
	s.startWorkers(ctx)

	<-ctx.Done()
}

func (s *Supervisor) startWorkers(ctx context.Context) {
	for i := 0; i < s.concurrency; i++ {
		go func() {
			for job := range s.GetJob() {
				s.handleJob(ctx, job)
			}
		}()
	}
}

func (s *Supervisor) handleJob(ctx context.Context, job *ots.Job) {
	path, err := os.MkdirTemp("", "ots-job")
	if err != nil {
		// TODO: update job status with error
		s.Error(err, "unable to create temp path")
		return
	}

	if err := s.StartJob(job.ID, ots.JobStartOptions{AgentID: DefaultID}); err != nil {
		s.Error(err, "unable to start job")
		return
	}

	s.Info("processing job", "run", job.ID, "status", job.Status, "dir", path)

	// For logs
	out := new(bytes.Buffer)

	jobStatus := ots.JobCompleted
	msteps := ots.NewMultiStep(job.Steps(s.RunService, s.ConfigurationVersionService, s.StateVersionService))

	if err := msteps.Run(ctx, path, out, job); err != nil {
		s.Error(err, "unable to run job")
		jobStatus = ots.JobErrored
	}

	if err := s.FinishJob(job.ID, ots.JobFinishOptions{Status: jobStatus}); err != nil {
		s.Error(err, "unable to finish job")
	}

	if err := s.UploadJobLogs(job.ID, out.Bytes()); err != nil {
		s.Error(err, "unable to upload logs for job")
	}
}
