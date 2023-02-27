// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const upsertWorkspacePermissionSQL = `INSERT INTO workspace_permissions (
    workspace_id,
    team_id,
    role
) SELECT w.workspace_id, t.team_id, $1
    FROM teams t
    JOIN organizations o ON t.organization_name = o.name
    JOIN workspaces w ON w.organization_name = o.name
    WHERE t.name = $2
    AND w.workspace_id = $3
ON CONFLICT (workspace_id, team_id) DO UPDATE SET role = $1
;`

type UpsertWorkspacePermissionParams struct {
	Role        pgtype.Text
	TeamName    pgtype.Text
	WorkspaceID pgtype.Text
}

// UpsertWorkspacePermission implements Querier.UpsertWorkspacePermission.
func (q *DBQuerier) UpsertWorkspacePermission(ctx context.Context, params UpsertWorkspacePermissionParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpsertWorkspacePermission")
	cmdTag, err := q.conn.Exec(ctx, upsertWorkspacePermissionSQL, params.Role, params.TeamName, params.WorkspaceID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query UpsertWorkspacePermission: %w", err)
	}
	return cmdTag, err
}

// UpsertWorkspacePermissionBatch implements Querier.UpsertWorkspacePermissionBatch.
func (q *DBQuerier) UpsertWorkspacePermissionBatch(batch genericBatch, params UpsertWorkspacePermissionParams) {
	batch.Queue(upsertWorkspacePermissionSQL, params.Role, params.TeamName, params.WorkspaceID)
}

// UpsertWorkspacePermissionScan implements Querier.UpsertWorkspacePermissionScan.
func (q *DBQuerier) UpsertWorkspacePermissionScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec UpsertWorkspacePermissionBatch: %w", err)
	}
	return cmdTag, err
}

const findWorkspacePermissionsByIDSQL = `SELECT
    w.organization_name,
    w.workspace_id,
    (
        SELECT array_remove(array_agg(wp.*), NULL)
        FROM workspace_permissions wp
        WHERE wp.workspace_id = w.workspace_id
    ) AS workspace_permissions
FROM workspaces w
WHERE workspace_id = $1
;`

type FindWorkspacePermissionsByIDRow struct {
	OrganizationName     pgtype.Text            `json:"organization_name"`
	WorkspaceID          pgtype.Text            `json:"workspace_id"`
	WorkspacePermissions []WorkspacePermissions `json:"workspace_permissions"`
}

// FindWorkspacePermissionsByID implements Querier.FindWorkspacePermissionsByID.
func (q *DBQuerier) FindWorkspacePermissionsByID(ctx context.Context, workspaceID pgtype.Text) (FindWorkspacePermissionsByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindWorkspacePermissionsByID")
	row := q.conn.QueryRow(ctx, findWorkspacePermissionsByIDSQL, workspaceID)
	var item FindWorkspacePermissionsByIDRow
	workspacePermissionsArray := q.types.newWorkspacePermissionsArray()
	if err := row.Scan(&item.OrganizationName, &item.WorkspaceID, workspacePermissionsArray); err != nil {
		return item, fmt.Errorf("query FindWorkspacePermissionsByID: %w", err)
	}
	if err := workspacePermissionsArray.AssignTo(&item.WorkspacePermissions); err != nil {
		return item, fmt.Errorf("assign FindWorkspacePermissionsByID row: %w", err)
	}
	return item, nil
}

// FindWorkspacePermissionsByIDBatch implements Querier.FindWorkspacePermissionsByIDBatch.
func (q *DBQuerier) FindWorkspacePermissionsByIDBatch(batch genericBatch, workspaceID pgtype.Text) {
	batch.Queue(findWorkspacePermissionsByIDSQL, workspaceID)
}

// FindWorkspacePermissionsByIDScan implements Querier.FindWorkspacePermissionsByIDScan.
func (q *DBQuerier) FindWorkspacePermissionsByIDScan(results pgx.BatchResults) (FindWorkspacePermissionsByIDRow, error) {
	row := results.QueryRow()
	var item FindWorkspacePermissionsByIDRow
	workspacePermissionsArray := q.types.newWorkspacePermissionsArray()
	if err := row.Scan(&item.OrganizationName, &item.WorkspaceID, workspacePermissionsArray); err != nil {
		return item, fmt.Errorf("scan FindWorkspacePermissionsByIDBatch row: %w", err)
	}
	if err := workspacePermissionsArray.AssignTo(&item.WorkspacePermissions); err != nil {
		return item, fmt.Errorf("assign FindWorkspacePermissionsByID row: %w", err)
	}
	return item, nil
}

const deleteWorkspacePermissionByIDSQL = `DELETE
FROM workspace_permissions p
USING workspaces w, teams t
WHERE p.team_id = t.team_id
AND p.workspace_id = $1
AND t.name = $2
;`

// DeleteWorkspacePermissionByID implements Querier.DeleteWorkspacePermissionByID.
func (q *DBQuerier) DeleteWorkspacePermissionByID(ctx context.Context, workspaceID pgtype.Text, teamName pgtype.Text) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteWorkspacePermissionByID")
	cmdTag, err := q.conn.Exec(ctx, deleteWorkspacePermissionByIDSQL, workspaceID, teamName)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query DeleteWorkspacePermissionByID: %w", err)
	}
	return cmdTag, err
}

// DeleteWorkspacePermissionByIDBatch implements Querier.DeleteWorkspacePermissionByIDBatch.
func (q *DBQuerier) DeleteWorkspacePermissionByIDBatch(batch genericBatch, workspaceID pgtype.Text, teamName pgtype.Text) {
	batch.Queue(deleteWorkspacePermissionByIDSQL, workspaceID, teamName)
}

// DeleteWorkspacePermissionByIDScan implements Querier.DeleteWorkspacePermissionByIDScan.
func (q *DBQuerier) DeleteWorkspacePermissionByIDScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec DeleteWorkspacePermissionByIDBatch: %w", err)
	}
	return cmdTag, err
}
