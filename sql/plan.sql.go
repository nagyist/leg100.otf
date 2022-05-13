// Code generated by pggen. DO NOT EDIT.

package sql

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

const insertPlanStatusTimestampSQL = `INSERT INTO plan_status_timestamps (
    run_id,
    status,
    timestamp
) VALUES (
    $1,
    $2,
    current_timestamp
)
RETURNING *;`

type InsertPlanStatusTimestampRow struct {
	RunID     string    `json:"run_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func (s InsertPlanStatusTimestampRow) GetRunID() string { return s.RunID }
func (s InsertPlanStatusTimestampRow) GetStatus() string { return s.Status }
func (s InsertPlanStatusTimestampRow) GetTimestamp() time.Time { return s.Timestamp }


// InsertPlanStatusTimestamp implements Querier.InsertPlanStatusTimestamp.
func (q *DBQuerier) InsertPlanStatusTimestamp(ctx context.Context, id string, status string) (InsertPlanStatusTimestampRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertPlanStatusTimestamp")
	row := q.conn.QueryRow(ctx, insertPlanStatusTimestampSQL, id, status)
	var item InsertPlanStatusTimestampRow
	if err := row.Scan(&item.RunID, &item.Status, &item.Timestamp); err != nil {
		return item, fmt.Errorf("query InsertPlanStatusTimestamp: %w", err)
	}
	return item, nil
}

// InsertPlanStatusTimestampBatch implements Querier.InsertPlanStatusTimestampBatch.
func (q *DBQuerier) InsertPlanStatusTimestampBatch(batch genericBatch, id string, status string) {
	batch.Queue(insertPlanStatusTimestampSQL, id, status)
}

// InsertPlanStatusTimestampScan implements Querier.InsertPlanStatusTimestampScan.
func (q *DBQuerier) InsertPlanStatusTimestampScan(results pgx.BatchResults) (InsertPlanStatusTimestampRow, error) {
	row := results.QueryRow()
	var item InsertPlanStatusTimestampRow
	if err := row.Scan(&item.RunID, &item.Status, &item.Timestamp); err != nil {
		return item, fmt.Errorf("scan InsertPlanStatusTimestampBatch row: %w", err)
	}
	return item, nil
}

const updatePlanStatusSQL = `UPDATE runs
SET
    status = $1,
    updated_at = current_timestamp
WHERE plan_id = $2
RETURNING *;`

type UpdatePlanStatusRow struct {
	RunID                       string    `json:"run_id"`
	PlanID                      string    `json:"plan_id"`
	ApplyID                     string    `json:"apply_id"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
	IsDestroy                   bool      `json:"is_destroy"`
	PositionInQueue             int32     `json:"position_in_queue"`
	Refresh                     bool      `json:"refresh"`
	RefreshOnly                 bool      `json:"refresh_only"`
	Status                      string    `json:"status"`
	ReplaceAddrs                []string  `json:"replace_addrs"`
	TargetAddrs                 []string  `json:"target_addrs"`
	PlanStatus                  string    `json:"plan_status"`
	PlanBin                     []byte    `json:"plan_bin"`
	PlanJson                    []byte    `json:"plan_json"`
	PlannedResourceAdditions    *int32    `json:"planned_resource_additions"`
	PlannedResourceChanges      *int32    `json:"planned_resource_changes"`
	PlannedResourceDestructions *int32    `json:"planned_resource_destructions"`
	ApplyStatus                 string    `json:"apply_status"`
	AppliedResourceAdditions    *int32    `json:"applied_resource_additions"`
	AppliedResourceChanges      *int32    `json:"applied_resource_changes"`
	AppliedResourceDestructions *int32    `json:"applied_resource_destructions"`
	WorkspaceID                 string    `json:"workspace_id"`
	ConfigurationVersionID      string    `json:"configuration_version_id"`
}

func (s UpdatePlanStatusRow) GetRunID() string { return s.RunID }
func (s UpdatePlanStatusRow) GetPlanID() string { return s.PlanID }
func (s UpdatePlanStatusRow) GetApplyID() string { return s.ApplyID }
func (s UpdatePlanStatusRow) GetCreatedAt() time.Time { return s.CreatedAt }
func (s UpdatePlanStatusRow) GetUpdatedAt() time.Time { return s.UpdatedAt }
func (s UpdatePlanStatusRow) GetIsDestroy() bool { return s.IsDestroy }
func (s UpdatePlanStatusRow) GetPositionInQueue() int32 { return s.PositionInQueue }
func (s UpdatePlanStatusRow) GetRefresh() bool { return s.Refresh }
func (s UpdatePlanStatusRow) GetRefreshOnly() bool { return s.RefreshOnly }
func (s UpdatePlanStatusRow) GetStatus() string { return s.Status }
func (s UpdatePlanStatusRow) GetReplaceAddrs() []string { return s.ReplaceAddrs }
func (s UpdatePlanStatusRow) GetTargetAddrs() []string { return s.TargetAddrs }
func (s UpdatePlanStatusRow) GetPlanStatus() string { return s.PlanStatus }
func (s UpdatePlanStatusRow) GetPlanBin() []byte { return s.PlanBin }
func (s UpdatePlanStatusRow) GetPlanJson() []byte { return s.PlanJson }
func (s UpdatePlanStatusRow) GetPlannedResourceAdditions() *int32 { return s.PlannedResourceAdditions }
func (s UpdatePlanStatusRow) GetPlannedResourceChanges() *int32 { return s.PlannedResourceChanges }
func (s UpdatePlanStatusRow) GetPlannedResourceDestructions() *int32 { return s.PlannedResourceDestructions }
func (s UpdatePlanStatusRow) GetApplyStatus() string { return s.ApplyStatus }
func (s UpdatePlanStatusRow) GetAppliedResourceAdditions() *int32 { return s.AppliedResourceAdditions }
func (s UpdatePlanStatusRow) GetAppliedResourceChanges() *int32 { return s.AppliedResourceChanges }
func (s UpdatePlanStatusRow) GetAppliedResourceDestructions() *int32 { return s.AppliedResourceDestructions }
func (s UpdatePlanStatusRow) GetWorkspaceID() string { return s.WorkspaceID }
func (s UpdatePlanStatusRow) GetConfigurationVersionID() string { return s.ConfigurationVersionID }


// UpdatePlanStatus implements Querier.UpdatePlanStatus.
func (q *DBQuerier) UpdatePlanStatus(ctx context.Context, status string, id string) (UpdatePlanStatusRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlanStatus")
	row := q.conn.QueryRow(ctx, updatePlanStatusSQL, status, id)
	var item UpdatePlanStatusRow
	if err := row.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlanStatus, &item.PlanBin, &item.PlanJson, &item.PlannedResourceAdditions, &item.PlannedResourceChanges, &item.PlannedResourceDestructions, &item.ApplyStatus, &item.AppliedResourceAdditions, &item.AppliedResourceChanges, &item.AppliedResourceDestructions, &item.WorkspaceID, &item.ConfigurationVersionID); err != nil {
		return item, fmt.Errorf("query UpdatePlanStatus: %w", err)
	}
	return item, nil
}

// UpdatePlanStatusBatch implements Querier.UpdatePlanStatusBatch.
func (q *DBQuerier) UpdatePlanStatusBatch(batch genericBatch, status string, id string) {
	batch.Queue(updatePlanStatusSQL, status, id)
}

// UpdatePlanStatusScan implements Querier.UpdatePlanStatusScan.
func (q *DBQuerier) UpdatePlanStatusScan(results pgx.BatchResults) (UpdatePlanStatusRow, error) {
	row := results.QueryRow()
	var item UpdatePlanStatusRow
	if err := row.Scan(&item.RunID, &item.PlanID, &item.ApplyID, &item.CreatedAt, &item.UpdatedAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.ReplaceAddrs, &item.TargetAddrs, &item.PlanStatus, &item.PlanBin, &item.PlanJson, &item.PlannedResourceAdditions, &item.PlannedResourceChanges, &item.PlannedResourceDestructions, &item.ApplyStatus, &item.AppliedResourceAdditions, &item.AppliedResourceChanges, &item.AppliedResourceDestructions, &item.WorkspaceID, &item.ConfigurationVersionID); err != nil {
		return item, fmt.Errorf("scan UpdatePlanStatusBatch row: %w", err)
	}
	return item, nil
}

const getPlanBinByRunIDSQL = `SELECT plan_bin
FROM runs
WHERE run_id = $1
;`

// GetPlanBinByRunID implements Querier.GetPlanBinByRunID.
func (q *DBQuerier) GetPlanBinByRunID(ctx context.Context, runID string) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetPlanBinByRunID")
	row := q.conn.QueryRow(ctx, getPlanBinByRunIDSQL, runID)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query GetPlanBinByRunID: %w", err)
	}
	return item, nil
}

// GetPlanBinByRunIDBatch implements Querier.GetPlanBinByRunIDBatch.
func (q *DBQuerier) GetPlanBinByRunIDBatch(batch genericBatch, runID string) {
	batch.Queue(getPlanBinByRunIDSQL, runID)
}

// GetPlanBinByRunIDScan implements Querier.GetPlanBinByRunIDScan.
func (q *DBQuerier) GetPlanBinByRunIDScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan GetPlanBinByRunIDBatch row: %w", err)
	}
	return item, nil
}

const getPlanJSONByRunIDSQL = `SELECT plan_json
FROM runs
WHERE run_id = $1
;`

// GetPlanJSONByRunID implements Querier.GetPlanJSONByRunID.
func (q *DBQuerier) GetPlanJSONByRunID(ctx context.Context, runID string) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetPlanJSONByRunID")
	row := q.conn.QueryRow(ctx, getPlanJSONByRunIDSQL, runID)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query GetPlanJSONByRunID: %w", err)
	}
	return item, nil
}

// GetPlanJSONByRunIDBatch implements Querier.GetPlanJSONByRunIDBatch.
func (q *DBQuerier) GetPlanJSONByRunIDBatch(batch genericBatch, runID string) {
	batch.Queue(getPlanJSONByRunIDSQL, runID)
}

// GetPlanJSONByRunIDScan implements Querier.GetPlanJSONByRunIDScan.
func (q *DBQuerier) GetPlanJSONByRunIDScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan GetPlanJSONByRunIDBatch row: %w", err)
	}
	return item, nil
}

const putPlanBinByRunIDSQL = `UPDATE runs
SET plan_bin = $1
WHERE run_id = $2
;`

// PutPlanBinByRunID implements Querier.PutPlanBinByRunID.
func (q *DBQuerier) PutPlanBinByRunID(ctx context.Context, planBin []byte, runID string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "PutPlanBinByRunID")
	cmdTag, err := q.conn.Exec(ctx, putPlanBinByRunIDSQL, planBin, runID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query PutPlanBinByRunID: %w", err)
	}
	return cmdTag, err
}

// PutPlanBinByRunIDBatch implements Querier.PutPlanBinByRunIDBatch.
func (q *DBQuerier) PutPlanBinByRunIDBatch(batch genericBatch, planBin []byte, runID string) {
	batch.Queue(putPlanBinByRunIDSQL, planBin, runID)
}

// PutPlanBinByRunIDScan implements Querier.PutPlanBinByRunIDScan.
func (q *DBQuerier) PutPlanBinByRunIDScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec PutPlanBinByRunIDBatch: %w", err)
	}
	return cmdTag, err
}

const putPlanJSONByRunIDSQL = `UPDATE runs
SET plan_json = $1
WHERE run_id = $2
;`

// PutPlanJSONByRunID implements Querier.PutPlanJSONByRunID.
func (q *DBQuerier) PutPlanJSONByRunID(ctx context.Context, planJson []byte, runID string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "PutPlanJSONByRunID")
	cmdTag, err := q.conn.Exec(ctx, putPlanJSONByRunIDSQL, planJson, runID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query PutPlanJSONByRunID: %w", err)
	}
	return cmdTag, err
}

// PutPlanJSONByRunIDBatch implements Querier.PutPlanJSONByRunIDBatch.
func (q *DBQuerier) PutPlanJSONByRunIDBatch(batch genericBatch, planJson []byte, runID string) {
	batch.Queue(putPlanJSONByRunIDSQL, planJson, runID)
}

// PutPlanJSONByRunIDScan implements Querier.PutPlanJSONByRunIDScan.
func (q *DBQuerier) PutPlanJSONByRunIDScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec PutPlanJSONByRunIDBatch: %w", err)
	}
	return cmdTag, err
}

const updatePlanResourcesSQL = `UPDATE runs
SET
    planned_resource_additions = $1,
    planned_resource_changes = $2,
    planned_resource_destructions = $3
WHERE plan_id = $4
;`

type UpdatePlanResourcesParams struct {
	PlannedResourceAdditions    int32
	PlannedResourceChanges      int32
	PlannedResourceDestructions int32
	PlanID                      string
}

// UpdatePlanResources implements Querier.UpdatePlanResources.
func (q *DBQuerier) UpdatePlanResources(ctx context.Context, params UpdatePlanResourcesParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlanResources")
	cmdTag, err := q.conn.Exec(ctx, updatePlanResourcesSQL, params.PlannedResourceAdditions, params.PlannedResourceChanges, params.PlannedResourceDestructions, params.PlanID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query UpdatePlanResources: %w", err)
	}
	return cmdTag, err
}

// UpdatePlanResourcesBatch implements Querier.UpdatePlanResourcesBatch.
func (q *DBQuerier) UpdatePlanResourcesBatch(batch genericBatch, params UpdatePlanResourcesParams) {
	batch.Queue(updatePlanResourcesSQL, params.PlannedResourceAdditions, params.PlannedResourceChanges, params.PlannedResourceDestructions, params.PlanID)
}

// UpdatePlanResourcesScan implements Querier.UpdatePlanResourcesScan.
func (q *DBQuerier) UpdatePlanResourcesScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec UpdatePlanResourcesBatch: %w", err)
	}
	return cmdTag, err
}
