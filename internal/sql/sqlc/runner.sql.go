// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: runner.sql

package sqlc

import (
	"context"
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leg100/otf/internal/resource"
)

const deleteRunner = `-- name: DeleteRunner :one
DELETE
FROM runners
WHERE runner_id = $1
RETURNING runner_id, name, version, max_jobs, ip_address, last_ping_at, last_status_at, status, agent_pool_id
`

func (q *Queries) DeleteRunner(ctx context.Context, runnerID resource.ID) (Runner, error) {
	row := q.db.QueryRow(ctx, deleteRunner, runnerID)
	var i Runner
	err := row.Scan(
		&i.RunnerID,
		&i.Name,
		&i.Version,
		&i.MaxJobs,
		&i.IPAddress,
		&i.LastPingAt,
		&i.LastStatusAt,
		&i.Status,
		&i.AgentPoolID,
	)
	return i, err
}

const findRunnerByID = `-- name: FindRunnerByID :one
SELECT
    a.runner_id, a.name, a.version, a.max_jobs, a.ip_address, a.last_ping_at, a.last_status_at, a.status, a.agent_pool_id,
    ap::"agent_pools" AS agent_pool,
    ( SELECT count(*)
      FROM jobs j
      WHERE a.runner_id = j.runner_id
      AND j.status IN ('allocated', 'running')
    ) AS current_jobs
FROM runners a
LEFT JOIN agent_pools ap USING (agent_pool_id)
LEFT JOIN jobs j USING (runner_id)
WHERE a.runner_id = $1
`

type FindRunnerByIDRow struct {
	RunnerID     resource.ID
	Name         pgtype.Text
	Version      pgtype.Text
	MaxJobs      pgtype.Int4
	IPAddress    netip.Addr
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	Status       pgtype.Text
	AgentPoolID  *resource.ID
	AgentPool    *AgentPool
	CurrentJobs  int64
}

func (q *Queries) FindRunnerByID(ctx context.Context, runnerID resource.ID) (FindRunnerByIDRow, error) {
	row := q.db.QueryRow(ctx, findRunnerByID, runnerID)
	var i FindRunnerByIDRow
	err := row.Scan(
		&i.RunnerID,
		&i.Name,
		&i.Version,
		&i.MaxJobs,
		&i.IPAddress,
		&i.LastPingAt,
		&i.LastStatusAt,
		&i.Status,
		&i.AgentPoolID,
		&i.AgentPool,
		&i.CurrentJobs,
	)
	return i, err
}

const findRunnerByIDForUpdate = `-- name: FindRunnerByIDForUpdate :one
SELECT
    a.runner_id, a.name, a.version, a.max_jobs, a.ip_address, a.last_ping_at, a.last_status_at, a.status, a.agent_pool_id,
    ap::"agent_pools" AS agent_pool,
    ( SELECT count(*)
      FROM jobs j
      WHERE a.runner_id = j.runner_id
      AND j.status IN ('allocated', 'running')
    ) AS current_jobs
FROM runners a
LEFT JOIN agent_pools ap USING (agent_pool_id)
WHERE a.runner_id = $1
FOR UPDATE OF a
`

type FindRunnerByIDForUpdateRow struct {
	RunnerID     resource.ID
	Name         pgtype.Text
	Version      pgtype.Text
	MaxJobs      pgtype.Int4
	IPAddress    netip.Addr
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	Status       pgtype.Text
	AgentPoolID  *resource.ID
	AgentPool    *AgentPool
	CurrentJobs  int64
}

func (q *Queries) FindRunnerByIDForUpdate(ctx context.Context, runnerID resource.ID) (FindRunnerByIDForUpdateRow, error) {
	row := q.db.QueryRow(ctx, findRunnerByIDForUpdate, runnerID)
	var i FindRunnerByIDForUpdateRow
	err := row.Scan(
		&i.RunnerID,
		&i.Name,
		&i.Version,
		&i.MaxJobs,
		&i.IPAddress,
		&i.LastPingAt,
		&i.LastStatusAt,
		&i.Status,
		&i.AgentPoolID,
		&i.AgentPool,
		&i.CurrentJobs,
	)
	return i, err
}

const findRunners = `-- name: FindRunners :many
SELECT
    a.runner_id, a.name, a.version, a.max_jobs, a.ip_address, a.last_ping_at, a.last_status_at, a.status, a.agent_pool_id,
    ap::"agent_pools" AS agent_pool,
    ( SELECT count(*)
      FROM jobs j
      WHERE a.runner_id = j.runner_id
      AND j.status IN ('allocated', 'running')
    ) AS current_jobs
FROM runners a
LEFT JOIN agent_pools ap USING (agent_pool_id)
ORDER BY a.last_ping_at DESC
`

type FindRunnersRow struct {
	RunnerID     resource.ID
	Name         pgtype.Text
	Version      pgtype.Text
	MaxJobs      pgtype.Int4
	IPAddress    netip.Addr
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	Status       pgtype.Text
	AgentPoolID  *resource.ID
	AgentPool    *AgentPool
	CurrentJobs  int64
}

func (q *Queries) FindRunners(ctx context.Context) ([]FindRunnersRow, error) {
	rows, err := q.db.Query(ctx, findRunners)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindRunnersRow
	for rows.Next() {
		var i FindRunnersRow
		if err := rows.Scan(
			&i.RunnerID,
			&i.Name,
			&i.Version,
			&i.MaxJobs,
			&i.IPAddress,
			&i.LastPingAt,
			&i.LastStatusAt,
			&i.Status,
			&i.AgentPoolID,
			&i.AgentPool,
			&i.CurrentJobs,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findRunnersByOrganization = `-- name: FindRunnersByOrganization :many
SELECT
    a.runner_id, a.name, a.version, a.max_jobs, a.ip_address, a.last_ping_at, a.last_status_at, a.status, a.agent_pool_id,
    ap::"agent_pools" AS agent_pool,
    ( SELECT count(*)
      FROM jobs j
      WHERE a.runner_id = j.runner_id
      AND j.status IN ('allocated', 'running')
    ) AS current_jobs
FROM runners a
JOIN agent_pools ap USING (agent_pool_id)
WHERE ap.organization_name = $1
ORDER BY last_ping_at DESC
`

type FindRunnersByOrganizationRow struct {
	RunnerID     resource.ID
	Name         pgtype.Text
	Version      pgtype.Text
	MaxJobs      pgtype.Int4
	IPAddress    netip.Addr
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	Status       pgtype.Text
	AgentPoolID  *resource.ID
	AgentPool    *AgentPool
	CurrentJobs  int64
}

func (q *Queries) FindRunnersByOrganization(ctx context.Context, organizationName pgtype.Text) ([]FindRunnersByOrganizationRow, error) {
	rows, err := q.db.Query(ctx, findRunnersByOrganization, organizationName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindRunnersByOrganizationRow
	for rows.Next() {
		var i FindRunnersByOrganizationRow
		if err := rows.Scan(
			&i.RunnerID,
			&i.Name,
			&i.Version,
			&i.MaxJobs,
			&i.IPAddress,
			&i.LastPingAt,
			&i.LastStatusAt,
			&i.Status,
			&i.AgentPoolID,
			&i.AgentPool,
			&i.CurrentJobs,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findRunnersByPoolID = `-- name: FindRunnersByPoolID :many
SELECT
    a.runner_id, a.name, a.version, a.max_jobs, a.ip_address, a.last_ping_at, a.last_status_at, a.status, a.agent_pool_id,
    ap::"agent_pools" AS agent_pool,
    ( SELECT count(*)
      FROM jobs j
      WHERE a.runner_id = j.runner_id
      AND j.status IN ('allocated', 'running')
    ) AS current_jobs
FROM runners a
JOIN agent_pools ap USING (agent_pool_id)
WHERE ap.agent_pool_id = $1
ORDER BY last_ping_at DESC
`

type FindRunnersByPoolIDRow struct {
	RunnerID     resource.ID
	Name         pgtype.Text
	Version      pgtype.Text
	MaxJobs      pgtype.Int4
	IPAddress    netip.Addr
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	Status       pgtype.Text
	AgentPoolID  *resource.ID
	AgentPool    *AgentPool
	CurrentJobs  int64
}

func (q *Queries) FindRunnersByPoolID(ctx context.Context, agentPoolID resource.ID) ([]FindRunnersByPoolIDRow, error) {
	rows, err := q.db.Query(ctx, findRunnersByPoolID, agentPoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindRunnersByPoolIDRow
	for rows.Next() {
		var i FindRunnersByPoolIDRow
		if err := rows.Scan(
			&i.RunnerID,
			&i.Name,
			&i.Version,
			&i.MaxJobs,
			&i.IPAddress,
			&i.LastPingAt,
			&i.LastStatusAt,
			&i.Status,
			&i.AgentPoolID,
			&i.AgentPool,
			&i.CurrentJobs,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findServerRunners = `-- name: FindServerRunners :many
SELECT
    a.runner_id, a.name, a.version, a.max_jobs, a.ip_address, a.last_ping_at, a.last_status_at, a.status, a.agent_pool_id,
    ap::"agent_pools" AS agent_pool,
    ( SELECT count(*)
      FROM jobs j
      WHERE a.runner_id = j.runner_id
      AND j.status IN ('allocated', 'running')
    ) AS current_jobs
FROM runners a
LEFT JOIN agent_pools ap USING (agent_pool_id)
WHERE agent_pool_id IS NULL
ORDER BY last_ping_at DESC
`

type FindServerRunnersRow struct {
	RunnerID     resource.ID
	Name         pgtype.Text
	Version      pgtype.Text
	MaxJobs      pgtype.Int4
	IPAddress    netip.Addr
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	Status       pgtype.Text
	AgentPoolID  *resource.ID
	AgentPool    *AgentPool
	CurrentJobs  int64
}

func (q *Queries) FindServerRunners(ctx context.Context) ([]FindServerRunnersRow, error) {
	rows, err := q.db.Query(ctx, findServerRunners)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindServerRunnersRow
	for rows.Next() {
		var i FindServerRunnersRow
		if err := rows.Scan(
			&i.RunnerID,
			&i.Name,
			&i.Version,
			&i.MaxJobs,
			&i.IPAddress,
			&i.LastPingAt,
			&i.LastStatusAt,
			&i.Status,
			&i.AgentPoolID,
			&i.AgentPool,
			&i.CurrentJobs,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertRunner = `-- name: InsertRunner :exec
INSERT INTO runners (
    runner_id,
    name,
    version,
    max_jobs,
    ip_address,
    last_ping_at,
    last_status_at,
    status,
    agent_pool_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
`

type InsertRunnerParams struct {
	RunnerID     resource.ID
	Name         pgtype.Text
	Version      pgtype.Text
	MaxJobs      pgtype.Int4
	IPAddress    netip.Addr
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	Status       pgtype.Text
	AgentPoolID  *resource.ID
}

func (q *Queries) InsertRunner(ctx context.Context, arg InsertRunnerParams) error {
	_, err := q.db.Exec(ctx, insertRunner,
		arg.RunnerID,
		arg.Name,
		arg.Version,
		arg.MaxJobs,
		arg.IPAddress,
		arg.LastPingAt,
		arg.LastStatusAt,
		arg.Status,
		arg.AgentPoolID,
	)
	return err
}

const updateRunner = `-- name: UpdateRunner :one
UPDATE runners
SET status = $1,
    last_ping_at = $2,
    last_status_at = $3
WHERE runner_id = $4
RETURNING runner_id, name, version, max_jobs, ip_address, last_ping_at, last_status_at, status, agent_pool_id
`

type UpdateRunnerParams struct {
	Status       pgtype.Text
	LastPingAt   pgtype.Timestamptz
	LastStatusAt pgtype.Timestamptz
	RunnerID     resource.ID
}

func (q *Queries) UpdateRunner(ctx context.Context, arg UpdateRunnerParams) (Runner, error) {
	row := q.db.QueryRow(ctx, updateRunner,
		arg.Status,
		arg.LastPingAt,
		arg.LastStatusAt,
		arg.RunnerID,
	)
	var i Runner
	err := row.Scan(
		&i.RunnerID,
		&i.Name,
		&i.Version,
		&i.MaxJobs,
		&i.IPAddress,
		&i.LastPingAt,
		&i.LastStatusAt,
		&i.Status,
		&i.AgentPoolID,
	)
	return i, err
}
