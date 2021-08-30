package ots

import tfe "github.com/leg100/go-tfe"

var _ Operation = (*PlanOperation)(nil)
var _ Operation = (*ApplyOperation)(nil)

type Operation interface {
	Start(job *Job)
	Finish(job *Job, opts JobFinishOptions)
}

type PlanOperation string

func NewPlanOperation() PlanOperation { return PlanOperation("plan") }

func (o PlanOperation) Start(job *Job) {
	job.Run.Status = tfe.RunPlanning
	job.Run.StatusTimestamps.PlanningAt = TimeNow()

	job.Run.Plan.Status = tfe.PlanRunning
	job.Run.Plan.StatusTimestamps.StartedAt = TimeNow()
}

func (o PlanOperation) Finish(job *Job, success bool) {
	if !success {
		job.Run.Status = tfe.RunErrored
		job.Run.StatusTimestamps.ErroredAt = TimeNow()

		job.Run.Plan.Status = tfe.PlanErrored
		job.Run.Plan.StatusTimestamps.ErroredAt = TimeNow()

		return
	}

	if job.ConfigurationVersion.Speculative {
		job.Run.Status = tfe.RunPlannedAndFinished
		job.Run.StatusTimestamps.PlannedAndFinishedAt = TimeNow()
	} else {
		job.Run.Status = tfe.RunPlanned
		job.Run.StatusTimestamps.PlannedAt = TimeNow()
	}

	job.Run.Plan.Status = tfe.PlanFinished
	job.Run.Plan.StatusTimestamps.FinishedAt = TimeNow()
}

type ApplyOperation string

func NewApplyOperation() ApplyOperation { return ApplyOperation("apply") }

func (o ApplyOperation) Start(job *Job) {
	job.Run.Status = tfe.RunApplying
	job.Run.StatusTimestamps.ApplyingAt = TimeNow()

	job.Run.Apply.Status = tfe.ApplyRunning
	job.Run.Apply.StatusTimestamps.StartedAt = TimeNow()
}

func (o ApplyOperation) Finish(job *Job, opts JobFinishOptions) {
	job.Run.Status = tfe.RunApplied
	job.Run.StatusTimestamps.AppliedAt = TimeNow()

	job.Run.Apply.Status = tfe.ApplyFinished
	job.Run.Apply.StatusTimestamps.FinishedAt = TimeNow()
}
