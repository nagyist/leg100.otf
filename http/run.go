package http

import (
	"net/http"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-tfe"
	"github.com/leg100/ots"
)

func (h *Server) CreateRun(w http.ResponseWriter, r *http.Request) {
	CreateObject(w, r, &tfe.RunCreateOptions{}, func(opts interface{}) (interface{}, error) {
		return h.RunService.CreateRun(opts.(*tfe.RunCreateOptions))
	})
}

func (h *Server) ApplyRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunApplyOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		ErrUnprocessable(w, err)
		return
	}

	if err := h.RunService.ApplyRun(vars["id"], opts); err != nil {
		ErrNotFound(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-type", jsonapi.MediaType)
}

func (h *Server) GetRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	GetObject(w, r, func() (interface{}, error) {
		return h.RunService.GetRun(vars["id"])
	})
}

func (h *Server) ListRuns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var opts ots.ListOptions
	if err := DecodeAndSanitize(&opts, r.URL.Query()); err != nil {
		ErrUnprocessable(w, err)
		return
	}

	ListObjects(w, r, func() (interface{}, error) {
		return h.RunService.ListRuns(vars["workspace_id"], opts)
	})
}

func (h *Server) DiscardRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunDiscardOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		ErrUnprocessable(w, err)
		return
	}

	if err := h.RunService.DiscardRun(vars["id"], opts); err != nil {
		ErrNotFound(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-type", jsonapi.MediaType)
}

func (h *Server) CancelRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunCancelOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		ErrUnprocessable(w, err)
		return
	}

	if err := h.RunService.CancelRun(vars["id"], opts); err != nil {
		ErrNotFound(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-type", jsonapi.MediaType)
}

func (h *Server) ForceCancelRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	opts := &tfe.RunForceCancelOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		ErrUnprocessable(w, err)
		return
	}

	if err := h.RunService.ForceCancelRun(vars["id"], opts); err != nil {
		ErrNotFound(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-type", jsonapi.MediaType)
}
