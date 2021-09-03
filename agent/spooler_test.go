package agent

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/leg100/ots"
	"github.com/leg100/ots/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockJobLister struct {
	jobs []*ots.Job
}

func (l *mockJobLister) ListJobs(opts ots.JobListOptions) ([]*ots.Job, error) {
	return l.jobs, nil
}

type mockSubscription struct {
	c chan ots.Event
}

func (s *mockSubscription) C() <-chan ots.Event { return s.c }

func (s *mockSubscription) Close() error { return nil }

// TestSpooler_New tests the spooler constructor
func TestSpooler_New(t *testing.T) {
	want := &ots.Job{ID: "run-123", Status: ots.JobPending}

	spooler, err := NewSpooler(
		&mockJobLister{jobs: []*ots.Job{want}},
		&mock.EventService{},
		logr.Discard(),
	)
	require.NoError(t, err)

	assert.Equal(t, want, <-spooler.queue)
}

// TestSpooler_Start tests the spooler daemon start op
func TestSpooler_Start(t *testing.T) {
	spooler := &SpoolerDaemon{
		EventService: &mock.EventService{
			SubscribeFn: func(id string) ots.Subscription {
				return &mockSubscription{}
			},
		},
		Logger: logr.Discard(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		spooler.Start(ctx)
		done <- struct{}{}
	}()

	cancel()

	<-done
}

// TestSpooler_GetJob tests retrieving a job from the spooler
func TestSpooler_GetJob(t *testing.T) {
	want := &ots.Job{ID: "job-123", Status: ots.JobPending}

	spooler := &SpoolerDaemon{queue: make(chan *ots.Job, 1)}
	spooler.queue <- want

	assert.Equal(t, want, <-spooler.GetJob())
}

// TestSpooler_GetJobFromEvent tests retrieving a job from the spooler after an
// event is received
func TestSpooler_GetJobFromEvent(t *testing.T) {
	want := &ots.Job{ID: "job-123", Status: ots.JobPending}

	sub := mockSubscription{c: make(chan ots.Event, 1)}

	spooler := &SpoolerDaemon{
		queue: make(chan *ots.Job, 1),
		EventService: &mock.EventService{
			SubscribeFn: func(id string) ots.Subscription {
				return &sub
			},
		},
		Logger: logr.Discard(),
	}

	go spooler.Start(context.Background())

	// send event
	sub.c <- ots.Event{Type: ots.JobCreatedEvent, Payload: want}

	assert.Equal(t, want, <-spooler.GetJob())
}
