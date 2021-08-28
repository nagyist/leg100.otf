package agent

import (
	"github.com/go-logr/logr"
	"github.com/leg100/go-tfe"
	"github.com/leg100/ots"
)

type NewApplyRunnerFn func(
	*ots.Run,
	ots.ConfigurationVersionService,
	ots.StateVersionService,
	ots.RunService,
	logr.Logger) *ots.MultiStep

func NewApplyRunner(run *ots.Run,
	cvs ots.ConfigurationVersionService,
	svs ots.StateVersionService,
	rs ots.RunService,
	log logr.Logger) *ots.MultiStep {

	return ots.NewMultiStep(
		[]ots.Step{
			NewDownloadConfigStep(run, cvs),
			NewDeleteBackendStep,
			DownloadPlanFileStep(run, rs),
			NewDownloadStateStep(run, svs, log),
			UpdateApplyStatusStep(run, rs, tfe.ApplyRunning),
			NewInitStep,
			ApplyStep,
			UploadStateStep(run, svs),
			FinishApplyStep(run, rs, log),
		},
	)
}
