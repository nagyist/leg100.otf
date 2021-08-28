package agent

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/leg100/go-tfe"
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

	cvs ots.ConfigurationVersionService
	svs ots.StateVersionService
	rs  ots.RunService
}

type RunLister interface {
	List(ots.RunListOptions) (*ots.RunList, error)
}

const (
	// SpoolerCapacity is the max number of queued runs the spooler can store
	SpoolerCapacity = 100
)

var (
	// QueuedStatuses are the list of run statuses that indicate it is in a
	// queued state
	QueuedStatuses = []tfe.RunStatus{tfe.RunPlanQueued, tfe.RunApplyQueued}
)

// NewSpooler is a constructor for a Spooler pre-populated with queued runs
func NewSpooler(
	cvs ots.ConfigurationVersionService,
	svs ots.StateVersionService,
	rs ots.RunService,
	es ots.EventService,
	logger logr.Logger) (*SpoolerDaemon, error) {

	// TODO: order runs by created_at date
	runs, err := rs.List(ots.RunListOptions{Statuses: QueuedStatuses})
	if err != nil {
		return nil, err
	}

	// Populate queue
	queue := make(chan *ots.Job, SpoolerCapacity)
	for _, r := range runs.Items {
		queue <- ots.NewJobFromRun(r, cvs, svs, rs)
	}

	return &SpoolerDaemon{
		queue:        queue,
		EventService: es,
		cvs:          cvs,
		svs:          svs,
		rs:           rs,
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
	case *ots.Run:
		s.Info("run event received", "run", obj.ID, "type", ev.Type, "status", obj.Status)

		switch ev.Type {
		case ots.PlanQueued, ots.ApplyQueued:
			s.queue <- ots.NewJobFromRun(obj, s.cvs, s.svs, s.rs)
		case ots.RunCanceled:
			// TODO: forward event immediately to job supervisor
		}
	}
}
