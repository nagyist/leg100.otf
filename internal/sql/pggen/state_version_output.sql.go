// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertStateVersionOutputSQL = `INSERT INTO state_version_outputs (
    state_version_output_id,
    name,
    sensitive,
    type,
    value,
    state_version_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);`

type InsertStateVersionOutputParams struct {
	ID             pgtype.Text
	Name           pgtype.Text
	Sensitive      bool
	Type           pgtype.Text
	Value          []byte
	StateVersionID pgtype.Text
}

// InsertStateVersionOutput implements Querier.InsertStateVersionOutput.
func (q *DBQuerier) InsertStateVersionOutput(ctx context.Context, params InsertStateVersionOutputParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertStateVersionOutput")
	cmdTag, err := q.conn.Exec(ctx, insertStateVersionOutputSQL, params.ID, params.Name, params.Sensitive, params.Type, params.Value, params.StateVersionID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertStateVersionOutput: %w", err)
	}
	return cmdTag, err
}

// InsertStateVersionOutputBatch implements Querier.InsertStateVersionOutputBatch.
func (q *DBQuerier) InsertStateVersionOutputBatch(batch genericBatch, params InsertStateVersionOutputParams) {
	batch.Queue(insertStateVersionOutputSQL, params.ID, params.Name, params.Sensitive, params.Type, params.Value, params.StateVersionID)
}

// InsertStateVersionOutputScan implements Querier.InsertStateVersionOutputScan.
func (q *DBQuerier) InsertStateVersionOutputScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertStateVersionOutputBatch: %w", err)
	}
	return cmdTag, err
}

const findStateVersionOutputByIDSQL = `SELECT *
FROM state_version_outputs
WHERE state_version_output_id = $1
;`

type FindStateVersionOutputByIDRow struct {
	StateVersionOutputID pgtype.Text `json:"state_version_output_id"`
	Name                 pgtype.Text `json:"name"`
	Sensitive            bool        `json:"sensitive"`
	Type                 pgtype.Text `json:"type"`
	Value                []byte      `json:"value"`
	StateVersionID       pgtype.Text `json:"state_version_id"`
}

// FindStateVersionOutputByID implements Querier.FindStateVersionOutputByID.
func (q *DBQuerier) FindStateVersionOutputByID(ctx context.Context, id pgtype.Text) (FindStateVersionOutputByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindStateVersionOutputByID")
	row := q.conn.QueryRow(ctx, findStateVersionOutputByIDSQL, id)
	var item FindStateVersionOutputByIDRow
	if err := row.Scan(&item.StateVersionOutputID, &item.Name, &item.Sensitive, &item.Type, &item.Value, &item.StateVersionID); err != nil {
		return item, fmt.Errorf("query FindStateVersionOutputByID: %w", err)
	}
	return item, nil
}

// FindStateVersionOutputByIDBatch implements Querier.FindStateVersionOutputByIDBatch.
func (q *DBQuerier) FindStateVersionOutputByIDBatch(batch genericBatch, id pgtype.Text) {
	batch.Queue(findStateVersionOutputByIDSQL, id)
}

// FindStateVersionOutputByIDScan implements Querier.FindStateVersionOutputByIDScan.
func (q *DBQuerier) FindStateVersionOutputByIDScan(results pgx.BatchResults) (FindStateVersionOutputByIDRow, error) {
	row := results.QueryRow()
	var item FindStateVersionOutputByIDRow
	if err := row.Scan(&item.StateVersionOutputID, &item.Name, &item.Sensitive, &item.Type, &item.Value, &item.StateVersionID); err != nil {
		return item, fmt.Errorf("scan FindStateVersionOutputByIDBatch row: %w", err)
	}
	return item, nil
}