package otf

type planningState struct {
	run *Run
	*runStateMixin
}

func newPlanningState(r *Run) *planningState {
	return &planningState{
		run: r,
		runStateMixin: &runStateMixin{
			run: r,
		},
	}
}

func (s *planningState) String() string { return "planning" }

// TODO: compile plan report
func (s *planningState) Finish(svc RunService) (*ResourceReport, error) {
	report, err := CompilePlanReport(svc.GetPlanFile())
	if err != nil {
		s.Error(err, "compiling planned changes report", "id", planID)
		return err
	}

	if !s.run.HasChanges() || s.run.Speculative() {
		s.run.setState(s.run.plannedAndFinishedState)
	} else if s.run.autoApply {
		s.run.setState(s.run.applyQueuedState)
	} else {
		s.run.setState(s.run.plannedState)
	}
	return nil
}

func (s *planningState) Cancel() error {
	s.run.setState(s.run.canceledState)
	return nil
}
