package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leg100/otf"
	"github.com/leg100/otf/http/decode"
)

// PlanFileOptions represents the options for retrieving the plan file for a
// run.
type PlanFileOptions struct {
	// Format of plan file. Valid values are json and binary.
	Format otf.PlanFormat `schema:"format"`
}

func (s *Server) GetPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	plan, err := s.PlanService().Get(r.Context(), vars["id"])
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeResponse(w, r, plan)
}

func (s *Server) GetPlanJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	json, err := s.RunService().GetPlanFile(r.Context(), otf.RunGetOptions{PlanID: &id}, otf.PlanFormatJSON)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if _, err := w.Write(json); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetPlanLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var opts otf.GetChunkOptions
	if err := decode.Query(&opts, r.URL.Query()); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	chunk, err := s.PlanService().Get(r.Context(), vars["job_id"], opts)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if _, err := w.Write(chunk.Marshal()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
