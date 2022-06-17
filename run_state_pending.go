package otf

type pendingState struct {
	run *Run
	*runStateMixin
}

func newPendingState(r *Run) *pendingState {
	return &pendingState{
		run: r,
		runStateMixin: &runStateMixin{
			run: r,
		},
	}
}

func (s *pendingState) Start() error {
	s.run.setState(s.run.planningState)
	return nil
}

func (s *pendingState) Cancel() error {
	s.run.setState(s.run.canceledState)
	return nil
}
