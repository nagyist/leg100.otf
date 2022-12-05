// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertVCSProviderSQL = `INSERT INTO vcs_providers (
    vcs_provider_id,
    token,
    created_at,
    name,
    cloud,
    organization_name
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);`

type InsertVCSProviderParams struct {
	VCSProviderID    pgtype.Text
	Token            pgtype.Text
	CreatedAt        pgtype.Timestamptz
	Name             pgtype.Text
	Cloud            pgtype.Text
	OrganizationName pgtype.Text
}

// InsertVCSProvider implements Querier.InsertVCSProvider.
func (q *DBQuerier) InsertVCSProvider(ctx context.Context, params InsertVCSProviderParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertVCSProvider")
	cmdTag, err := q.conn.Exec(ctx, insertVCSProviderSQL, params.VCSProviderID, params.Token, params.CreatedAt, params.Name, params.Cloud, params.OrganizationName)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertVCSProvider: %w", err)
	}
	return cmdTag, err
}

// InsertVCSProviderBatch implements Querier.InsertVCSProviderBatch.
func (q *DBQuerier) InsertVCSProviderBatch(batch genericBatch, params InsertVCSProviderParams) {
	batch.Queue(insertVCSProviderSQL, params.VCSProviderID, params.Token, params.CreatedAt, params.Name, params.Cloud, params.OrganizationName)
}

// InsertVCSProviderScan implements Querier.InsertVCSProviderScan.
func (q *DBQuerier) InsertVCSProviderScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertVCSProviderBatch: %w", err)
	}
	return cmdTag, err
}

const findVCSProvidersSQL = `SELECT *
FROM vcs_providers
WHERE organization_name = $1
;`

type FindVCSProvidersRow struct {
	VCSProviderID    pgtype.Text        `json:"vcs_provider_id"`
	Token            pgtype.Text        `json:"token"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	Name             pgtype.Text        `json:"name"`
	Cloud            pgtype.Text        `json:"cloud"`
	OrganizationName pgtype.Text        `json:"organization_name"`
}

// FindVCSProviders implements Querier.FindVCSProviders.
func (q *DBQuerier) FindVCSProviders(ctx context.Context, organizationName pgtype.Text) ([]FindVCSProvidersRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindVCSProviders")
	rows, err := q.conn.Query(ctx, findVCSProvidersSQL, organizationName)
	if err != nil {
		return nil, fmt.Errorf("query FindVCSProviders: %w", err)
	}
	defer rows.Close()
	items := []FindVCSProvidersRow{}
	for rows.Next() {
		var item FindVCSProvidersRow
		if err := rows.Scan(&item.VCSProviderID, &item.Token, &item.CreatedAt, &item.Name, &item.Cloud, &item.OrganizationName); err != nil {
			return nil, fmt.Errorf("scan FindVCSProviders row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindVCSProviders rows: %w", err)
	}
	return items, err
}

// FindVCSProvidersBatch implements Querier.FindVCSProvidersBatch.
func (q *DBQuerier) FindVCSProvidersBatch(batch genericBatch, organizationName pgtype.Text) {
	batch.Queue(findVCSProvidersSQL, organizationName)
}

// FindVCSProvidersScan implements Querier.FindVCSProvidersScan.
func (q *DBQuerier) FindVCSProvidersScan(results pgx.BatchResults) ([]FindVCSProvidersRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindVCSProvidersBatch: %w", err)
	}
	defer rows.Close()
	items := []FindVCSProvidersRow{}
	for rows.Next() {
		var item FindVCSProvidersRow
		if err := rows.Scan(&item.VCSProviderID, &item.Token, &item.CreatedAt, &item.Name, &item.Cloud, &item.OrganizationName); err != nil {
			return nil, fmt.Errorf("scan FindVCSProvidersBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindVCSProvidersBatch rows: %w", err)
	}
	return items, err
}

const findVCSProviderSQL = `SELECT *
FROM vcs_providers
WHERE vcs_provider_id = $1
;`

type FindVCSProviderRow struct {
	VCSProviderID    pgtype.Text        `json:"vcs_provider_id"`
	Token            pgtype.Text        `json:"token"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	Name             pgtype.Text        `json:"name"`
	Cloud            pgtype.Text        `json:"cloud"`
	OrganizationName pgtype.Text        `json:"organization_name"`
}

// FindVCSProvider implements Querier.FindVCSProvider.
func (q *DBQuerier) FindVCSProvider(ctx context.Context, vcsProviderID pgtype.Text) (FindVCSProviderRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindVCSProvider")
	row := q.conn.QueryRow(ctx, findVCSProviderSQL, vcsProviderID)
	var item FindVCSProviderRow
	if err := row.Scan(&item.VCSProviderID, &item.Token, &item.CreatedAt, &item.Name, &item.Cloud, &item.OrganizationName); err != nil {
		return item, fmt.Errorf("query FindVCSProvider: %w", err)
	}
	return item, nil
}

// FindVCSProviderBatch implements Querier.FindVCSProviderBatch.
func (q *DBQuerier) FindVCSProviderBatch(batch genericBatch, vcsProviderID pgtype.Text) {
	batch.Queue(findVCSProviderSQL, vcsProviderID)
}

// FindVCSProviderScan implements Querier.FindVCSProviderScan.
func (q *DBQuerier) FindVCSProviderScan(results pgx.BatchResults) (FindVCSProviderRow, error) {
	row := results.QueryRow()
	var item FindVCSProviderRow
	if err := row.Scan(&item.VCSProviderID, &item.Token, &item.CreatedAt, &item.Name, &item.Cloud, &item.OrganizationName); err != nil {
		return item, fmt.Errorf("scan FindVCSProviderBatch row: %w", err)
	}
	return item, nil
}

const deleteVCSProviderByIDSQL = `DELETE
FROM vcs_providers
WHERE vcs_provider_id = $1
RETURNING vcs_provider_id
;`

// DeleteVCSProviderByID implements Querier.DeleteVCSProviderByID.
func (q *DBQuerier) DeleteVCSProviderByID(ctx context.Context, vcsProviderID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteVCSProviderByID")
	row := q.conn.QueryRow(ctx, deleteVCSProviderByIDSQL, vcsProviderID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteVCSProviderByID: %w", err)
	}
	return item, nil
}

// DeleteVCSProviderByIDBatch implements Querier.DeleteVCSProviderByIDBatch.
func (q *DBQuerier) DeleteVCSProviderByIDBatch(batch genericBatch, vcsProviderID pgtype.Text) {
	batch.Queue(deleteVCSProviderByIDSQL, vcsProviderID)
}

// DeleteVCSProviderByIDScan implements Querier.DeleteVCSProviderByIDScan.
func (q *DBQuerier) DeleteVCSProviderByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteVCSProviderByIDBatch row: %w", err)
	}
	return item, nil
}
