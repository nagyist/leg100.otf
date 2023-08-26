package variable

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/http/html"
	"github.com/leg100/otf/internal/organization"
	"github.com/leg100/otf/internal/rbac"
	"github.com/leg100/otf/internal/run"
	"github.com/leg100/otf/internal/sql"
	"github.com/leg100/otf/internal/sql/pggen"
	"github.com/leg100/otf/internal/tfeapi"
	"github.com/leg100/otf/internal/workspace"
)

type (
	VariableService = Service

	Service interface {
		// MergeVariables merges variables for a workspace according to the
		// precedence rules documented here:
		//
		// https://developer.hashicorp.com/terraform/cloud-docs/workspaces/variables#precedence
		MergeVariables(ctx context.Context, run *run.Run) ([]*Variable, error)

		ListVariables(ctx context.Context, workspaceID string) ([]*Variable, error)

		CreateWorkspaceVariable(ctx context.Context, workspaceID string, opts CreateVariableOptions) (*WorkspaceVariable, error)
		UpdateWorkspaceVariable(ctx context.Context, variableID string, opts UpdateVariableOptions) (*WorkspaceVariable, error)
		ListWorkspaceVariables(ctx context.Context, workspaceID string) ([]*WorkspaceVariable, error)
		GetWorkspaceVariable(ctx context.Context, variableID string) (*WorkspaceVariable, error)
		DeleteWorkspaceVariable(ctx context.Context, variableID string) (*WorkspaceVariable, error)

		createVariableSet(ctx context.Context, organization string, opts CreateVariableSetOptions) (*VariableSet, error)
		updateVariableSet(ctx context.Context, setID string, opts UpdateVariableSetOptions) (*VariableSet, error)
		listVariableSets(ctx context.Context, organization string) ([]*VariableSet, error)
		listWorkspaceVariableSets(ctx context.Context, workspaceID string) ([]*VariableSet, error)
		getVariableSet(ctx context.Context, setID string) (*VariableSet, error)
		deleteVariableSet(ctx context.Context, setID string) error

		addVariableToSet(ctx context.Context, setID string, opts CreateVariableOptions) error
		updateVariableSetVariable(ctx context.Context, variableID string, opts UpdateVariableOptions) (*Variable, error)
		deleteVariableFromSet(ctx context.Context, setID, variableID string) error
		applySetToWorkspaces(ctx context.Context, setID string, workspaceIDs []string) error
		deleteSetFromWorkspaces(ctx context.Context, setID string, workspaceIDs []string) error
	}

	service struct {
		logr.Logger

		db           *pgdb
		web          *web
		api          *tfe
		workspace    internal.Authorizer
		organization internal.Authorizer

		*factory
	}

	Options struct {
		WorkspaceAuthorizer internal.Authorizer
		WorkspaceService    workspace.Service

		*sql.DB
		*tfeapi.Responder
		html.Renderer
		logr.Logger
	}
)

func NewService(opts Options) *service {
	svc := service{
		Logger:       opts.Logger,
		db:           &pgdb{opts.DB},
		factory:      &factory{generateVersion: versionGenerator},
		workspace:    opts.WorkspaceAuthorizer,
		organization: &organization.Authorizer{Logger: opts.Logger},
	}

	svc.web = &web{
		Renderer: opts.Renderer,
		Service:  opts.WorkspaceService,
		svc:      &svc,
	}
	svc.api = &tfe{
		Service:   &svc,
		Responder: opts.Responder,
	}

	return &svc
}

func (s *service) AddHandlers(r *mux.Router) {
	s.web.addHandlers(r)
	s.api.addHandlers(r)
}

func (s *service) MergeVariables(ctx context.Context, run *run.Run) ([]*Variable, error) {
	sets, err := s.listVariableSets(ctx, run.Organization)
	if err != nil {
		return nil, err
	}
	workspaceVariables, err := s.ListWorkspaceVariables(ctx, run.WorkspaceID)
	if err != nil {
		return nil, err
	}
	return mergeVariables(run, workspaceVariables, sets), nil
}

func (s *service) ListVariables(ctx context.Context, workspaceID string) ([]*Variable, error) {
	workspaceVariables, err := s.ListWorkspaceVariables(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	// convert from []*WorkspaceVariable to []*Variable
	variables := make([]*Variable, len(workspaceVariables))
	for i, v := range workspaceVariables {
		variables[i] = v.Variable
	}

	return variables, nil
}

func (s *service) CreateWorkspaceVariable(ctx context.Context, workspaceID string, opts CreateVariableOptions) (*WorkspaceVariable, error) {
	subject, err := s.workspace.CanAccess(ctx, rbac.CreateWorkspaceVariableAction, workspaceID)
	if err != nil {
		return nil, err
	}

	v, err := s.newWorkspaceVariable(workspaceID, opts)
	if err != nil {
		s.Error(err, "constructing workspace variable", "subject", subject, "workspace", workspaceID, "key", opts.Key)
		return nil, err
	}

	err = s.db.Lock(ctx, "variables", func(ctx context.Context, q pggen.Querier) (err error) {
		// check for conflict with other workspace variables
		all, err := s.db.listWorkspaceVariables(ctx, workspaceID)
		if err != nil {
			return err
		}
		for _, v := range all {
			if err := v.conflicts(v.Variable); err != nil {
				return err
			}
		}
		if err := s.db.createWorkspaceVariable(ctx, v); err != nil {
			s.Error(err, "creating workspace variable", "subject", subject, "variable", v)
			return err
		}
		return nil
	})
	if err != nil {
		s.Error(err, "creating workspace variable", "subject", subject, "variable", v)
		return nil, err
	}

	s.V(1).Info("created workspace variable", "subject", subject, "variable", v)

	return v, nil
}

func (s *service) UpdateWorkspaceVariable(ctx context.Context, variableID string, opts UpdateVariableOptions) (*WorkspaceVariable, error) {
	var (
		subject  internal.Subject
		existing *WorkspaceVariable
		updated  Variable
	)
	err := s.db.Lock(ctx, "variables", func(ctx context.Context, q pggen.Querier) (err error) {
		existing, err = s.db.getWorkspaceVariable(ctx, variableID)
		if err != nil {
			return err
		}

		subject, err = s.workspace.CanAccess(ctx, rbac.UpdateWorkspaceVariableAction, existing.WorkspaceID)
		if err != nil {
			return err
		}

		updated = *existing.Variable
		if err := s.update(&updated, opts); err != nil {
			return err
		}

		// check for conflict with other workspace variables
		all, err := s.db.listWorkspaceVariables(ctx, existing.WorkspaceID)
		if err != nil {
			return err
		}
		for _, v := range all {
			if err := existing.conflicts(v.Variable); err != nil {
				return err
			}
		}
		if err := s.db.updateVariable(ctx, &updated); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.Error(err, "updating workspace variable", "subject", subject, "variable_id", variableID)
		return nil, err
	}
	s.V(1).Info("updated workspace variable", "subject", subject, "workspace_id", existing.WorkspaceID, "before", existing.Variable, "after", &updated)

	return &WorkspaceVariable{WorkspaceID: existing.WorkspaceID, Variable: &updated}, nil
}

func (s *service) ListWorkspaceVariables(ctx context.Context, workspaceID string) ([]*WorkspaceVariable, error) {
	subject, err := s.workspace.CanAccess(ctx, rbac.ListWorkspaceVariablesAction, workspaceID)
	if err != nil {
		return nil, err
	}

	variables, err := s.db.listWorkspaceVariables(ctx, workspaceID)
	if err != nil {
		s.Error(err, "listing workspace variables", "subject", subject, "workspace_id", workspaceID)
		return nil, err
	}

	s.V(9).Info("listed workspace variables", "subject", subject, "workspace_id", workspaceID, "count", len(variables))

	return variables, nil
}

func (s *service) GetWorkspaceVariable(ctx context.Context, variableID string) (*WorkspaceVariable, error) {
	// get workspace variable first for authorization purposes
	wv, err := s.db.getWorkspaceVariable(ctx, variableID)
	if err != nil {
		s.Error(err, "retrieving workspace variable", "variable_id", variableID)
		return nil, err
	}

	subject, err := s.workspace.CanAccess(ctx, rbac.GetWorkspaceVariableAction, wv.WorkspaceID)
	if err != nil {
		return nil, err
	}

	s.V(9).Info("retrieved variable", "subject", subject, "workspace_variable", wv)

	return wv, nil
}

func (s *service) DeleteWorkspaceVariable(ctx context.Context, variableID string) (*WorkspaceVariable, error) {
	// get workspace variable first for authorization purposes
	wv, err := s.db.getWorkspaceVariable(ctx, variableID)
	if err != nil {
		return nil, err
	}

	subject, err := s.workspace.CanAccess(ctx, rbac.DeleteWorkspaceVariableAction, wv.WorkspaceID)
	if err != nil {
		return nil, err
	}

	if err := s.db.deleteVariable(ctx, variableID); err != nil {
		s.Error(err, "deleting workspace variable", "subject", subject, "variable", wv)
		return nil, err
	}
	s.V(1).Info("deleted workspace variable", "subject", subject, "variable", wv)

	return wv, nil
}

func (s *service) createVariableSet(ctx context.Context, organization string, opts CreateVariableSetOptions) (*VariableSet, error) {
	subject, err := s.organization.CanAccess(ctx, rbac.CreateVariableSetAction, organization)
	if err != nil {
		return nil, err
	}

	set, err := s.newSet(organization, opts)
	if err != nil {
		s.Error(err, "constructing variable set", "subject", subject, "organization", organization)
		return nil, err
	}

	if err := s.db.createVariableSet(ctx, set); err != nil {
		s.Error(err, "creating variable set", "subject", subject, "set", set)
		return nil, err
	}

	s.V(1).Info("created variable set", "subject", subject, "set", set)

	return set, nil
}

func (s *service) updateVariableSet(ctx context.Context, setID string, opts UpdateVariableSetOptions) (*VariableSet, error) {
	var (
		subject internal.Subject
		before  VariableSet
		after   *VariableSet
	)
	err := s.db.Lock(ctx, "variables", func(ctx context.Context, q pggen.Querier) (err error) {
		after, err = s.db.updateVariableSet(ctx, setID, func(existing *VariableSet) error {
			subject, err = s.organization.CanAccess(ctx, rbac.UpdateVariableSetAction, existing.Organization)
			if err != nil {
				return err
			}

			before = *existing
			if err := existing.update(opts); err != nil {
				return err
			}

			// if set has been promoted to global then we need to check for
			// conflicts
			if !existing.Global && existing.Global {
				sets, err := s.db.listVariableSets(ctx, existing.Organization)
				if err != nil {
					return err
				}
				for _, v := range existing.Variables {
					if err := checkConflicts(v, existing, sets); err != nil {
						return err
					}
				}
			}
			return nil
		})
		return err
	})
	if err != nil {
		s.Error(err, "updating variable set", "subject", subject, "set_id", setID)
		return nil, err
	}
	s.V(1).Info("updated variable set", "subject", subject, "before", &before, "after", after)

	return after, nil
}

func (s *service) listVariableSets(ctx context.Context, organization string) ([]*VariableSet, error) {
	subject, err := s.organization.CanAccess(ctx, rbac.ListVariableSetsAction, organization)
	if err != nil {
		return nil, err
	}

	sets, err := s.db.listVariableSets(ctx, organization)
	if err != nil {
		s.Error(err, "listing variable sets", "subject", subject, "organization", organization)
		return nil, err
	}
	s.V(9).Info("listed variable sets", "subject", subject, "organization", organization, "count", len(sets))

	return sets, nil
}

func (s *service) listWorkspaceVariableSets(ctx context.Context, workspaceID string) ([]*VariableSet, error) {
	subject, err := s.workspace.CanAccess(ctx, rbac.ListVariableSetsAction, workspaceID)
	if err != nil {
		return nil, err
	}

	sets, err := s.db.listVariableSetsByWorkspace(ctx, workspaceID)
	if err != nil {
		s.Error(err, "listing variable sets", "subject", subject, "workspace_id", workspaceID)
		return nil, err
	}
	s.V(9).Info("listed variable sets", "subject", subject, "workspace_id", workspaceID, "count", len(sets))

	return sets, nil
}

func (s *service) getVariableSet(ctx context.Context, setID string) (*VariableSet, error) {
	// retrieve set first in order to retrieve organization name for authorization
	set, err := s.db.getVariableSet(ctx, setID)
	if err != nil {
		s.Error(err, "retrieving variable set", "set_id", setID)
		return nil, err
	}

	subject, err := s.organization.CanAccess(ctx, rbac.GetVariableSetAction, set.Organization)
	if err != nil {
		s.Error(err, "retrieving variable set", "subject", subject, "set", set)
		return nil, err
	}
	s.V(9).Info("retrieved variable set", "subject", subject, "set", set)

	return set, nil
}

func (s *service) deleteVariableSet(ctx context.Context, setID string) error {
	// retrieve existing set in order to retrieve organization for authorization
	existing, err := s.db.getVariableSet(ctx, setID)
	if err != nil {
		return err
	}

	subject, err := s.organization.CanAccess(ctx, rbac.DeleteVariableSetAction, existing.Organization)
	if err != nil {
		return err
	}

	if err := s.db.deleteVariableSet(ctx, setID); err != nil {
		s.Error(err, "deleting variable set", "subject", subject, "set", existing)
		return err
	}
	s.V(1).Info("deleted variable set", "subject", subject, "set", existing)

	return nil
}

func (s *service) addVariableToSet(ctx context.Context, setID string, opts CreateVariableOptions) error {
	var (
		subject internal.Subject
		set     *VariableSet
		v       *Variable
	)
	err := s.db.Lock(ctx, "variables", func(ctx context.Context, q pggen.Querier) (err error) {
		set, err = s.db.getVariableSet(ctx, setID)
		if err != nil {
			return err
		}

		subject, err := s.organization.CanAccess(ctx, rbac.AddVariableToSetAction, set.Organization)
		if err != nil {
			return err
		}

		v, err = s.new(opts)
		if err != nil {
			s.Error(err, "adding variable to set: constructing variable", "subject", subject)
			return err
		}

		sets, err := s.db.listVariableSets(ctx, set.Organization)
		if err != nil {
			return err
		}
		if err := checkConflicts(v, set, sets); err != nil {
			return err
		}

		if err := s.db.addVariableToSet(ctx, setID, v); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.Error(err, "adding variable to set", "subject", subject, "set", set, "variable", v)
		return err
	}

	s.V(1).Info("added variable to set", "subject", subject, "set", set, "variable", v)

	return nil
}

func (s *service) updateVariableSetVariable(ctx context.Context, variableID string, opts UpdateVariableOptions) (*Variable, error) {
	var (
		subject  internal.Subject
		set      *VariableSet
		existing *Variable
		updated  Variable
	)
	err := s.db.Lock(ctx, "variables", func(ctx context.Context, q pggen.Querier) (err error) {
		set, existing, err = s.db.getVariableSetByVariableID(ctx, variableID)
		if err != nil {
			return err
		}
		subject, err = s.organization.CanAccess(ctx, rbac.UpdateVariableSetAction, set.Organization)
		if err != nil {
			return err
		}

		// update a copy of variable
		updated = *existing
		if err := s.update(&updated, opts); err != nil {
			return err
		}

		// check for conflicts
		sets, err := s.db.listVariableSets(ctx, set.Organization)
		if err != nil {
			return err
		}
		if err := checkConflicts(&updated, set, sets); err != nil {
			return err
		}

		if err := s.db.updateVariable(ctx, &updated); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.Error(err, "updating variable set variable", "subject", subject, "variable_id", variableID)
		return nil, err
	}
	s.V(1).Info("updated variable set variable", "subject", subject, "set", set, "before", existing, "after", updated)

	return &updated, nil
}

func (s *service) deleteVariableFromSet(ctx context.Context, setID, variableID string) error {
	set, v, err := s.db.getVariableSetByVariableID(ctx, variableID)
	if err != nil {
		return err
	}

	subject, err := s.organization.CanAccess(ctx, rbac.RemoveVariableFromSetAction, set.Organization)
	if err != nil {
		return err
	}

	if err := s.db.deleteVariable(ctx, variableID); err != nil {
		s.Error(err, "deleting variable from set", "subject", subject, "variable", v, "set", set)
		return err
	}
	s.V(1).Info("deleted variable from set", "subject", subject, "variable", v, "set", set)

	return nil
}

func (s *service) applySetToWorkspaces(ctx context.Context, setID string, workspaceIDs []string) error {
	// retrieve set first in order to retrieve organization name for authorization
	set, err := s.db.getVariableSet(ctx, setID)
	if err != nil {
		return err
	}

	subject, err := s.organization.CanAccess(ctx, rbac.ApplyVariableSetToWorkspacesAction, set.Organization)
	if err != nil {
		return err
	}

	if err := s.db.createVariableSetWorkspaces(ctx, setID, workspaceIDs); err != nil {
		s.Error(err, "applying variable set to workspaces", "subject", subject, "set", set, "workspaces", workspaceIDs)
		return err
	}
	s.V(1).Info("applied variable set to workspaces", "subject", subject, "set", set, "workspaces", workspaceIDs)

	return nil
}

func (s *service) deleteSetFromWorkspaces(ctx context.Context, setID string, workspaceIDs []string) error {
	// retrieve set first in order to retrieve organization name for authorization
	set, err := s.db.getVariableSet(ctx, setID)
	if err != nil {
		return err
	}

	subject, err := s.organization.CanAccess(ctx, rbac.DeleteVariableSetFromWorkspacesAction, set.Organization)
	if err != nil {
		return err
	}

	if err := s.db.deleteVariableSetWorkspaces(ctx, setID, workspaceIDs); err != nil {
		s.Error(err, "removing variable set from workspaces", "subject", subject, "set", set, "workspaces", workspaceIDs)
		return err
	}
	s.V(1).Info("removed variable set from workspaces", "subject", subject, "set", set, "workspaces", workspaceIDs)

	return nil
}
