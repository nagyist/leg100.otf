// Code generated by pggen. DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertRunSQL = `INSERT INTO runs (
    run_id,
    created_at,
    updated_at,
    is_destroy,
    position_in_queue,
    refresh,
    refresh_only,
    status,
    replace_addrs,
    target_addrs,
    workspace_id
) VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;`

type InsertRunParams struct {
	ID              string
	IsDestroy       bool
	PositionInQueue int32
	Refresh         bool
	RefreshOnly     bool
	Status          string
	ReplaceAddrs    []string
	TargetAddrs     []string
	WorkspaceID     string
}

type InsertRunRow struct {
	RunID                  string             `json:"run_id"`
	CreatedAt              pgtype.Timestamptz `json:"created_at"`
	UpdatedAt              pgtype.Timestamptz `json:"updated_at"`
	IsDestroy              bool               `json:"is_destroy"`
	PositionInQueue        int32              `json:"position_in_queue"`
	Refresh                bool               `json:"refresh"`
	RefreshOnly            bool               `json:"refresh_only"`
	Status                 string             `json:"status"`
	ReplaceAddrs           []string           `json:"replace_addrs"`
	TargetAddrs            []string           `json:"target_addrs"`
	WorkspaceID            string             `json:"workspace_id"`
	ConfigurationVersionID string             `json:"configuration_version_id"`
}

// InsertRun implements Querier.InsertRun.
func (q *DBQuerier) InsertRun(ctx context.Context, params InsertRunParams) (InsertRunRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRun")
	row := q.conn.QueryRow(ctx, insertRunSQL, params.ID, params.IsDestroy, params.PositionInQueue, params.Refresh, params.RefreshOnly, params.Status, params.ReplaceAddrs, params.TargetAddrs, params.WorkspaceID)
	var item InsertRunRow
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID); err != nil {
		return item, fmt.Errorf("query InsertRun: %w", err)
	}
	return item, nil
}

// InsertRunBatch implements Querier.InsertRunBatch.
func (q *DBQuerier) InsertRunBatch(batch genericBatch, params InsertRunParams) {
	batch.Queue(insertRunSQL, params.ID, params.IsDestroy, params.PositionInQueue, params.Refresh, params.RefreshOnly, params.Status, params.ReplaceAddrs, params.TargetAddrs, params.WorkspaceID)
}

// InsertRunScan implements Querier.InsertRunScan.
func (q *DBQuerier) InsertRunScan(results pgx.BatchResults) (InsertRunRow, error) {
	row := results.QueryRow()
	var item InsertRunRow
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID); err != nil {
		return item, fmt.Errorf("scan InsertRunBatch row: %w", err)
	}
	return item, nil
}

const insertRunStatusTimestampSQL = `INSERT INTO run_status_timestamps (
    run_id,
    status,
    timestamp
) VALUES (
    $1,
    $2,
    NOW()
)
RETURNING *;`

type InsertRunStatusTimestampRow struct {
	RunID     string             `json:"run_id"`
	Status    string             `json:"status"`
	Timestamp pgtype.Timestamptz `json:"timestamp"`
}

// InsertRunStatusTimestamp implements Querier.InsertRunStatusTimestamp.
func (q *DBQuerier) InsertRunStatusTimestamp(ctx context.Context, id string, status string) (InsertRunStatusTimestampRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRunStatusTimestamp")
	row := q.conn.QueryRow(ctx, insertRunStatusTimestampSQL, id, status)
	var item InsertRunStatusTimestampRow
	if err := row.Scan(&item.RunID, &item.Status, &item.Timestamp); err != nil {
		return item, fmt.Errorf("query InsertRunStatusTimestamp: %w", err)
	}
	return item, nil
}

// InsertRunStatusTimestampBatch implements Querier.InsertRunStatusTimestampBatch.
func (q *DBQuerier) InsertRunStatusTimestampBatch(batch genericBatch, id string, status string) {
	batch.Queue(insertRunStatusTimestampSQL, id, status)
}

// InsertRunStatusTimestampScan implements Querier.InsertRunStatusTimestampScan.
func (q *DBQuerier) InsertRunStatusTimestampScan(results pgx.BatchResults) (InsertRunStatusTimestampRow, error) {
	row := results.QueryRow()
	var item InsertRunStatusTimestampRow
	if err := row.Scan(&item.RunID, &item.Status, &item.Timestamp); err != nil {
		return item, fmt.Errorf("scan InsertRunStatusTimestampBatch row: %w", err)
	}
	return item, nil
}

const findRunsByWorkspaceIDSQL = `SELECT runs.*,
    (plans.*)::"plans" AS plan,
    (applies.*)::"applies" AS apply,
    (configuration_versions.*)::"configuration_versions" AS configuration_version,
    (workspaces.*)::"workspaces" AS workspace,
    count(*) OVER() AS full_count
FROM runs
JOIN plans USING(run_id)
JOIN applies USING(run_id)
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
WHERE workspaces.workspace_id = $1
LIMIT $2 OFFSET $3
;`

type FindRunsByWorkspaceIDParams struct {
	WorkspaceID string
	Limit       int
	Offset      int
}

type FindRunsByWorkspaceIDRow struct {
	RunID                  *string               `json:"run_id"`
	CreatedAt              pgtype.Timestamptz    `json:"created_at"`
	UpdatedAt              pgtype.Timestamptz    `json:"updated_at"`
	IsDestroy              *bool                 `json:"is_destroy"`
	PositionInQueue        *int32                `json:"position_in_queue"`
	Refresh                *bool                 `json:"refresh"`
	RefreshOnly            *bool                 `json:"refresh_only"`
	Status                 *string               `json:"status"`
	ReplaceAddrs           []string              `json:"replace_addrs"`
	TargetAddrs            []string              `json:"target_addrs"`
	WorkspaceID            *string               `json:"workspace_id"`
	ConfigurationVersionID *string               `json:"configuration_version_id"`
	Plan                   Plans                 `json:"plan"`
	Apply                  Applies               `json:"apply"`
	ConfigurationVersion   ConfigurationVersions `json:"configuration_version"`
	Workspace              Workspaces            `json:"workspace"`
	FullCount              *int                  `json:"full_count"`
}

// FindRunsByWorkspaceID implements Querier.FindRunsByWorkspaceID.
func (q *DBQuerier) FindRunsByWorkspaceID(ctx context.Context, params FindRunsByWorkspaceIDParams) ([]FindRunsByWorkspaceIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunsByWorkspaceID")
	rows, err := q.conn.Query(ctx, findRunsByWorkspaceIDSQL, params.WorkspaceID, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindRunsByWorkspaceID: %w", err)
	}
	defer rows.Close()
	items := []FindRunsByWorkspaceIDRow{}
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	for rows.Next() {
		var item FindRunsByWorkspaceIDRow
		if err := rows.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow, &item.FullCount); err != nil {
			return nil, fmt.Errorf("scan FindRunsByWorkspaceID row: %w", err)
		}
		if err := planRow.AssignTo(&item.Plan); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		if err := applyRow.AssignTo(&item.Apply); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRunsByWorkspaceID rows: %w", err)
	}
	return items, err
}

// FindRunsByWorkspaceIDBatch implements Querier.FindRunsByWorkspaceIDBatch.
func (q *DBQuerier) FindRunsByWorkspaceIDBatch(batch genericBatch, params FindRunsByWorkspaceIDParams) {
	batch.Queue(findRunsByWorkspaceIDSQL, params.WorkspaceID, params.Limit, params.Offset)
}

// FindRunsByWorkspaceIDScan implements Querier.FindRunsByWorkspaceIDScan.
func (q *DBQuerier) FindRunsByWorkspaceIDScan(results pgx.BatchResults) ([]FindRunsByWorkspaceIDRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindRunsByWorkspaceIDBatch: %w", err)
	}
	defer rows.Close()
	items := []FindRunsByWorkspaceIDRow{}
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	for rows.Next() {
		var item FindRunsByWorkspaceIDRow
		if err := rows.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow, &item.FullCount); err != nil {
			return nil, fmt.Errorf("scan FindRunsByWorkspaceIDBatch row: %w", err)
		}
		if err := planRow.AssignTo(&item.Plan); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		if err := applyRow.AssignTo(&item.Apply); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceID row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRunsByWorkspaceIDBatch rows: %w", err)
	}
	return items, err
}

const findRunsByWorkspaceNameSQL = `SELECT runs.*,
    (plans.*)::"plans" AS plan,
    (applies.*)::"applies" AS apply,
    (configuration_versions.*)::"configuration_versions" AS configuration_version,
    (workspaces.*)::"workspaces" AS workspace,
    count(*) OVER() AS full_count
FROM runs
JOIN plans USING(run_id)
JOIN applies USING(run_id)
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
JOIN organizations USING(organization_id)
WHERE workspaces.name = $1
AND organizations.name = $2
LIMIT $3 OFFSET $4
;`

type FindRunsByWorkspaceNameParams struct {
	WorkspaceName    string
	OrganizationName string
	Limit            int
	Offset           int
}

type FindRunsByWorkspaceNameRow struct {
	RunID                  *string               `json:"run_id"`
	CreatedAt              pgtype.Timestamptz    `json:"created_at"`
	UpdatedAt              pgtype.Timestamptz    `json:"updated_at"`
	IsDestroy              *bool                 `json:"is_destroy"`
	PositionInQueue        *int32                `json:"position_in_queue"`
	Refresh                *bool                 `json:"refresh"`
	RefreshOnly            *bool                 `json:"refresh_only"`
	Status                 *string               `json:"status"`
	ReplaceAddrs           []string              `json:"replace_addrs"`
	TargetAddrs            []string              `json:"target_addrs"`
	WorkspaceID            *string               `json:"workspace_id"`
	ConfigurationVersionID *string               `json:"configuration_version_id"`
	Plan                   Plans                 `json:"plan"`
	Apply                  Applies               `json:"apply"`
	ConfigurationVersion   ConfigurationVersions `json:"configuration_version"`
	Workspace              Workspaces            `json:"workspace"`
	FullCount              *int                  `json:"full_count"`
}

// FindRunsByWorkspaceName implements Querier.FindRunsByWorkspaceName.
func (q *DBQuerier) FindRunsByWorkspaceName(ctx context.Context, params FindRunsByWorkspaceNameParams) ([]FindRunsByWorkspaceNameRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunsByWorkspaceName")
	rows, err := q.conn.Query(ctx, findRunsByWorkspaceNameSQL, params.WorkspaceName, params.OrganizationName, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindRunsByWorkspaceName: %w", err)
	}
	defer rows.Close()
	items := []FindRunsByWorkspaceNameRow{}
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	for rows.Next() {
		var item FindRunsByWorkspaceNameRow
		if err := rows.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow, &item.FullCount); err != nil {
			return nil, fmt.Errorf("scan FindRunsByWorkspaceName row: %w", err)
		}
		if err := planRow.AssignTo(&item.Plan); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		if err := applyRow.AssignTo(&item.Apply); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRunsByWorkspaceName rows: %w", err)
	}
	return items, err
}

// FindRunsByWorkspaceNameBatch implements Querier.FindRunsByWorkspaceNameBatch.
func (q *DBQuerier) FindRunsByWorkspaceNameBatch(batch genericBatch, params FindRunsByWorkspaceNameParams) {
	batch.Queue(findRunsByWorkspaceNameSQL, params.WorkspaceName, params.OrganizationName, params.Limit, params.Offset)
}

// FindRunsByWorkspaceNameScan implements Querier.FindRunsByWorkspaceNameScan.
func (q *DBQuerier) FindRunsByWorkspaceNameScan(results pgx.BatchResults) ([]FindRunsByWorkspaceNameRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindRunsByWorkspaceNameBatch: %w", err)
	}
	defer rows.Close()
	items := []FindRunsByWorkspaceNameRow{}
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	for rows.Next() {
		var item FindRunsByWorkspaceNameRow
		if err := rows.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow, &item.FullCount); err != nil {
			return nil, fmt.Errorf("scan FindRunsByWorkspaceNameBatch row: %w", err)
		}
		if err := planRow.AssignTo(&item.Plan); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		if err := applyRow.AssignTo(&item.Apply); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
			return nil, fmt.Errorf("assign FindRunsByWorkspaceName row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRunsByWorkspaceNameBatch rows: %w", err)
	}
	return items, err
}

const findRunByIDSQL = `SELECT runs.*,
    (plans.*)::"plans" AS plan,
    (applies.*)::"applies" AS apply,
    (configuration_versions.*)::"configuration_versions" AS configuration_version,
    (workspaces.*)::"workspaces" AS workspace
FROM runs
JOIN plans USING(run_id)
JOIN applies USING(run_id)
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
WHERE runs.run_id = $1
LIMIT $2 OFFSET $3
;`

type FindRunByIDParams struct {
	RunID  string
	Limit  int
	Offset int
}

type FindRunByIDRow struct {
	RunID                  *string               `json:"run_id"`
	CreatedAt              pgtype.Timestamptz    `json:"created_at"`
	UpdatedAt              pgtype.Timestamptz    `json:"updated_at"`
	IsDestroy              *bool                 `json:"is_destroy"`
	PositionInQueue        *int32                `json:"position_in_queue"`
	Refresh                *bool                 `json:"refresh"`
	RefreshOnly            *bool                 `json:"refresh_only"`
	Status                 *string               `json:"status"`
	ReplaceAddrs           []string              `json:"replace_addrs"`
	TargetAddrs            []string              `json:"target_addrs"`
	WorkspaceID            *string               `json:"workspace_id"`
	ConfigurationVersionID *string               `json:"configuration_version_id"`
	Plan                   Plans                 `json:"plan"`
	Apply                  Applies               `json:"apply"`
	ConfigurationVersion   ConfigurationVersions `json:"configuration_version"`
	Workspace              Workspaces            `json:"workspace"`
}

// FindRunByID implements Querier.FindRunByID.
func (q *DBQuerier) FindRunByID(ctx context.Context, params FindRunByIDParams) (FindRunByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByID")
	row := q.conn.QueryRow(ctx, findRunByIDSQL, params.RunID, params.Limit, params.Offset)
	var item FindRunByIDRow
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow); err != nil {
		return item, fmt.Errorf("query FindRunByID: %w", err)
	}
	if err := planRow.AssignTo(&item.Plan); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := applyRow.AssignTo(&item.Apply); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	return item, nil
}

// FindRunByIDBatch implements Querier.FindRunByIDBatch.
func (q *DBQuerier) FindRunByIDBatch(batch genericBatch, params FindRunByIDParams) {
	batch.Queue(findRunByIDSQL, params.RunID, params.Limit, params.Offset)
}

// FindRunByIDScan implements Querier.FindRunByIDScan.
func (q *DBQuerier) FindRunByIDScan(results pgx.BatchResults) (FindRunByIDRow, error) {
	row := results.QueryRow()
	var item FindRunByIDRow
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow); err != nil {
		return item, fmt.Errorf("scan FindRunByIDBatch row: %w", err)
	}
	if err := planRow.AssignTo(&item.Plan); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := applyRow.AssignTo(&item.Apply); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	return item, nil
}

const findRunByPlanIDSQL = `SELECT runs.*,
    (plans.*)::"plans" AS plan,
    (applies.*)::"applies" AS apply,
    (configuration_versions.*)::"configuration_versions" AS configuration_version,
    (workspaces.*)::"workspaces" AS workspace
FROM runs
JOIN plans USING(run_id)
JOIN applies USING(run_id)
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
WHERE plans.plan_id = $1
LIMIT $2 OFFSET $3
;`

type FindRunByPlanIDParams struct {
	PlanID string
	Limit  int
	Offset int
}

type FindRunByPlanIDRow struct {
	RunID                  *string               `json:"run_id"`
	CreatedAt              pgtype.Timestamptz    `json:"created_at"`
	UpdatedAt              pgtype.Timestamptz    `json:"updated_at"`
	IsDestroy              *bool                 `json:"is_destroy"`
	PositionInQueue        *int32                `json:"position_in_queue"`
	Refresh                *bool                 `json:"refresh"`
	RefreshOnly            *bool                 `json:"refresh_only"`
	Status                 *string               `json:"status"`
	ReplaceAddrs           []string              `json:"replace_addrs"`
	TargetAddrs            []string              `json:"target_addrs"`
	WorkspaceID            *string               `json:"workspace_id"`
	ConfigurationVersionID *string               `json:"configuration_version_id"`
	Plan                   Plans                 `json:"plan"`
	Apply                  Applies               `json:"apply"`
	ConfigurationVersion   ConfigurationVersions `json:"configuration_version"`
	Workspace              Workspaces            `json:"workspace"`
}

// FindRunByPlanID implements Querier.FindRunByPlanID.
func (q *DBQuerier) FindRunByPlanID(ctx context.Context, params FindRunByPlanIDParams) (FindRunByPlanIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByPlanID")
	row := q.conn.QueryRow(ctx, findRunByPlanIDSQL, params.PlanID, params.Limit, params.Offset)
	var item FindRunByPlanIDRow
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow); err != nil {
		return item, fmt.Errorf("query FindRunByPlanID: %w", err)
	}
	if err := planRow.AssignTo(&item.Plan); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	if err := applyRow.AssignTo(&item.Apply); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	return item, nil
}

// FindRunByPlanIDBatch implements Querier.FindRunByPlanIDBatch.
func (q *DBQuerier) FindRunByPlanIDBatch(batch genericBatch, params FindRunByPlanIDParams) {
	batch.Queue(findRunByPlanIDSQL, params.PlanID, params.Limit, params.Offset)
}

// FindRunByPlanIDScan implements Querier.FindRunByPlanIDScan.
func (q *DBQuerier) FindRunByPlanIDScan(results pgx.BatchResults) (FindRunByPlanIDRow, error) {
	row := results.QueryRow()
	var item FindRunByPlanIDRow
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow); err != nil {
		return item, fmt.Errorf("scan FindRunByPlanIDBatch row: %w", err)
	}
	if err := planRow.AssignTo(&item.Plan); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	if err := applyRow.AssignTo(&item.Apply); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByPlanID row: %w", err)
	}
	return item, nil
}

const findRunByApplyIDSQL = `SELECT runs.*,
    (plans.*)::"plans" AS plan,
    (applies.*)::"applies" AS apply,
    (configuration_versions.*)::"configuration_versions" AS configuration_version,
    (workspaces.*)::"workspaces" AS workspace
FROM runs
JOIN plans USING(run_id)
JOIN applies USING(run_id)
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
WHERE applies.apply_id = $1
LIMIT $2 OFFSET $3
;`

type FindRunByApplyIDParams struct {
	ApplyID string
	Limit   int
	Offset  int
}

type FindRunByApplyIDRow struct {
	RunID                  *string               `json:"run_id"`
	CreatedAt              pgtype.Timestamptz    `json:"created_at"`
	UpdatedAt              pgtype.Timestamptz    `json:"updated_at"`
	IsDestroy              *bool                 `json:"is_destroy"`
	PositionInQueue        *int32                `json:"position_in_queue"`
	Refresh                *bool                 `json:"refresh"`
	RefreshOnly            *bool                 `json:"refresh_only"`
	Status                 *string               `json:"status"`
	ReplaceAddrs           []string              `json:"replace_addrs"`
	TargetAddrs            []string              `json:"target_addrs"`
	WorkspaceID            *string               `json:"workspace_id"`
	ConfigurationVersionID *string               `json:"configuration_version_id"`
	Plan                   Plans                 `json:"plan"`
	Apply                  Applies               `json:"apply"`
	ConfigurationVersion   ConfigurationVersions `json:"configuration_version"`
	Workspace              Workspaces            `json:"workspace"`
}

// FindRunByApplyID implements Querier.FindRunByApplyID.
func (q *DBQuerier) FindRunByApplyID(ctx context.Context, params FindRunByApplyIDParams) (FindRunByApplyIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByApplyID")
	row := q.conn.QueryRow(ctx, findRunByApplyIDSQL, params.ApplyID, params.Limit, params.Offset)
	var item FindRunByApplyIDRow
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow); err != nil {
		return item, fmt.Errorf("query FindRunByApplyID: %w", err)
	}
	if err := planRow.AssignTo(&item.Plan); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	if err := applyRow.AssignTo(&item.Apply); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	return item, nil
}

// FindRunByApplyIDBatch implements Querier.FindRunByApplyIDBatch.
func (q *DBQuerier) FindRunByApplyIDBatch(batch genericBatch, params FindRunByApplyIDParams) {
	batch.Queue(findRunByApplyIDSQL, params.ApplyID, params.Limit, params.Offset)
}

// FindRunByApplyIDScan implements Querier.FindRunByApplyIDScan.
func (q *DBQuerier) FindRunByApplyIDScan(results pgx.BatchResults) (FindRunByApplyIDRow, error) {
	row := results.QueryRow()
	var item FindRunByApplyIDRow
	planRow := q.types.newPlans()
	applyRow := q.types.newApplies()
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID, planRow, applyRow, configurationVersionRow, workspaceRow); err != nil {
		return item, fmt.Errorf("scan FindRunByApplyIDBatch row: %w", err)
	}
	if err := planRow.AssignTo(&item.Plan); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	if err := applyRow.AssignTo(&item.Apply); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByApplyID row: %w", err)
	}
	return item, nil
}

const updateRunStatusSQL = `UPDATE runs
SET
    status = $1,
    updated_at = NOW()
WHERE run_id = $2
RETURNING *;`

type UpdateRunStatusRow struct {
	RunID                  string             `json:"run_id"`
	CreatedAt              pgtype.Timestamptz `json:"created_at"`
	UpdatedAt              pgtype.Timestamptz `json:"updated_at"`
	IsDestroy              bool               `json:"is_destroy"`
	PositionInQueue        int32              `json:"position_in_queue"`
	Refresh                bool               `json:"refresh"`
	RefreshOnly            bool               `json:"refresh_only"`
	Status                 string             `json:"status"`
	ReplaceAddrs           []string           `json:"replace_addrs"`
	TargetAddrs            []string           `json:"target_addrs"`
	WorkspaceID            string             `json:"workspace_id"`
	ConfigurationVersionID string             `json:"configuration_version_id"`
}

// UpdateRunStatus implements Querier.UpdateRunStatus.
func (q *DBQuerier) UpdateRunStatus(ctx context.Context, status string, id string) (UpdateRunStatusRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRunStatus")
	row := q.conn.QueryRow(ctx, updateRunStatusSQL, status, id)
	var item UpdateRunStatusRow
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID); err != nil {
		return item, fmt.Errorf("query UpdateRunStatus: %w", err)
	}
	return item, nil
}

// UpdateRunStatusBatch implements Querier.UpdateRunStatusBatch.
func (q *DBQuerier) UpdateRunStatusBatch(batch genericBatch, status string, id string) {
	batch.Queue(updateRunStatusSQL, status, id)
}

// UpdateRunStatusScan implements Querier.UpdateRunStatusScan.
func (q *DBQuerier) UpdateRunStatusScan(results pgx.BatchResults) (UpdateRunStatusRow, error) {
	row := results.QueryRow()
	var item UpdateRunStatusRow
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.WorkspaceID, &item.ConfigurationVersionID); err != nil {
		return item, fmt.Errorf("scan UpdateRunStatusBatch row: %w", err)
	}
	return item, nil
}

const deleteRunByIDSQL = `DELETE
FROM runs
WHERE run_id = $1;`

// DeleteRunByID implements Querier.DeleteRunByID.
func (q *DBQuerier) DeleteRunByID(ctx context.Context, runID string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteRunByID")
	cmdTag, err := q.conn.Exec(ctx, deleteRunByIDSQL, runID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query DeleteRunByID: %w", err)
	}
	return cmdTag, err
}

// DeleteRunByIDBatch implements Querier.DeleteRunByIDBatch.
func (q *DBQuerier) DeleteRunByIDBatch(batch genericBatch, runID string) {
	batch.Queue(deleteRunByIDSQL, runID)
}

// DeleteRunByIDScan implements Querier.DeleteRunByIDScan.
func (q *DBQuerier) DeleteRunByIDScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec DeleteRunByIDBatch: %w", err)
	}
	return cmdTag, err
}
