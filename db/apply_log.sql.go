// Code generated by pggen. DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertApplyLogChunkSQL = `INSERT INTO apply_logs (
    apply_id,
    chunk,
    start,
    _end,
    size
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;`

type InsertApplyLogChunkParams struct {
	ApplyID string
	Chunk   []byte
	Start   bool
	End     bool
	Size    int32
}

type InsertApplyLogChunkRow struct {
	ApplyID string `json:"apply_id"`
	ChunkID int32  `json:"chunk_id"`
	Chunk   []byte `json:"chunk"`
	Size    int32  `json:"size"`
	Start   bool   `json:"start"`
	End     bool   `json:"_end"`
}

// InsertApplyLogChunk implements Querier.InsertApplyLogChunk.
func (q *DBQuerier) InsertApplyLogChunk(ctx context.Context, params InsertApplyLogChunkParams) (InsertApplyLogChunkRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertApplyLogChunk")
	row := q.conn.QueryRow(ctx, insertApplyLogChunkSQL, params.ApplyID, params.Chunk, params.Start, params.End, params.Size)
	var item InsertApplyLogChunkRow
	if err := row.Scan(&item.ApplyID, &item.ChunkID, &item.Chunk, &item.Size, &item.Start, &item.End); err != nil {
		return item, fmt.Errorf("query InsertApplyLogChunk: %w", err)
	}
	return item, nil
}

// InsertApplyLogChunkBatch implements Querier.InsertApplyLogChunkBatch.
func (q *DBQuerier) InsertApplyLogChunkBatch(batch genericBatch, params InsertApplyLogChunkParams) {
	batch.Queue(insertApplyLogChunkSQL, params.ApplyID, params.Chunk, params.Start, params.End, params.Size)
}

// InsertApplyLogChunkScan implements Querier.InsertApplyLogChunkScan.
func (q *DBQuerier) InsertApplyLogChunkScan(results pgx.BatchResults) (InsertApplyLogChunkRow, error) {
	row := results.QueryRow()
	var item InsertApplyLogChunkRow
	if err := row.Scan(&item.ApplyID, &item.ChunkID, &item.Chunk, &item.Size, &item.Start, &item.End); err != nil {
		return item, fmt.Errorf("scan InsertApplyLogChunkBatch row: %w", err)
	}
	return item, nil
}

const findApplyLogChunksSQL = `SELECT chunk, start, _end
FROM apply_logs
WHERE apply_id = $1
ORDER BY chunk_id ASC
;`

type FindApplyLogChunksRow struct {
	Chunk pgtype.Bytea `json:"chunk"`
	Start *bool        `json:"start"`
	End   *bool        `json:"_end"`
}

// FindApplyLogChunks implements Querier.FindApplyLogChunks.
func (q *DBQuerier) FindApplyLogChunks(ctx context.Context, applyID string) ([]FindApplyLogChunksRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindApplyLogChunks")
	rows, err := q.conn.Query(ctx, findApplyLogChunksSQL, applyID)
	if err != nil {
		return nil, fmt.Errorf("query FindApplyLogChunks: %w", err)
	}
	defer rows.Close()
	items := []FindApplyLogChunksRow{}
	for rows.Next() {
		var item FindApplyLogChunksRow
		if err := rows.Scan(&item.Chunk, &item.Start, &item.End); err != nil {
			return nil, fmt.Errorf("scan FindApplyLogChunks row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindApplyLogChunks rows: %w", err)
	}
	return items, err
}

// FindApplyLogChunksBatch implements Querier.FindApplyLogChunksBatch.
func (q *DBQuerier) FindApplyLogChunksBatch(batch genericBatch, applyID string) {
	batch.Queue(findApplyLogChunksSQL, applyID)
}

// FindApplyLogChunksScan implements Querier.FindApplyLogChunksScan.
func (q *DBQuerier) FindApplyLogChunksScan(results pgx.BatchResults) ([]FindApplyLogChunksRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindApplyLogChunksBatch: %w", err)
	}
	defer rows.Close()
	items := []FindApplyLogChunksRow{}
	for rows.Next() {
		var item FindApplyLogChunksRow
		if err := rows.Scan(&item.Chunk, &item.Start, &item.End); err != nil {
			return nil, fmt.Errorf("scan FindApplyLogChunksBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindApplyLogChunksBatch rows: %w", err)
	}
	return items, err
}
