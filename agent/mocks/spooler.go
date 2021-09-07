package mocks

import (
	"context"

	"github.com/leg100/ots"
)

type Spooler struct {
	queue chan *ots.Job
}

func NewSpooler(run ...*ots.Job) *Spooler {
	queue := make(chan *ots.Job, len(run))
	for _, r := range run {
		queue <- r
	}
	return &Spooler{queue: queue}
}

func (s *Spooler) GetJob() <-chan *ots.Job {
	return s.queue
}

func (s *Spooler) Start(context.Context) {}
