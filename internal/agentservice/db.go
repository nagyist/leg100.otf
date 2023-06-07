package agentservice

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/sql"
	"github.com/leg100/otf/internal/sql/pggen"
)

type (
	// pgdb is an agent database on postgres
	pgdb struct {
		internal.DB // provides access to generated SQL queries
	}

	pgresult struct {
		AgentID          pgtype.Text        `json:"agent_id"`
		Status           pgtype.Text        `json:"status"`
		IpAddress        pgtype.Text        `json:"ip_address"`
		Version          pgtype.Text        `json:"version"`
		Name             pgtype.Text        `json:"name"`
		External         bool               `json:"external"`
		LastSeen         pgtype.Timestamptz `json:"last_seen"`
		OrganizationName pgtype.Text        `json:"organization_name"`
	}
)

func (db *pgdb) create(ctx context.Context, agent *Agent) error {
	params := pggen.InsertAgentParams{
		AgentID:          sql.String(agent.ID),
		Status:           sql.String(string(agent.Status)),
		IpAddress:        sql.String(agent.IPAddress),
		Version:          sql.String(agent.Version),
		External:         agent.External,
		LastSeen:         sql.Timestamptz(agent.LastSeen),
		Name:             sql.NullString(),
		OrganizationName: sql.NullString(),
	}
	if agent.Name != nil {
		params.Name = sql.String(*agent.Name)
	}
	if agent.Organization != nil {
		params.OrganizationName = sql.String(*agent.Organization)
	}

	_, err := db.InsertAgent(ctx, params)
	return sql.Error(err)
}

func (db *pgdb) updateStatus(ctx context.Context, id string, status Status) error {
	_, err := db.UpdateAgentStatus(ctx, pggen.UpdateAgentStatusParams{
		Status:   sql.String(string(status)),
		AgentID:  sql.String(id),
		LastSeen: sql.Timestamptz(internal.CurrentTimestamp()),
	})
	return sql.Error(err)
}
