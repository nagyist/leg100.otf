package agent

import (
	"context"

	"github.com/leg100/ots"
)

type mockSpooler struct {
	queue chan *ots.Job
}

func newMockSpooler(run ...*ots.Job) *mockSpooler {
	queue := make(chan *ots.Job, len(run))
	for _, r := range run {
		queue <- r
	}
	return &mockSpooler{queue: queue}
}

func (s *mockSpooler) GetJob() <-chan *ots.Job {
	return s.queue
}

func (s *mockSpooler) Start(context.Context) {}
