package auth

import (
	"context"

	"github.com/leg100/otf/rbac"
)

type teamApp interface {
	createTeam(ctx context.Context, opts createTeamOptions) (*Team, error)
	getTeam(ctx context.Context, organization, team string) (*Team, error)
	getTeamByID(ctx context.Context, teamID string) (*Team, error)
	listTeams(ctx context.Context, organization string) ([]*Team, error)
	listTeamMembers(ctx context.Context, teamID string) ([]*User, error)
	updateTeam(ctx context.Context, teamID string, opts UpdateTeamOptions) (*Team, error)
}

func (a *Application) createTeam(ctx context.Context, opts createTeamOptions) (*Team, error) {
	subject, err := a.CanAccessOrganization(ctx, rbac.CreateTeamAction, opts.Organization)
	if err != nil {
		return nil, err
	}

	team := newTeam(opts)

	if err := a.db.createTeam(ctx, team); err != nil {
		a.Error(err, "creating team", "name", opts.Name, "organization", opts.Organization, "subject", subject)
		return nil, err
	}
	a.V(0).Info("created team", "name", opts.Name, "organization", opts.Organization, "subject", subject)

	return team, nil
}

func (a *Application) updateTeam(ctx context.Context, teamID string, opts UpdateTeamOptions) (*Team, error) {
	team, err := a.db.getTeamByID(ctx, teamID)
	if err != nil {
		a.Error(err, "retrieving team", "team_id", teamID)
		return nil, err
	}
	subject, err := a.CanAccessOrganization(ctx, rbac.UpdateTeamAction, team.Organization())
	if err != nil {
		return nil, err
	}

	team, err = a.db.UpdateTeam(ctx, teamID, func(team *Team) error {
		return team.Update(opts)
	})
	if err != nil {
		a.Error(err, "updating team", "name", team.Name(), "organization", team.Organization(), "subject", subject)
		return nil, err
	}

	a.V(2).Info("updated team", "name", team.Name(), "organization", team.Organization(), "subject", subject)

	return team, nil
}

func (a *Application) getTeam(ctx context.Context, organization, name string) (*Team, error) {
	subject, err := a.CanAccessOrganization(ctx, rbac.GetTeamAction, organization)
	if err != nil {
		return nil, err
	}

	team, err := a.db.getTeam(ctx, name, organization)
	if err != nil {
		a.Error(err, "retrieving team", "team", name, "organization", organization, "subject", subject)
		return nil, err
	}

	a.V(2).Info("retrieved team", "team", name, "organization", organization, "subject", subject)

	return team, nil
}

func (a *Application) getTeamByID(ctx context.Context, teamID string) (*Team, error) {
	team, err := a.db.getTeamByID(ctx, teamID)
	if err != nil {
		a.Error(err, "retrieving team", "team_id", teamID)
		return nil, err
	}

	subject, err := a.CanAccessOrganization(ctx, rbac.GetTeamAction, team.organization)
	if err != nil {
		return nil, err
	}

	a.V(2).Info("retrieved team", "team", team.name, "organization", team.organization, "subject", subject)

	return team, nil
}

// listTeams lists teams in the organization.
func (a *Application) listTeams(ctx context.Context, organization string) ([]*Team, error) {
	subject, err := a.CanAccessOrganization(ctx, rbac.ListTeamsAction, organization)
	if err != nil {
		return nil, err
	}

	teams, err := a.db.listTeams(ctx, organization)
	if err != nil {
		a.V(2).Info("listing teams", "organization", organization, "subject", subject)
		return nil, err
	}
	a.V(2).Info("listed teams", "organization", organization, "subject", subject)

	return teams, nil
}

// listTeamMembers lists users that are members of the given team. The caller
// needs either organization-wide authority to call this endpoint, or they need
// to be a member of the team.
func (a *Application) listTeamMembers(ctx context.Context, teamID string) ([]*User, error) {
	team, err := a.db.getTeamByID(ctx, teamID)
	if err != nil {
		a.Error(err, "retrieving team", "team_id", teamID)
		return nil, err
	}

	subject, err := a.CanAccessOrganization(ctx, rbac.ListUsersAction, team.Organization())
	if err != nil {
		return nil, err
	}

	members, err := a.db.listTeamMembers(ctx, teamID)
	if err != nil {
		a.Error(err, "listing team members", "team_id", teamID, "subject", subject)
		return nil, err
	}

	a.V(2).Info("listed team members", "team_id", teamID, "subject", subject)

	return members, nil
}