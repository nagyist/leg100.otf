package agent

import (
	"github.com/go-logr/logr"
	"github.com/leg100/ots"
)

func mockNewPlanRunnerFn(run *ots.Run,
	cvs ots.ConfigurationVersionService,
	svs ots.StateVersionService,
	rs ots.RunService,
	log logr.Logger) *ots.MultiStep {

	return ots.NewMultiStep(
		[]ots.Step{},
	)
}

func mockNewPlanRunnerFnWithError(run *ots.Run,
	cvs ots.ConfigurationVersionService,
	svs ots.StateVersionService,
	rs ots.RunService,
	log logr.Logger) *ots.MultiStep {

	return ots.NewMultiStep(
		[]ots.Step{
			ots.NewCommandStep("/bin/false"),
		},
	)
}
