// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertTeamMembershipSQL = `WITH
    users AS (
        SELECT username
        FROM unnest($1::text[]) t(username)
    )
INSERT INTO team_memberships (username, team_id)
SELECT username, $2
FROM users
RETURNING username
;`

// InsertTeamMembership implements Querier.InsertTeamMembership.
func (q *DBQuerier) InsertTeamMembership(ctx context.Context, usernames []string, teamID pgtype.Text) ([]pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertTeamMembership")
	rows, err := q.conn.Query(ctx, insertTeamMembershipSQL, usernames, teamID)
	if err != nil {
		return nil, fmt.Errorf("query InsertTeamMembership: %w", err)
	}
	defer rows.Close()
	items := []pgtype.Text{}
	for rows.Next() {
		var item pgtype.Text
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan InsertTeamMembership row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close InsertTeamMembership rows: %w", err)
	}
	return items, err
}

// InsertTeamMembershipBatch implements Querier.InsertTeamMembershipBatch.
func (q *DBQuerier) InsertTeamMembershipBatch(batch genericBatch, usernames []string, teamID pgtype.Text) {
	batch.Queue(insertTeamMembershipSQL, usernames, teamID)
}

// InsertTeamMembershipScan implements Querier.InsertTeamMembershipScan.
func (q *DBQuerier) InsertTeamMembershipScan(results pgx.BatchResults) ([]pgtype.Text, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query InsertTeamMembershipBatch: %w", err)
	}
	defer rows.Close()
	items := []pgtype.Text{}
	for rows.Next() {
		var item pgtype.Text
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan InsertTeamMembershipBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close InsertTeamMembershipBatch rows: %w", err)
	}
	return items, err
}

const deleteTeamMembershipSQL = `WITH
    users AS (
        SELECT username
        FROM unnest($1::text[]) t(username)
    )
DELETE
FROM team_memberships tm
USING users
WHERE
    tm.username = users.username AND
    tm.team_id  = $2
RETURNING tm.username
;`

// DeleteTeamMembership implements Querier.DeleteTeamMembership.
func (q *DBQuerier) DeleteTeamMembership(ctx context.Context, usernames []string, teamID pgtype.Text) ([]pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteTeamMembership")
	rows, err := q.conn.Query(ctx, deleteTeamMembershipSQL, usernames, teamID)
	if err != nil {
		return nil, fmt.Errorf("query DeleteTeamMembership: %w", err)
	}
	defer rows.Close()
	items := []pgtype.Text{}
	for rows.Next() {
		var item pgtype.Text
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan DeleteTeamMembership row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close DeleteTeamMembership rows: %w", err)
	}
	return items, err
}

// DeleteTeamMembershipBatch implements Querier.DeleteTeamMembershipBatch.
func (q *DBQuerier) DeleteTeamMembershipBatch(batch genericBatch, usernames []string, teamID pgtype.Text) {
	batch.Queue(deleteTeamMembershipSQL, usernames, teamID)
}

// DeleteTeamMembershipScan implements Querier.DeleteTeamMembershipScan.
func (q *DBQuerier) DeleteTeamMembershipScan(results pgx.BatchResults) ([]pgtype.Text, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query DeleteTeamMembershipBatch: %w", err)
	}
	defer rows.Close()
	items := []pgtype.Text{}
	for rows.Next() {
		var item pgtype.Text
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan DeleteTeamMembershipBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close DeleteTeamMembershipBatch rows: %w", err)
	}
	return items, err
}