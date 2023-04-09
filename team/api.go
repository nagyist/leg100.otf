package team

import (
	"net/http"

	"github.com/gorilla/mux"
	otfhttp "github.com/leg100/otf/http"
	"github.com/leg100/otf/http/decode"
	"github.com/leg100/otf/http/jsonapi"
)

// api provides handlers for json:api endpoints
type api struct {
	svc AuthService
}

func (h *api) addHandlers(r *mux.Router) {
	r = otfhttp.APIRouter(r)

	// Team routes
	r.HandleFunc("/organizations/{organization_name}/teams", h.createTeam).Methods("POST")
	r.HandleFunc("/organizations/{organization_name}/teams/{team_name}", h.getTeam).Methods("GET")
	r.HandleFunc("/teams/{team_id}", h.deleteTeam).Methods("DELETE")
}

func (h *api) createTeam(w http.ResponseWriter, r *http.Request) {
	var params jsonapi.CreateTeamOptions
	if err := decode.Route(&params, r); err != nil {
		jsonapi.Error(w, err)
		return
	}
	if err := jsonapi.UnmarshalPayload(r.Body, &params); err != nil {
		jsonapi.Error(w, err)
		return
	}

	team, err := h.svc.CreateTeam(r.Context(), CreateTeamOptions{
		Name:         *params.Name,
		Organization: *params.Organization,
	})
	if err != nil {
		jsonapi.Error(w, err)
		return
	}

	jsonapi.WriteResponse(w, r,
		&jsonapi.Team{ID: team.ID, Name: team.Name},
		jsonapi.WithCode(http.StatusCreated))
}

func (h *api) getTeam(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Organization *string `schema:"organization_name,required"`
		Name         *string `schema:"team_name,required"`
	}
	if err := decode.All(&params, r); err != nil {
		jsonapi.Error(w, err)
		return
	}

	team, err := h.svc.GetTeam(r.Context(), *params.Organization, *params.Name)
	if err != nil {
		jsonapi.Error(w, err)
		return
	}

	jsonapi.WriteResponse(w, r, &jsonapi.Team{ID: team.ID, Name: team.Name})
}

func (h *api) deleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := decode.Param("team_id", r)
	if err != nil {
		jsonapi.Error(w, err)
		return
	}

	if err := h.svc.DeleteTeam(r.Context(), id); err != nil {
		jsonapi.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
