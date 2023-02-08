package agenttoken

import (
	"net/http"

	"github.com/leg100/otf"
	"github.com/leg100/otf/http/jsonapi"
)

type handlers struct {
	app appService
}

// Implements TFC state versions API:
//
// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/state-versions#state-versions-api
//
func (h *handlers) AddHandlers(r *mux.Router) {
	r.HandleFunc("/workspaces/{workspace_id}/state-versions", h.createVersion).Methods("POST")
	r.HandleFunc("/workspaces/{workspace_id}/current-state-version", h.getCurrentVersion).Methods("GET")
	r.HandleFunc("/state-versions/{id}", h.getVersion).Methods("GET")
	r.HandleFunc("/state-versions", h.listVersions).Methods("GET")
	r.HandleFunc("/state-versions/{id}/download", h.downloadState).Methods("GET")
}

func (s *Server) CreateAgentToken(w http.ResponseWriter, r *http.Request) {
	opts := jsonapiCreateOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, &opts); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	at, err := s.Application.CreateAgentToken(r.Context(), otf.CreateAgentTokenOptions{
		Description:  opts.Description,
		Organization: opts.Organization,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	jsonapi.WriteResponse(w, r, &AgentToken{at, true})
}

func (s *Server) GetCurrentAgent(w http.ResponseWriter, r *http.Request) {
	agent, err := otf.AgentFromContext(r.Context())
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	jsonapi.WriteResponse(w, r, &AgentToken{agent, false})
}

type AgentToken struct {
	*otf.AgentToken
	revealToken bool // toggle send auth token over wire
}

// ToJSONAPI assembles a JSON-API DTO.
func (t *AgentToken) ToJSONAPI() any {
	json := jsonapi.AgentToken{
		ID:           t.ID(),
		Organization: t.Organization(),
	}
	if t.revealToken {
		json.Token = t.Token()
	}
	return &json
}