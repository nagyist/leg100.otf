package jobs

import (
	"context"
	"fmt"
	"sync"

	"github.com/leg100/otf/internal/agentservice"
	"github.com/leg100/otf/internal/pubsub"
)

type (
	// Assigner assigns jobs to agents
	Assigner struct {
		JobService
		agentservice.AgentService
		pubsub.Subscriber

		// queue of unassigned jobs
		queue []*Job
		// store of active agents, keyed by ID
		active map[string]*agentservice.Agent
		mu     sync.Mutex
	}
)

func (s *Assigner) Start(ctx context.Context) error {
	// ensure data structures are allocated/emptied whenever assigner is
	// started/re-started
	s.queue = make([]*Job, 0)
	s.active = make(map[string]*agentservice.Agent)

	// Unsubscribe whenever exiting this routine.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// subscribe to job and agent events
	sub, err := s.Subscribe(ctx, "job-assigner-")
	if err != nil {
		return err
	}
	// retrieve existing agents
	existingAgents, err := s.ListAgents(ctx)
	if err != nil {
		return err
	}
	// retrieve existing unassigned jobs
	existingJobs, err := s.ListJobs(ctx, ListJobsOptions{Status: Unassigned})
	if err != nil {
		return err
	}
	// spool existing unassigned jobs in reverse order; ListJobs returns jobs newest first,
	// whereas we want oldest first.
	for i := len(existingJobs) - 1; i >= 0; i-- {
		if err := s.handleJob(ctx, existingJobs[i]); err != nil {
			return fmt.Errorf("spooling existing unassigned job: %w", err)
		}
	}
	// spool existing agents
	for _, agent := range existingAgents {
		if err := s.handleAgent(ctx, agent); err != nil {
			return fmt.Errorf("spooling existing agent: %w", err)
		}
	}
	// then relay events
	for event := range sub {
		switch payload := event.Payload.(type) {
		case *Job:
			if err := s.handleJob(ctx, payload); err != nil {
				return fmt.Errorf("relaying job event: %w", err)
			}
		case *agentservice.Agent:
			if err := s.handleAgent(ctx, payload); err != nil {
				return fmt.Errorf("relaying agent event: %w", err)
			}
		}
	}
	return nil
}

func (s *Assigner) handleAgent(ctx context.Context, agent *agentservice.Agent) error {
	if !agent.IsActive() {
		// remove inactive agent
		delete(s.active, agent.ID)
		return nil
	}
	// keep record of latest state of agent
	s.active[agent.ID] = agent

	if len(s.queue) == 0 {
		// no jobs to assign
		return nil
	}
	if agent.Status == agentservice.Idle {
		// pop job off queue and assign to idle agent
		return s.assignJob(ctx, agent)
	}
	return nil
}

func (s *Assigner) handleJob(ctx context.Context, job *Job) error {
	switch job.Status {
	case Running, Assigned:
		// skip in-progress jobs
		return nil
	case Unassigned:
		s.queue = append(s.queue, job)
	}
	if len(s.queue) == 0 {
		// no jobs to assign
		return nil
	}
	// find idle agent to assign job to
	for _, agent := range s.active {
		if agent.Status != agentservice.Idle {
			continue
		}
		// pop job off queue and assign to agent
		return s.assignJob(ctx, agent)
	}
	return nil
}

func (s *Assigner) assignJob(ctx context.Context, agent *agentservice.Agent) error {
	job := s.queue[0]
	s.queue = s.queue[1:]
	s.active[agent.ID].Status = agentservice.Busy
	return s.AssignJob(ctx, AssignJobOptions{
		RunID:   job.RunID,
		Phase:   job.Phase,
		AgentID: agent.ID,
	})
}
