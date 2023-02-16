package run

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leg100/otf"
	"github.com/leg100/otf/http/decode"
	"github.com/leg100/otf/http/html"
	"github.com/leg100/otf/http/html/paths"
	"github.com/leg100/otf/logs"
)

type web struct {
	otf.Renderer

	app app
}

type htmlLogChunk struct {
	otf.Chunk
}

func (h *web) addHandlers(r *mux.Router) {
	r.HandleFunc("/workspaces/{workspace_id}/runs", h.listRuns)
	r.HandleFunc("/runs/{run_id}", h.getRun)
	r.HandleFunc("/runs/{run_id}/delete", h.deleteRun)
	r.HandleFunc("/runs/{run_id}/cancel", h.cancelRun)
	r.HandleFunc("/runs/{run_id}/apply", h.applyRun)
	r.HandleFunc("/runs/{run_id}/discard", h.discardRun)

	// this handles the link the terraform CLI shows during a plan/apply.
	r.HandleFunc("/app/{organization_name}/{workspace_id}/runs/{run_id}", h.getRun)
}

func (h *web) listRuns(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		WorkspaceID string `schema:"workspace_id,required"`
		otf.ListOptions
	}
	var params parameters
	if err := decode.All(&params, r); err != nil {
		html.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	ws, err := h.app.GetWorkspace(r.Context(), params.WorkspaceID)
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	runs, err := h.ListRuns(r.Context(), otf.RunListOptions{
		ListOptions: params.ListOptions,
		WorkspaceID: &params.WorkspaceID,
	})
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Render("run_list.tmpl", w, r, struct {
		*otf.RunList
		otf.Workspace
		StreamID string
	}{
		RunList:   runs,
		Workspace: ws,
		StreamID:  "watch-ws-runs-" + otf.GenerateRandomString(5),
	})
}

func (h *web) getRun(w http.ResponseWriter, r *http.Request) {
	runID, err := decode.Param("run_id", r)
	if err != nil {
		html.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	run, err := h.app.get(r.Context(), runID)
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ws, err := h.GetWorkspace(r.Context(), run.WorkspaceID())
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get existing logs thus far received for each phase. If none are found then don't treat
	// that as an error because it merely means no logs have yet been received.
	planLogs, err := h.GetChunk(r.Context(), logs.GetChunkOptions{
		RunID: run.ID(),
		Phase: otf.PlanPhase,
	})
	if err != nil && !errors.Is(err, otf.ErrResourceNotFound) {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	applyLogs, err := h.GetChunk(r.Context(), otf.GetChunkOptions{
		RunID: run.ID(),
		Phase: otf.ApplyPhase,
	})
	if err != nil && !errors.Is(err, otf.ErrResourceNotFound) {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Render("run_get.tmpl", w, r, struct {
		*Run
		Workspace *otf.Workspace
		PlanLogs  *htmlLogChunk
		ApplyLogs *htmlLogChunk
	}{
		Run:       run,
		Workspace: ws,
		PlanLogs:  &htmlLogChunk{planLogs},
		ApplyLogs: &htmlLogChunk{applyLogs},
	})
}

func (h *web) deleteRun(w http.ResponseWriter, r *http.Request) {
	runID, err := decode.Param("run_id", r)
	if err != nil {
		html.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	run, err := h.app.get(r.Context(), runID)
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.DeleteRun(r.Context(), runID)
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, paths.Workspace(run.WorkspaceID()), http.StatusFound)
}

func (app *web) cancelRun(w http.ResponseWriter, r *http.Request) {
	runID, err := decode.Param("run_id", r)
	if err != nil {
		html.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	run, err := app.GetRun(r.Context(), runID)
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = app.CancelRun(r.Context(), runID, otf.RunCancelOptions{})
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, paths.Runs(run.WorkspaceID()), http.StatusFound)
}

func (app *web) applyRun(w http.ResponseWriter, r *http.Request) {
	run, err := app.GetRun(r.Context(), mux.Vars(r)["run_id"])
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = app.ApplyRun(r.Context(), mux.Vars(r)["run_id"], otf.RunApplyOptions{})
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, paths.Run(run.ID())+"#apply", http.StatusFound)
}

func (app *web) discardRun(w http.ResponseWriter, r *http.Request) {
	run, err := app.GetRun(r.Context(), mux.Vars(r)["run_id"])
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = app.DiscardRun(r.Context(), mux.Vars(r)["run_id"], otf.RunDiscardOptions{})
	if err != nil {
		html.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, paths.Run(run.ID()), http.StatusFound)
}
