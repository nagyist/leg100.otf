// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertAgentTokenSQL = `INSERT INTO agent_tokens (
    agent_token_id,
    created_at,
    description,
    organization_name
) VALUES (
    $1,
    $2,
    $3,
    $4
);`

type InsertAgentTokenParams struct {
	AgentTokenID     pgtype.Text
	CreatedAt        pgtype.Timestamptz
	Description      pgtype.Text
	OrganizationName pgtype.Text
}

// InsertAgentToken implements Querier.InsertAgentToken.
func (q *DBQuerier) InsertAgentToken(ctx context.Context, params InsertAgentTokenParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertAgentToken")
	cmdTag, err := q.conn.Exec(ctx, insertAgentTokenSQL, params.AgentTokenID, params.CreatedAt, params.Description, params.OrganizationName)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertAgentToken: %w", err)
	}
	return cmdTag, err
}

// InsertAgentTokenBatch implements Querier.InsertAgentTokenBatch.
func (q *DBQuerier) InsertAgentTokenBatch(batch genericBatch, params InsertAgentTokenParams) {
	batch.Queue(insertAgentTokenSQL, params.AgentTokenID, params.CreatedAt, params.Description, params.OrganizationName)
}

// InsertAgentTokenScan implements Querier.InsertAgentTokenScan.
func (q *DBQuerier) InsertAgentTokenScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertAgentTokenBatch: %w", err)
	}
	return cmdTag, err
}

const findAgentTokenByIDSQL = `SELECT *
FROM agent_tokens
WHERE agent_token_id = $1
;`

type FindAgentTokenByIDRow struct {
	AgentTokenID     pgtype.Text        `json:"agent_token_id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	Description      pgtype.Text        `json:"description"`
	OrganizationName pgtype.Text        `json:"organization_name"`
	AgentPoolID      pgtype.Text        `json:"agent_pool_id"`
}

// FindAgentTokenByID implements Querier.FindAgentTokenByID.
func (q *DBQuerier) FindAgentTokenByID(ctx context.Context, agentTokenID pgtype.Text) (FindAgentTokenByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAgentTokenByID")
	row := q.conn.QueryRow(ctx, findAgentTokenByIDSQL, agentTokenID)
	var item FindAgentTokenByIDRow
	if err := row.Scan(&item.AgentTokenID, &item.CreatedAt, &item.Description, &item.OrganizationName, &item.AgentPoolID); err != nil {
		return item, fmt.Errorf("query FindAgentTokenByID: %w", err)
	}
	return item, nil
}

// FindAgentTokenByIDBatch implements Querier.FindAgentTokenByIDBatch.
func (q *DBQuerier) FindAgentTokenByIDBatch(batch genericBatch, agentTokenID pgtype.Text) {
	batch.Queue(findAgentTokenByIDSQL, agentTokenID)
}

// FindAgentTokenByIDScan implements Querier.FindAgentTokenByIDScan.
func (q *DBQuerier) FindAgentTokenByIDScan(results pgx.BatchResults) (FindAgentTokenByIDRow, error) {
	row := results.QueryRow()
	var item FindAgentTokenByIDRow
	if err := row.Scan(&item.AgentTokenID, &item.CreatedAt, &item.Description, &item.OrganizationName, &item.AgentPoolID); err != nil {
		return item, fmt.Errorf("scan FindAgentTokenByIDBatch row: %w", err)
	}
	return item, nil
}

const findAgentTokensSQL = `SELECT *
FROM agent_tokens
WHERE organization_name = $1
ORDER BY created_at DESC
;`

type FindAgentTokensRow struct {
	AgentTokenID     pgtype.Text        `json:"agent_token_id"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	Description      pgtype.Text        `json:"description"`
	OrganizationName pgtype.Text        `json:"organization_name"`
	AgentPoolID      pgtype.Text        `json:"agent_pool_id"`
}

// FindAgentTokens implements Querier.FindAgentTokens.
func (q *DBQuerier) FindAgentTokens(ctx context.Context, organizationName pgtype.Text) ([]FindAgentTokensRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindAgentTokens")
	rows, err := q.conn.Query(ctx, findAgentTokensSQL, organizationName)
	if err != nil {
		return nil, fmt.Errorf("query FindAgentTokens: %w", err)
	}
	defer rows.Close()
	items := []FindAgentTokensRow{}
	for rows.Next() {
		var item FindAgentTokensRow
		if err := rows.Scan(&item.AgentTokenID, &item.CreatedAt, &item.Description, &item.OrganizationName, &item.AgentPoolID); err != nil {
			return nil, fmt.Errorf("scan FindAgentTokens row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAgentTokens rows: %w", err)
	}
	return items, err
}

// FindAgentTokensBatch implements Querier.FindAgentTokensBatch.
func (q *DBQuerier) FindAgentTokensBatch(batch genericBatch, organizationName pgtype.Text) {
	batch.Queue(findAgentTokensSQL, organizationName)
}

// FindAgentTokensScan implements Querier.FindAgentTokensScan.
func (q *DBQuerier) FindAgentTokensScan(results pgx.BatchResults) ([]FindAgentTokensRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindAgentTokensBatch: %w", err)
	}
	defer rows.Close()
	items := []FindAgentTokensRow{}
	for rows.Next() {
		var item FindAgentTokensRow
		if err := rows.Scan(&item.AgentTokenID, &item.CreatedAt, &item.Description, &item.OrganizationName, &item.AgentPoolID); err != nil {
			return nil, fmt.Errorf("scan FindAgentTokensBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAgentTokensBatch rows: %w", err)
	}
	return items, err
}

const deleteAgentTokenByIDSQL = `DELETE
FROM agent_tokens
WHERE agent_token_id = $1
RETURNING agent_token_id
;`

// DeleteAgentTokenByID implements Querier.DeleteAgentTokenByID.
func (q *DBQuerier) DeleteAgentTokenByID(ctx context.Context, agentTokenID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteAgentTokenByID")
	row := q.conn.QueryRow(ctx, deleteAgentTokenByIDSQL, agentTokenID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteAgentTokenByID: %w", err)
	}
	return item, nil
}

// DeleteAgentTokenByIDBatch implements Querier.DeleteAgentTokenByIDBatch.
func (q *DBQuerier) DeleteAgentTokenByIDBatch(batch genericBatch, agentTokenID pgtype.Text) {
	batch.Queue(deleteAgentTokenByIDSQL, agentTokenID)
}

// DeleteAgentTokenByIDScan implements Querier.DeleteAgentTokenByIDScan.
func (q *DBQuerier) DeleteAgentTokenByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteAgentTokenByIDBatch row: %w", err)
	}
	return item, nil
}
