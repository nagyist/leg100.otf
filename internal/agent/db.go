package agent

import (
	"github.com/jackc/pgtype"
	"github.com/leg100/otf/internal"
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
		LastSeen         pgtype.Timestamptz `json:"last_seen"`
		OrganizationName pgtype.Text        `json:"organization_name"`
	}
)
