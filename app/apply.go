package app

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/leg100/otf"
	"github.com/leg100/otf/sql"
)

var _ otf.ApplyService = (*ApplyService)(nil)

type ApplyService struct {
	db *sql.DB
	otf.ChunkService
	otf.EventService
	logr.Logger
}

func NewApplyService(db *sql.DB, logService otf.ChunkService, logger logr.Logger, es otf.EventService) *ApplyService {
	return &ApplyService{
		db:           db,
		EventService: es,
		ChunkService: logService,
		Logger:       logger,
	}
}

func (s ApplyService) Get(ctx context.Context, id string) (*otf.Apply, error) {
	run, err := s.db.GetRun(ctx, otf.RunGetOptions{ApplyID: &id})
	if err != nil {
		return nil, err
	}
	return run.Apply, nil
}

// Claim claims an apply job on behalf of an agent.
func (s ApplyService) Claim(ctx context.Context, jobID string, opts otf.JobClaimOptions) (otf.Job, error) {
	run, err := s.db.UpdateStatus(ctx, otf.RunGetOptions{JobID: &jobID}, func(run *otf.Run) error {
		return run.Apply.Start()
	})
	if err != nil {
		s.Error(err, "starting apply", "id", jobID)
		return nil, err
	}
	s.V(0).Info("started apply", "run_id", run.ID(), "id", jobID)
	return run, nil
}

// Finish an apply job.
func (s ApplyService) Finish(ctx context.Context, jobID string, opts otf.JobFinishOptions) (otf.Job, error) {
	var run *otf.Run
	var report otf.ResourceReport
	err := s.db.Tx(ctx, func(db otf.DB) error {
		chunk, err := db.GetChunk(ctx, jobID, otf.GetChunkOptions{})
		if err != nil {
			return err
		}
		report, err = otf.ParseApplyOutput(string(chunk.Data))
		if err != nil {
			return fmt.Errorf("compiling report of applied changes: %w", err)
		}
		if err := db.CreateApplyReport(ctx, jobID, report); err != nil {
			return fmt.Errorf("saving applied changes report: %w", err)
		}
		run, err = db.UpdateStatus(ctx, otf.RunGetOptions{ApplyID: &jobID}, func(run *otf.Run) (err error) {
			return run.Apply.Finish()
		})
		return err
	})
	if err != nil {
		s.Error(err, "finishing apply", "id", jobID)
		return nil, err
	}
	s.V(0).Info("finished apply",
		"id", jobID,
		"adds", report.Additions,
		"changes", report.Changes,
		"destructions", report.Destructions)

	s.Publish(otf.Event{Payload: run, Type: otf.EventRunApplied})
	return run, nil
}
