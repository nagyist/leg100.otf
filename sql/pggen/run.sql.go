// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const insertRunSQL = `INSERT INTO runs (
    run_id,
    plan_id,
    apply_id,
    created_at,
    is_destroy,
    position_in_queue,
    refresh,
    refresh_only,
    status,
    plan_status,
    apply_status,
    replace_addrs,
    target_addrs,
    planned_additions,
    planned_changes,
    planned_destructions,
    applied_additions,
    applied_changes,
    applied_destructions,
    configuration_version_id,
    workspace_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    $18,
    $19,
    $20,
    $21
);`

type InsertRunParams struct {
	ID                     string
	PlanID                 string
	ApplyID                string
	CreatedAt              time.Time
	IsDestroy              bool
	PositionInQueue        int
	Refresh                bool
	RefreshOnly            bool
	Status                 string
	PlanStatus             string
	ApplyStatus            string
	ReplaceAddrs           []string
	TargetAddrs            []string
	PlannedAdditions       int
	PlannedChanges         int
	PlannedDestructions    int
	AppliedAdditions       int
	AppliedChanges         int
	AppliedDestructions    int
	ConfigurationVersionID string
	WorkspaceID            string
}

// InsertRun implements Querier.InsertRun.
func (q *DBQuerier) InsertRun(ctx context.Context, params InsertRunParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRun")
	cmdTag, err := q.conn.Exec(ctx, insertRunSQL, params.ID, params.PlanID, params.ApplyID, params.CreatedAt, params.IsDestroy, params.PositionInQueue, params.Refresh, params.RefreshOnly, params.Status, params.PlanStatus, params.ApplyStatus, params.ReplaceAddrs, params.TargetAddrs, params.PlannedAdditions, params.PlannedChanges, params.PlannedDestructions, params.AppliedAdditions, params.AppliedChanges, params.AppliedDestructions, params.ConfigurationVersionID, params.WorkspaceID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertRun: %w", err)
	}
	return cmdTag, err
}

// InsertRunBatch implements Querier.InsertRunBatch.
func (q *DBQuerier) InsertRunBatch(batch genericBatch, params InsertRunParams) {
	batch.Queue(insertRunSQL, params.ID, params.PlanID, params.ApplyID, params.CreatedAt, params.IsDestroy, params.PositionInQueue, params.Refresh, params.RefreshOnly, params.Status, params.PlanStatus, params.ApplyStatus, params.ReplaceAddrs, params.TargetAddrs, params.PlannedAdditions, params.PlannedChanges, params.PlannedDestructions, params.AppliedAdditions, params.AppliedChanges, params.AppliedDestructions, params.ConfigurationVersionID, params.WorkspaceID)
}

// InsertRunScan implements Querier.InsertRunScan.
func (q *DBQuerier) InsertRunScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertRunBatch: %w", err)
	}
	return cmdTag, err
}

const insertRunStatusTimestampSQL = `INSERT INTO run_status_timestamps (
    run_id,
    status,
    timestamp
) VALUES (
    $1,
    $2,
    $3
);`

type InsertRunStatusTimestampParams struct {
	ID        string
	Status    string
	Timestamp time.Time
}

// InsertRunStatusTimestamp implements Querier.InsertRunStatusTimestamp.
func (q *DBQuerier) InsertRunStatusTimestamp(ctx context.Context, params InsertRunStatusTimestampParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRunStatusTimestamp")
	cmdTag, err := q.conn.Exec(ctx, insertRunStatusTimestampSQL, params.ID, params.Status, params.Timestamp)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertRunStatusTimestamp: %w", err)
	}
	return cmdTag, err
}

// InsertRunStatusTimestampBatch implements Querier.InsertRunStatusTimestampBatch.
func (q *DBQuerier) InsertRunStatusTimestampBatch(batch genericBatch, params InsertRunStatusTimestampParams) {
	batch.Queue(insertRunStatusTimestampSQL, params.ID, params.Status, params.Timestamp)
}

// InsertRunStatusTimestampScan implements Querier.InsertRunStatusTimestampScan.
func (q *DBQuerier) InsertRunStatusTimestampScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertRunStatusTimestampBatch: %w", err)
	}
	return cmdTag, err
}

const findRunsSQL = `SELECT
    runs.run_id,
    runs.plan_id,
    runs.apply_id,
    runs.created_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.status,
    runs.plan_status,
    runs.apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    runs.planned_additions,
    runs.planned_changes,
    runs.planned_destructions,
    runs.applied_additions,
    runs.applied_changes,
    runs.applied_destructions,
    runs.configuration_version_id,
    runs.workspace_id,
    configuration_versions.speculative,
    workspaces.auto_apply,
    CASE WHEN $1 THEN (configuration_versions.*)::"configuration_versions" END AS configuration_version,
    CASE WHEN $2 THEN (workspaces.*)::"workspaces" END AS workspace,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(pst.*) AS plan_status_timestamps
        FROM plan_status_timestamps pst
        WHERE pst.run_id = runs.run_id
        GROUP BY run_id
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(ast.*) AS apply_status_timestamps
        FROM apply_status_timestamps ast
        WHERE ast.run_id = runs.run_id
        GROUP BY run_id
    ) AS apply_status_timestamps
FROM runs
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
JOIN organizations USING(organization_id)
WHERE runs.workspace_id LIKE ANY($3)
AND runs.status LIKE ANY($4)
ORDER BY runs.created_at ASC
LIMIT $5 OFFSET $6
;`

type FindRunsParams struct {
	IncludeConfigurationVersion bool
	IncludeWorkspace            bool
	WorkspaceIds                []string
	Statuses                    []string
	Limit                       int
	Offset                      int
}

type FindRunsRow struct {
	RunID                  string                  `json:"run_id"`
	PlanID                 string                  `json:"plan_id"`
	ApplyID                string                  `json:"apply_id"`
	CreatedAt              time.Time               `json:"created_at"`
	IsDestroy              bool                    `json:"is_destroy"`
	PositionInQueue        int                     `json:"position_in_queue"`
	Refresh                bool                    `json:"refresh"`
	RefreshOnly            bool                    `json:"refresh_only"`
	Status                 string                  `json:"status"`
	PlanStatus             string                  `json:"plan_status"`
	ApplyStatus            string                  `json:"apply_status"`
	ReplaceAddrs           []string                `json:"replace_addrs"`
	TargetAddrs            []string                `json:"target_addrs"`
	PlannedAdditions       int                     `json:"planned_additions"`
	PlannedChanges         int                     `json:"planned_changes"`
	PlannedDestructions    int                     `json:"planned_destructions"`
	AppliedAdditions       int                     `json:"applied_additions"`
	AppliedChanges         int                     `json:"applied_changes"`
	AppliedDestructions    int                     `json:"applied_destructions"`
	ConfigurationVersionID string                  `json:"configuration_version_id"`
	WorkspaceID            string                  `json:"workspace_id"`
	Speculative            bool                    `json:"speculative"`
	AutoApply              bool                    `json:"auto_apply"`
	ConfigurationVersion   *ConfigurationVersions  `json:"configuration_version"`
	Workspace              *Workspaces             `json:"workspace"`
	RunStatusTimestamps    []RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []PlanStatusTimestamps  `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []ApplyStatusTimestamps `json:"apply_status_timestamps"`
}

// FindRuns implements Querier.FindRuns.
func (q *DBQuerier) FindRuns(ctx context.Context, params FindRunsParams) ([]FindRunsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRuns")
	rows, err := q.conn.Query(ctx, findRunsSQL, params.IncludeConfigurationVersion, params.IncludeWorkspace, params.WorkspaceIds, params.Statuses, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindRuns: %w", err)
	}
	defer rows.Close()
	items := []FindRunsRow{}
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPlanStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newApplyStatusTimestampsArray()
	for rows.Next() {
		var item FindRunsRow
		if err := rows.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlannedAdditions, &item.PlannedChanges, &item.PlannedDestructions, &item.AppliedAdditions, &item.AppliedChanges, &item.AppliedDestructions, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, configurationVersionRow, workspaceRow, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
			return nil, fmt.Errorf("scan FindRuns row: %w", err)
		}
		if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRuns rows: %w", err)
	}
	return items, err
}

// FindRunsBatch implements Querier.FindRunsBatch.
func (q *DBQuerier) FindRunsBatch(batch genericBatch, params FindRunsParams) {
	batch.Queue(findRunsSQL, params.IncludeConfigurationVersion, params.IncludeWorkspace, params.WorkspaceIds, params.Statuses, params.Limit, params.Offset)
}

// FindRunsScan implements Querier.FindRunsScan.
func (q *DBQuerier) FindRunsScan(results pgx.BatchResults) ([]FindRunsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindRunsBatch: %w", err)
	}
	defer rows.Close()
	items := []FindRunsRow{}
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPlanStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newApplyStatusTimestampsArray()
	for rows.Next() {
		var item FindRunsRow
		if err := rows.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlannedAdditions, &item.PlannedChanges, &item.PlannedDestructions, &item.AppliedAdditions, &item.AppliedChanges, &item.AppliedDestructions, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, configurationVersionRow, workspaceRow, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
			return nil, fmt.Errorf("scan FindRunsBatch row: %w", err)
		}
		if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRunsBatch rows: %w", err)
	}
	return items, err
}

const countRunsSQL = `SELECT count(*)
FROM runs
WHERE workspace_id LIKE ANY($1)
AND status LIKE ANY($2)
;`

// CountRuns implements Querier.CountRuns.
func (q *DBQuerier) CountRuns(ctx context.Context, workspaceIds []string, statuses []string) (*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountRuns")
	row := q.conn.QueryRow(ctx, countRunsSQL, workspaceIds, statuses)
	var item int
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("query CountRuns: %w", err)
	}
	return &item, nil
}

// CountRunsBatch implements Querier.CountRunsBatch.
func (q *DBQuerier) CountRunsBatch(batch genericBatch, workspaceIds []string, statuses []string) {
	batch.Queue(countRunsSQL, workspaceIds, statuses)
}

// CountRunsScan implements Querier.CountRunsScan.
func (q *DBQuerier) CountRunsScan(results pgx.BatchResults) (*int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("scan CountRunsBatch row: %w", err)
	}
	return &item, nil
}

const findRunByIDSQL = `SELECT
    runs.run_id,
    runs.plan_id,
    runs.apply_id,
    runs.created_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.status,
    runs.plan_status,
    runs.apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    runs.planned_additions,
    runs.planned_changes,
    runs.planned_destructions,
    runs.applied_additions,
    runs.applied_changes,
    runs.applied_destructions,
    runs.configuration_version_id,
    runs.workspace_id,
    configuration_versions.speculative,
    workspaces.auto_apply,
    CASE WHEN $1 THEN (configuration_versions.*)::"configuration_versions" END AS configuration_version,
    CASE WHEN $2 THEN (workspaces.*)::"workspaces" END AS workspace,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(pst.*) AS plan_status_timestamps
        FROM plan_status_timestamps pst
        WHERE pst.run_id = runs.run_id
        GROUP BY run_id
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(ast.*) AS apply_status_timestamps
        FROM apply_status_timestamps ast
        WHERE ast.run_id = runs.run_id
        GROUP BY run_id
    ) AS apply_status_timestamps
FROM runs
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
WHERE runs.run_id = $3
;`

type FindRunByIDParams struct {
	IncludeConfigurationVersion bool
	IncludeWorkspace            bool
	RunID                       string
}

type FindRunByIDRow struct {
	RunID                  string                  `json:"run_id"`
	PlanID                 string                  `json:"plan_id"`
	ApplyID                string                  `json:"apply_id"`
	CreatedAt              time.Time               `json:"created_at"`
	IsDestroy              bool                    `json:"is_destroy"`
	PositionInQueue        int                     `json:"position_in_queue"`
	Refresh                bool                    `json:"refresh"`
	RefreshOnly            bool                    `json:"refresh_only"`
	Status                 string                  `json:"status"`
	PlanStatus             string                  `json:"plan_status"`
	ApplyStatus            string                  `json:"apply_status"`
	ReplaceAddrs           []string                `json:"replace_addrs"`
	TargetAddrs            []string                `json:"target_addrs"`
	PlannedAdditions       int                     `json:"planned_additions"`
	PlannedChanges         int                     `json:"planned_changes"`
	PlannedDestructions    int                     `json:"planned_destructions"`
	AppliedAdditions       int                     `json:"applied_additions"`
	AppliedChanges         int                     `json:"applied_changes"`
	AppliedDestructions    int                     `json:"applied_destructions"`
	ConfigurationVersionID string                  `json:"configuration_version_id"`
	WorkspaceID            string                  `json:"workspace_id"`
	Speculative            bool                    `json:"speculative"`
	AutoApply              bool                    `json:"auto_apply"`
	ConfigurationVersion   *ConfigurationVersions  `json:"configuration_version"`
	Workspace              *Workspaces             `json:"workspace"`
	RunStatusTimestamps    []RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []PlanStatusTimestamps  `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []ApplyStatusTimestamps `json:"apply_status_timestamps"`
}

// FindRunByID implements Querier.FindRunByID.
func (q *DBQuerier) FindRunByID(ctx context.Context, params FindRunByIDParams) (FindRunByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByID")
	row := q.conn.QueryRow(ctx, findRunByIDSQL, params.IncludeConfigurationVersion, params.IncludeWorkspace, params.RunID)
	var item FindRunByIDRow
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPlanStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newApplyStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlannedAdditions, &item.PlannedChanges, &item.PlannedDestructions, &item.AppliedAdditions, &item.AppliedChanges, &item.AppliedDestructions, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, configurationVersionRow, workspaceRow, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("query FindRunByID: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	return item, nil
}

// FindRunByIDBatch implements Querier.FindRunByIDBatch.
func (q *DBQuerier) FindRunByIDBatch(batch genericBatch, params FindRunByIDParams) {
	batch.Queue(findRunByIDSQL, params.IncludeConfigurationVersion, params.IncludeWorkspace, params.RunID)
}

// FindRunByIDScan implements Querier.FindRunByIDScan.
func (q *DBQuerier) FindRunByIDScan(results pgx.BatchResults) (FindRunByIDRow, error) {
	row := results.QueryRow()
	var item FindRunByIDRow
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPlanStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newApplyStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlannedAdditions, &item.PlannedChanges, &item.PlannedDestructions, &item.AppliedAdditions, &item.AppliedChanges, &item.AppliedDestructions, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, configurationVersionRow, workspaceRow, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("scan FindRunByIDBatch row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	return item, nil
}

const findRunIDByPlanIDSQL = `SELECT run_id
FROM runs
WHERE plan_id = $1
;`

// FindRunIDByPlanID implements Querier.FindRunIDByPlanID.
func (q *DBQuerier) FindRunIDByPlanID(ctx context.Context, planID string) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunIDByPlanID")
	row := q.conn.QueryRow(ctx, findRunIDByPlanIDSQL, planID)
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindRunIDByPlanID: %w", err)
	}
	return item, nil
}

// FindRunIDByPlanIDBatch implements Querier.FindRunIDByPlanIDBatch.
func (q *DBQuerier) FindRunIDByPlanIDBatch(batch genericBatch, planID string) {
	batch.Queue(findRunIDByPlanIDSQL, planID)
}

// FindRunIDByPlanIDScan implements Querier.FindRunIDByPlanIDScan.
func (q *DBQuerier) FindRunIDByPlanIDScan(results pgx.BatchResults) (string, error) {
	row := results.QueryRow()
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindRunIDByPlanIDBatch row: %w", err)
	}
	return item, nil
}

const findRunIDByApplyIDSQL = `SELECT run_id
FROM runs
WHERE apply_id = $1
;`

// FindRunIDByApplyID implements Querier.FindRunIDByApplyID.
func (q *DBQuerier) FindRunIDByApplyID(ctx context.Context, applyID string) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunIDByApplyID")
	row := q.conn.QueryRow(ctx, findRunIDByApplyIDSQL, applyID)
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindRunIDByApplyID: %w", err)
	}
	return item, nil
}

// FindRunIDByApplyIDBatch implements Querier.FindRunIDByApplyIDBatch.
func (q *DBQuerier) FindRunIDByApplyIDBatch(batch genericBatch, applyID string) {
	batch.Queue(findRunIDByApplyIDSQL, applyID)
}

// FindRunIDByApplyIDScan implements Querier.FindRunIDByApplyIDScan.
func (q *DBQuerier) FindRunIDByApplyIDScan(results pgx.BatchResults) (string, error) {
	row := results.QueryRow()
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindRunIDByApplyIDBatch row: %w", err)
	}
	return item, nil
}

const findRunByIDForUpdateSQL = `SELECT
    runs.run_id,
    runs.plan_id,
    runs.apply_id,
    runs.created_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.status,
    runs.plan_status,
    runs.apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    runs.planned_additions,
    runs.planned_changes,
    runs.planned_destructions,
    runs.applied_additions,
    runs.applied_changes,
    runs.applied_destructions,
    runs.configuration_version_id,
    runs.workspace_id,
    configuration_versions.speculative,
    workspaces.auto_apply,
    CASE WHEN $1 THEN (configuration_versions.*)::"configuration_versions" END AS configuration_version,
    CASE WHEN $2 THEN (workspaces.*)::"workspaces" END AS workspace,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(pst.*) AS plan_status_timestamps
        FROM plan_status_timestamps pst
        WHERE pst.run_id = runs.run_id
        GROUP BY run_id
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(ast.*) AS apply_status_timestamps
        FROM apply_status_timestamps ast
        WHERE ast.run_id = runs.run_id
        GROUP BY run_id
    ) AS apply_status_timestamps
FROM runs
JOIN configuration_versions USING(workspace_id)
JOIN workspaces USING(workspace_id)
WHERE runs.run_id = $3
FOR UPDATE
;`

type FindRunByIDForUpdateParams struct {
	IncludeConfigurationVersion bool
	IncludeWorkspace            bool
	RunID                       string
}

type FindRunByIDForUpdateRow struct {
	RunID                  string                  `json:"run_id"`
	PlanID                 string                  `json:"plan_id"`
	ApplyID                string                  `json:"apply_id"`
	CreatedAt              time.Time               `json:"created_at"`
	IsDestroy              bool                    `json:"is_destroy"`
	PositionInQueue        int                     `json:"position_in_queue"`
	Refresh                bool                    `json:"refresh"`
	RefreshOnly            bool                    `json:"refresh_only"`
	Status                 string                  `json:"status"`
	PlanStatus             string                  `json:"plan_status"`
	ApplyStatus            string                  `json:"apply_status"`
	ReplaceAddrs           []string                `json:"replace_addrs"`
	TargetAddrs            []string                `json:"target_addrs"`
	PlannedAdditions       int                     `json:"planned_additions"`
	PlannedChanges         int                     `json:"planned_changes"`
	PlannedDestructions    int                     `json:"planned_destructions"`
	AppliedAdditions       int                     `json:"applied_additions"`
	AppliedChanges         int                     `json:"applied_changes"`
	AppliedDestructions    int                     `json:"applied_destructions"`
	ConfigurationVersionID string                  `json:"configuration_version_id"`
	WorkspaceID            string                  `json:"workspace_id"`
	Speculative            bool                    `json:"speculative"`
	AutoApply              bool                    `json:"auto_apply"`
	ConfigurationVersion   *ConfigurationVersions  `json:"configuration_version"`
	Workspace              *Workspaces             `json:"workspace"`
	RunStatusTimestamps    []RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []PlanStatusTimestamps  `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []ApplyStatusTimestamps `json:"apply_status_timestamps"`
}

// FindRunByIDForUpdate implements Querier.FindRunByIDForUpdate.
func (q *DBQuerier) FindRunByIDForUpdate(ctx context.Context, params FindRunByIDForUpdateParams) (FindRunByIDForUpdateRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByIDForUpdate")
	row := q.conn.QueryRow(ctx, findRunByIDForUpdateSQL, params.IncludeConfigurationVersion, params.IncludeWorkspace, params.RunID)
	var item FindRunByIDForUpdateRow
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPlanStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newApplyStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlannedAdditions, &item.PlannedChanges, &item.PlannedDestructions, &item.AppliedAdditions, &item.AppliedChanges, &item.AppliedDestructions, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, configurationVersionRow, workspaceRow, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("query FindRunByIDForUpdate: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	return item, nil
}

// FindRunByIDForUpdateBatch implements Querier.FindRunByIDForUpdateBatch.
func (q *DBQuerier) FindRunByIDForUpdateBatch(batch genericBatch, params FindRunByIDForUpdateParams) {
	batch.Queue(findRunByIDForUpdateSQL, params.IncludeConfigurationVersion, params.IncludeWorkspace, params.RunID)
}

// FindRunByIDForUpdateScan implements Querier.FindRunByIDForUpdateScan.
func (q *DBQuerier) FindRunByIDForUpdateScan(results pgx.BatchResults) (FindRunByIDForUpdateRow, error) {
	row := results.QueryRow()
	var item FindRunByIDForUpdateRow
	configurationVersionRow := q.types.newConfigurationVersions()
	workspaceRow := q.types.newWorkspaces()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPlanStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newApplyStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlannedAdditions, &item.PlannedChanges, &item.PlannedDestructions, &item.AppliedAdditions, &item.AppliedChanges, &item.AppliedDestructions, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, configurationVersionRow, workspaceRow, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("scan FindRunByIDForUpdateBatch row: %w", err)
	}
	if err := configurationVersionRow.AssignTo(&item.ConfigurationVersion); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := workspaceRow.AssignTo(&item.Workspace); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	return item, nil
}

const updateRunStatusSQL = `UPDATE runs
SET
    status = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdateRunStatus implements Querier.UpdateRunStatus.
func (q *DBQuerier) UpdateRunStatus(ctx context.Context, status string, id string) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRunStatus")
	row := q.conn.QueryRow(ctx, updateRunStatusSQL, status, id)
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateRunStatus: %w", err)
	}
	return item, nil
}

// UpdateRunStatusBatch implements Querier.UpdateRunStatusBatch.
func (q *DBQuerier) UpdateRunStatusBatch(batch genericBatch, status string, id string) {
	batch.Queue(updateRunStatusSQL, status, id)
}

// UpdateRunStatusScan implements Querier.UpdateRunStatusScan.
func (q *DBQuerier) UpdateRunStatusScan(results pgx.BatchResults) (string, error) {
	row := results.QueryRow()
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateRunStatusBatch row: %w", err)
	}
	return item, nil
}

const updateRunPlannedChangesByPlanIDSQL = `UPDATE runs
SET
    planned_additions = $1,
    planned_changes = $2,
    planned_destructions = $3
WHERE plan_id = $4
RETURNING plan_id
;`

type UpdateRunPlannedChangesByPlanIDParams struct {
	Additions    int
	Changes      int
	Destructions int
	PlanID       string
}

// UpdateRunPlannedChangesByPlanID implements Querier.UpdateRunPlannedChangesByPlanID.
func (q *DBQuerier) UpdateRunPlannedChangesByPlanID(ctx context.Context, params UpdateRunPlannedChangesByPlanIDParams) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRunPlannedChangesByPlanID")
	row := q.conn.QueryRow(ctx, updateRunPlannedChangesByPlanIDSQL, params.Additions, params.Changes, params.Destructions, params.PlanID)
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateRunPlannedChangesByPlanID: %w", err)
	}
	return item, nil
}

// UpdateRunPlannedChangesByPlanIDBatch implements Querier.UpdateRunPlannedChangesByPlanIDBatch.
func (q *DBQuerier) UpdateRunPlannedChangesByPlanIDBatch(batch genericBatch, params UpdateRunPlannedChangesByPlanIDParams) {
	batch.Queue(updateRunPlannedChangesByPlanIDSQL, params.Additions, params.Changes, params.Destructions, params.PlanID)
}

// UpdateRunPlannedChangesByPlanIDScan implements Querier.UpdateRunPlannedChangesByPlanIDScan.
func (q *DBQuerier) UpdateRunPlannedChangesByPlanIDScan(results pgx.BatchResults) (string, error) {
	row := results.QueryRow()
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateRunPlannedChangesByPlanIDBatch row: %w", err)
	}
	return item, nil
}

const updateRunAppliedChangesByApplyIDSQL = `UPDATE runs
SET
    applied_additions = $1,
    applied_changes = $2,
    applied_destructions = $3
WHERE apply_id = $4
RETURNING plan_id
;`

type UpdateRunAppliedChangesByApplyIDParams struct {
	Additions    int
	Changes      int
	Destructions int
	ApplyID      string
}

// UpdateRunAppliedChangesByApplyID implements Querier.UpdateRunAppliedChangesByApplyID.
func (q *DBQuerier) UpdateRunAppliedChangesByApplyID(ctx context.Context, params UpdateRunAppliedChangesByApplyIDParams) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRunAppliedChangesByApplyID")
	row := q.conn.QueryRow(ctx, updateRunAppliedChangesByApplyIDSQL, params.Additions, params.Changes, params.Destructions, params.ApplyID)
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateRunAppliedChangesByApplyID: %w", err)
	}
	return item, nil
}

// UpdateRunAppliedChangesByApplyIDBatch implements Querier.UpdateRunAppliedChangesByApplyIDBatch.
func (q *DBQuerier) UpdateRunAppliedChangesByApplyIDBatch(batch genericBatch, params UpdateRunAppliedChangesByApplyIDParams) {
	batch.Queue(updateRunAppliedChangesByApplyIDSQL, params.Additions, params.Changes, params.Destructions, params.ApplyID)
}

// UpdateRunAppliedChangesByApplyIDScan implements Querier.UpdateRunAppliedChangesByApplyIDScan.
func (q *DBQuerier) UpdateRunAppliedChangesByApplyIDScan(results pgx.BatchResults) (string, error) {
	row := results.QueryRow()
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateRunAppliedChangesByApplyIDBatch row: %w", err)
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
