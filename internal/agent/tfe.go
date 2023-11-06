package agent

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leg100/otf/internal/http/decode"
	"github.com/leg100/otf/internal/resource"
	"github.com/leg100/otf/internal/tfeapi"
	"github.com/leg100/otf/internal/tfeapi/types"
)

type tfe struct {
	*service
	*tfeapi.Responder
}

func (a *tfe) addHandlers(r *mux.Router) {
	r = r.PathPrefix(tfeapi.APIPrefixV2).Subrouter()

	// Agent Pools And Agents API
	r.HandleFunc("/organizations/{organization_name}/agent-pools", a.listAgentPools).Methods("GET")
	r.HandleFunc("/organizations/{organization_name}/agent-pools", a.createAgentPool).Methods("POST")
	r.HandleFunc("/agent-pools/{pool_id}", a.getAgentPool).Methods("GET")
	r.HandleFunc("/agent-pools/{pool_id}", a.updateAgentPool).Methods("PATCH")
	r.HandleFunc("/agent-pools/{pool_id}", a.deleteAgentPool).Methods("DELETE")

	// Feature sets API:
	//
	// https://developer.hashicorp.com/terraform/cloud-docs/api-docs/feature-sets
	//
	// This is only implemented in order to get the go-tfe integration tests
	// passing. Those tests first do some silly shit with the feature sets
	// API before hitting the agent pool API, so the former is stubbed out here
	// to make 'em happy.
	r.HandleFunc("/admin/feature-sets", func(w http.ResponseWriter, r *http.Request) {
		// tests expect one feature set to be returned but don't check contents
		// so return the bare minimum.
		fs := []struct {
			ID string `jsonapi:"primary,feature-sets"`
		}{
			{ID: "fs-123"},
		}
		a.RespondWithPage(w, r, &fs, &resource.Pagination{})
	})
	// Tests don't check response so return empty response.
	r.HandleFunc("/admin/organizations/{organization_name}/subscription", func(w http.ResponseWriter, r *http.Request) {})
}

func (a *tfe) createAgentPool(w http.ResponseWriter, r *http.Request) {
	organization, err := decode.Param("organization_name", r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	var params types.AgentPoolCreateOptions
	if err := tfeapi.Unmarshal(r.Body, &params); err != nil {
		tfeapi.Error(w, err)
		return
	}

	// convert tfe params to otf opts
	opts := createPoolOptions{
		Name:               params.Name,
		Organization:       organization,
		OrganizationScoped: params.OrganizationScoped,
	}
	opts.AllowedWorkspaces = make([]string, len(params.AllowedWorkspaces))
	for i, aw := range params.AllowedWorkspaces {
		opts.AllowedWorkspaces[i] = aw.ID
	}

	pool, err := a.service.createPool(r.Context(), opts)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	a.Respond(w, r, a.toPool(pool), http.StatusCreated)
}

func (a *tfe) updateAgentPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := decode.Param("pool_id", r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	var params types.AgentPoolUpdateOptions
	if err := tfeapi.Unmarshal(r.Body, &params); err != nil {
		tfeapi.Error(w, err)
		return
	}

	// convert tfe params to otf opts
	opts := updatePoolOptions{
		Name:               params.Name,
		OrganizationScoped: params.OrganizationScoped,
	}
	if params.AllowedWorkspaces != nil {
		opts.AllowedWorkspaces = make([]string, len(params.AllowedWorkspaces))
		for i, aw := range params.AllowedWorkspaces {
			opts.AllowedWorkspaces[i] = aw.ID
		}
	}

	pool, err := a.service.updatePool(r.Context(), poolID, opts)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	a.Respond(w, r, a.toPool(pool), http.StatusOK)
}

func (a *tfe) getAgentPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := decode.Param("pool_id", r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	pool, err := a.service.getPool(r.Context(), poolID)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	a.Respond(w, r, a.toPool(pool), http.StatusOK)
}

func (a *tfe) listAgentPools(w http.ResponseWriter, r *http.Request) {
	organization, err := decode.Param("organization_name", r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}
	var params types.AgentPoolListOptions
	if err := decode.All(&params, r); err != nil {
		tfeapi.Error(w, err)
		return
	}

	pools, err := a.service.listPools(r.Context(), listPoolOptions{
		Organization:         &organization,
		NameSubstring:        params.Query,
		AllowedWorkspaceName: params.AllowedWorkspacesName,
	})
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	// client expects a page, whereas listPools returns full result set, so
	// convert to page first
	page := resource.NewPage(pools, resource.PageOptions(params.ListOptions), nil)

	// convert items
	items := make([]*types.AgentPool, len(page.Items))
	for i, from := range page.Items {
		if err != nil {
			tfeapi.Error(w, err)
			return
		}
		items[i] = a.toPool(from)
	}
	a.RespondWithPage(w, r, items, page.Pagination)
}

func (a *tfe) deleteAgentPool(w http.ResponseWriter, r *http.Request) {
	poolID, err := decode.Param("pool_id", r)
	if err != nil {
		tfeapi.Error(w, err)
		return
	}

	if err := a.service.deletePool(r.Context(), poolID); err != nil {
		tfeapi.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *tfe) toPool(from *Pool) *types.AgentPool {
	to := &types.AgentPool{
		ID:   from.ID,
		Name: from.Name,
		Organization: &types.Organization{
			Name: from.Organization,
		},
		OrganizationScoped: from.OrganizationScoped,
	}
	to.Workspaces = make([]*types.Workspace, len(from.Workspaces))
	for i, workspaceID := range from.Workspaces {
		to.Workspaces[i] = &types.Workspace{ID: workspaceID}
	}
	to.AllowedWorkspaces = make([]*types.Workspace, len(from.AllowedWorkspaces))
	for i, workspaceID := range from.AllowedWorkspaces {
		to.AllowedWorkspaces[i] = &types.Workspace{ID: workspaceID}
	}
	return to
}