package agent

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/leg100/ots"
)

var _ Spooler = (*SpoolerDaemon)(nil)

// Spooler is a daemon that queues incoming run jobs
type Spooler interface {
	GetJob() <-chan *ots.Job
	Start(context.Context)
}

// SpoolerDaemon queues jobs.
type SpoolerDaemon struct {
	// Queue of queued runs
	queue chan *ots.Job
	// EventService allows subscribing to stream of events
	ots.EventService
	// Logger for logging various events
	logr.Logger
}

type JobLister interface {
	ListJobs(ots.JobListOptions) ([]*ots.Job, error)
}

const (
	// SpoolerCapacity is the max number of queued runs the spooler can store
	SpoolerCapacity = 100
)

// NewSpooler is a constructor for a Spooler pre-populated with queued runs
func NewSpooler(
	jl JobLister,
	es ots.EventService,
	logger logr.Logger) (*SpoolerDaemon, error) {

	// TODO: order jobs by created_at date
	jobs, err := jl.ListJobs(ots.JobListOptions{})
	if err != nil {
		return nil, err
	}

	// Populate queue
	queue := make(chan *ots.Job, SpoolerCapacity)
	for _, j := range jobs {
		queue <- j
	}

	return &SpoolerDaemon{
		queue:        queue,
		EventService: es,
		Logger:       logger,
	}, nil
}

// Start starts the spooler
func (s *SpoolerDaemon) Start(ctx context.Context) {
	sub := s.Subscribe(DefaultID)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-sub.C():
			s.handleEvent(event)
		}
	}
}

// GetJob retrieves receive-only job queue
func (s *SpoolerDaemon) GetJob() <-chan *ots.Job {
	return s.queue
}

func (s *SpoolerDaemon) handleEvent(ev ots.Event) {
	switch obj := ev.Payload.(type) {
	case *ots.Job:
		s.Info("job event received", "run", obj.ID, "type", ev.Type, "status", obj.Status)

		switch ev.Type {
		case ots.JobCreatedEvent:
			s.queue <- obj
		case ots.JobCanceledEvent:
			// TODO: forward event immediately to job supervisor
		}
	}
}
