package agent

import (
	"github.com/go-logr/logr"
	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
)

type NewPlanRunnerFn func(
	*ots.Run,
	ots.ConfigurationVersionService,
	ots.StateVersionService,
	ots.RunService,
	logr.Logger) *ots.MultiStep

func NewPlanRunner(run *ots.Run,
	cvs ots.ConfigurationVersionService,
	svs ots.StateVersionService,
	rs ots.RunService,
	log logr.Logger) *ots.MultiStep {

	return ots.NewMultiStep(
		[]ots.Step{
			NewDownloadConfigStep(run, cvs),
			NewDeleteBackendStep,
			NewDownloadStateStep(run, svs, log),
			NewUpdatePlanStatusStep(run, rs, tfe.PlanRunning),
			NewInitStep,
			NewPlanStep,
			NewJSONPlanStep,
			NewFinishPlanStep(run, rs, log),
		},
	)
}
