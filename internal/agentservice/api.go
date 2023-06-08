package agentservice

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	otfhttp "github.com/leg100/otf/internal/http"
)

type api struct {
	Service
}

func (a *api) addHandlers(r *mux.Router) {
	r.HandleFunc("/api/agent/register", a.register).Methods("POST")
	r.HandleFunc("/api/agent/status", a.updateStatus).Methods("PUT")
	r.HandleFunc("/api/agent/jobs", a.getJob).Methods("GET")
}

func (a *api) register(w http.ResponseWriter, r *http.Request) {
	ip, err := otfhttp.GetClientIP(r)
	if err != nil {
		http.Error(w, "retrieving client IP address: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var params struct {
		Name         *string
		Organization string
		Version      string
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := a.Service.Register(r.Context(), RegisterOptions{
		Name:         params.Name,
		IPAddress:    ip,
		External:     true,
		Organization: &params.Organization,
		Version:      params.Version,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(struct{ ID string }{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (a *api) updateStatus(w http.ResponseWriter, r *http.Request) {
	var params struct {
		ID     string
		Status Status
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := a.Service.UpdateStatus(r.Context(), params.ID, params.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *api) getJob(w http.ResponseWriter, r *http.Request) {
	var params struct {
		ID string
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ch, err := a.Service.GetJob(r.Context(), params.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(<-ch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
