// Code generated by pggen. DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertPlanSQL = `INSERT INTO plans (
    plan_id,
    created_at,
    updated_at,
    status,
    run_id
) VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;`

type InsertPlanParams struct {
	ID     string
	Status string
	RunID  string
}

type InsertPlanRow struct {
	PlanID               string             `json:"plan_id"`
	CreatedAt            pgtype.Timestamptz `json:"created_at"`
	UpdatedAt            pgtype.Timestamptz `json:"updated_at"`
	ResourceAdditions    int32              `json:"resource_additions"`
	ResourceChanges      int32              `json:"resource_changes"`
	ResourceDestructions int32              `json:"resource_destructions"`
	Status               string             `json:"status"`
	PlanFile             []byte             `json:"plan_file"`
	PlanJson             []byte             `json:"plan_json"`
	RunID                string             `json:"run_id"`
}

// InsertPlan implements Querier.InsertPlan.
func (q *DBQuerier) InsertPlan(ctx context.Context, params InsertPlanParams) (InsertPlanRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertPlan")
	row := q.conn.QueryRow(ctx, insertPlanSQL, params.ID, params.Status, params.RunID)
	var item InsertPlanRow
	if err := row.Scan(&item.PlanID, &item.CreatedAt, &item.UpdatedAt, &item.ResourceAdditions, &item.ResourceChanges, &item.ResourceDestructions, &item.Status, &item.PlanFile, &item.PlanJson, &item.RunID); err != nil {
		return item, fmt.Errorf("query InsertPlan: %w", err)
	}
	return item, nil
}

// InsertPlanBatch implements Querier.InsertPlanBatch.
func (q *DBQuerier) InsertPlanBatch(batch genericBatch, params InsertPlanParams) {
	batch.Queue(insertPlanSQL, params.ID, params.Status, params.RunID)
}

// InsertPlanScan implements Querier.InsertPlanScan.
func (q *DBQuerier) InsertPlanScan(results pgx.BatchResults) (InsertPlanRow, error) {
	row := results.QueryRow()
	var item InsertPlanRow
	if err := row.Scan(&item.PlanID, &item.CreatedAt, &item.UpdatedAt, &item.ResourceAdditions, &item.ResourceChanges, &item.ResourceDestructions, &item.Status, &item.PlanFile, &item.PlanJson, &item.RunID); err != nil {
		return item, fmt.Errorf("scan InsertPlanBatch row: %w", err)
	}
	return item, nil
}

const insertPlanStatusTimestampSQL = `INSERT INTO plan_status_timestamps (
    plan_id,
    status,
    timestamp
) VALUES (
    $1,
    $2,
    NOW()
)
RETURNING *;`

type InsertPlanStatusTimestampRow struct {
	PlanID    string             `json:"plan_id"`
	Status    string             `json:"status"`
	Timestamp pgtype.Timestamptz `json:"timestamp"`
}

// InsertPlanStatusTimestamp implements Querier.InsertPlanStatusTimestamp.
func (q *DBQuerier) InsertPlanStatusTimestamp(ctx context.Context, id string, status string) (InsertPlanStatusTimestampRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertPlanStatusTimestamp")
	row := q.conn.QueryRow(ctx, insertPlanStatusTimestampSQL, id, status)
	var item InsertPlanStatusTimestampRow
	if err := row.Scan(&item.PlanID, &item.Status, &item.Timestamp); err != nil {
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
	if err := row.Scan(&item.PlanID, &item.Status, &item.Timestamp); err != nil {
		return item, fmt.Errorf("scan InsertPlanStatusTimestampBatch row: %w", err)
	}
	return item, nil
}

const updatePlanStatusSQL = `UPDATE plans
SET
    status = $1,
    updated_at = NOW()
WHERE plan_id = $2
RETURNING *;`

type UpdatePlanStatusRow struct {
	PlanID               string             `json:"plan_id"`
	CreatedAt            pgtype.Timestamptz `json:"created_at"`
	UpdatedAt            pgtype.Timestamptz `json:"updated_at"`
	ResourceAdditions    int32              `json:"resource_additions"`
	ResourceChanges      int32              `json:"resource_changes"`
	ResourceDestructions int32              `json:"resource_destructions"`
	Status               string             `json:"status"`
	PlanFile             []byte             `json:"plan_file"`
	PlanJson             []byte             `json:"plan_json"`
	RunID                string             `json:"run_id"`
}

// UpdatePlanStatus implements Querier.UpdatePlanStatus.
func (q *DBQuerier) UpdatePlanStatus(ctx context.Context, status string, id string) (UpdatePlanStatusRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlanStatus")
	row := q.conn.QueryRow(ctx, updatePlanStatusSQL, status, id)
	var item UpdatePlanStatusRow
	if err := row.Scan(&item.PlanID, &item.CreatedAt, &item.UpdatedAt, &item.ResourceAdditions, &item.ResourceChanges, &item.ResourceDestructions, &item.Status, &item.PlanFile, &item.PlanJson, &item.RunID); err != nil {
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
	if err := row.Scan(&item.PlanID, &item.CreatedAt, &item.UpdatedAt, &item.ResourceAdditions, &item.ResourceChanges, &item.ResourceDestructions, &item.Status, &item.PlanFile, &item.PlanJson, &item.RunID); err != nil {
		return item, fmt.Errorf("scan UpdatePlanStatusBatch row: %w", err)
	}
	return item, nil
}

const getPlanFileByRunIDSQL = `SELECT plan_file
FROM plans
WHERE run_id = $1
;`

// GetPlanFileByRunID implements Querier.GetPlanFileByRunID.
func (q *DBQuerier) GetPlanFileByRunID(ctx context.Context, runID string) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetPlanFileByRunID")
	row := q.conn.QueryRow(ctx, getPlanFileByRunIDSQL, runID)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query GetPlanFileByRunID: %w", err)
	}
	return item, nil
}

// GetPlanFileByRunIDBatch implements Querier.GetPlanFileByRunIDBatch.
func (q *DBQuerier) GetPlanFileByRunIDBatch(batch genericBatch, runID string) {
	batch.Queue(getPlanFileByRunIDSQL, runID)
}

// GetPlanFileByRunIDScan implements Querier.GetPlanFileByRunIDScan.
func (q *DBQuerier) GetPlanFileByRunIDScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan GetPlanFileByRunIDBatch row: %w", err)
	}
	return item, nil
}

const getPlanJSONByRunIDSQL = `SELECT plan_json
FROM plans
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

const putPlanFileByRunIDSQL = `UPDATE plans
SET plan_file = $1
WHERE run_id = $2
;`

// PutPlanFileByRunID implements Querier.PutPlanFileByRunID.
func (q *DBQuerier) PutPlanFileByRunID(ctx context.Context, planFile []byte, runID string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "PutPlanFileByRunID")
	cmdTag, err := q.conn.Exec(ctx, putPlanFileByRunIDSQL, planFile, runID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query PutPlanFileByRunID: %w", err)
	}
	return cmdTag, err
}

// PutPlanFileByRunIDBatch implements Querier.PutPlanFileByRunIDBatch.
func (q *DBQuerier) PutPlanFileByRunIDBatch(batch genericBatch, planFile []byte, runID string) {
	batch.Queue(putPlanFileByRunIDSQL, planFile, runID)
}

// PutPlanFileByRunIDScan implements Querier.PutPlanFileByRunIDScan.
func (q *DBQuerier) PutPlanFileByRunIDScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec PutPlanFileByRunIDBatch: %w", err)
	}
	return cmdTag, err
}

const putPlanJSONByRunIDSQL = `UPDATE plans
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

const updatePlanResourcesSQL = `UPDATE plans
SET
    resource_additions = $1,
    resource_changes = $2,
    resource_destructions = $3
WHERE run_id = $4
;`

type UpdatePlanResourcesParams struct {
	ResourceAdditions    int32
	ResourceChanges      int32
	ResourceDestructions int32
	RunID                string
}

// UpdatePlanResources implements Querier.UpdatePlanResources.
func (q *DBQuerier) UpdatePlanResources(ctx context.Context, params UpdatePlanResourcesParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlanResources")
	cmdTag, err := q.conn.Exec(ctx, updatePlanResourcesSQL, params.ResourceAdditions, params.ResourceChanges, params.ResourceDestructions, params.RunID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query UpdatePlanResources: %w", err)
	}
	return cmdTag, err
}

// UpdatePlanResourcesBatch implements Querier.UpdatePlanResourcesBatch.
func (q *DBQuerier) UpdatePlanResourcesBatch(batch genericBatch, params UpdatePlanResourcesParams) {
	batch.Queue(updatePlanResourcesSQL, params.ResourceAdditions, params.ResourceChanges, params.ResourceDestructions, params.RunID)
}

// UpdatePlanResourcesScan implements Querier.UpdatePlanResourcesScan.
func (q *DBQuerier) UpdatePlanResourcesScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec UpdatePlanResourcesBatch: %w", err)
	}
	return cmdTag, err
}
