package ots

import tfe "github.com/leg100/go-tfe"

type RunState struct {
	tfe.RunStatusTimestamps
	Status tfe.RunStatus
}

type PlanState struct {
	tfe.PlanStatusTimestamps
	tfe.PlanStatus
}

func (rs *RunState) ChangeStatus(state tfe.RunStatus) {
	rs.Status = state
	rs.setTimestamp(state)
}

func (rs *RunState) setTimestamp(state tfe.RunStatus) {
	switch state {
	case tfe.RunPending:
		rs.PlanQueueableAt = TimeNow()
	case tfe.RunPlanQueued:
		rs.PlanQueuedAt = TimeNow()
	case tfe.RunPlanning:
		rs.PlanningAt = TimeNow()
	case tfe.RunPlanned:
		rs.PlannedAt = TimeNow()
	case tfe.RunPlannedAndFinished:
		rs.PlannedAndFinishedAt = TimeNow()
	case tfe.RunApplyQueued:
		rs.ApplyQueuedAt = TimeNow()
	case tfe.RunApplying:
		rs.ApplyingAt = TimeNow()
	case tfe.RunApplied:
		rs.AppliedAt = TimeNow()
	case tfe.RunErrored:
		rs.ErroredAt = TimeNow()
	case tfe.RunCanceled:
		rs.CanceledAt = TimeNow()
	case tfe.RunDiscarded:
		rs.DiscardedAt = TimeNow()
	}
}
