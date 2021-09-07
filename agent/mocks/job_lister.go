package mocks

import "github.com/leg100/ots"

type JobLister struct {
	Jobs []*ots.Job
}

func (l *JobLister) ListJobs(opts ots.JobListOptions) ([]*ots.Job, error) {
	return l.Jobs, nil
}
