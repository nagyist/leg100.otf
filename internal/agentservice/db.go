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

	agentresult struct {
		AgentID          pgtype.Text        `json:"agent_id"`
		Status           pgtype.Text        `json:"status"`
		IpAddress        pgtype.Text        `json:"ip_address"`
		Version          pgtype.Text        `json:"version"`
		Name             pgtype.Text        `json:"name"`
		External         bool               `json:"external"`
		LastSeen         pgtype.Timestamptz `json:"last_seen"`
		OrganizationName pgtype.Text        `json:"organization_name"`
	}

	jobresult struct {
		JobID   pgtype.Text `json:"job_id"`
		RunID   pgtype.Text `json:"run_id"`
		Phase   pgtype.Text `json:"phase"`
		AgentID pgtype.Text `json:"agent_id"`
		Status  pgtype.Text `json:"status"`
	}
)

func (r agentresult) toAgent() *Agent {
	to := &Agent{
		ID:        r.AgentID.String,
		Status:    Status(r.Status.String),
		IPAddress: r.IpAddress.String,
		Version:   r.Version.String,
		External:  r.External,
		LastSeen:  r.LastSeen.Time.UTC(),
	}
	if r.Name.Status == pgtype.Present {
		to.Name = &r.Name.String
	}
	if r.OrganizationName.Status == pgtype.Present {
		to.Organization = &r.OrganizationName.String
	}
	return to
}

func (r jobresult) toJob() *Job {
	to := &Job{
		ID:     r.JobID.String,
		Status: JobStatus(r.Status.String),
		RunID:  r.RunID.String,
		Phase:  internal.PhaseType(r.Phase.String),
	}
	if r.AgentID.Status == pgtype.Present {
		to.AgentID = &r.AgentID.String
	}
	return to
}

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

func (db *pgdb) list(ctx context.Context) (agents []*Agent, err error) {
	results, err := db.FindAgents(ctx)
	if err != nil {
		return nil, sql.Error(err)
	}
	for _, res := range results {
		agents = append(agents, agentresult(res).toAgent())
	}
	return
}

func (db *pgdb) listByOrganization(ctx context.Context, organization string) (agents []*Agent, err error) {
	results, err := db.FindAgentsByOrganization(ctx, sql.String(organization))
	if err != nil {
		return nil, sql.Error(err)
	}
	for _, res := range results {
		agents = append(agents, agentresult(res).toAgent())
	}
	return
}

func (db *pgdb) getAssignedJobByAgentID(ctx context.Context, agentID string) (job *Job, err error) {
	result, err := db.FindAssignedJobByAgentID(ctx, sql.String(agentID))
	if err != nil {
		return nil, sql.Error(err)
	}
	return jobresult(result).toJob(), nil
}

func (db *pgdb) getJobByRunIDAndPhase(ctx context.Context, runID string, phase internal.PhaseType) (job *Job, err error) {
	result, err := db.FindJobByRunIDAndPhase(ctx, sql.String(runID), sql.String(string(phase)))
	if err != nil {
		return nil, sql.Error(err)
	}
	return jobresult(result).toJob(), nil
}
