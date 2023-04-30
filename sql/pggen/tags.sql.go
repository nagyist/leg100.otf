// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertTagSQL = `INSERT INTO tags (
    tag_id,
    name,
    organization_name
) SELECT $1, $2, w.organization_name
  FROM workspaces w
  WHERE w.workspace_id = $3
ON CONFLICT (name, organization_name) DO NOTHING
;`

type InsertTagParams struct {
	TagID       pgtype.Text
	Name        pgtype.Text
	WorkspaceID pgtype.Text
}

// InsertTag implements Querier.InsertTag.
func (q *DBQuerier) InsertTag(ctx context.Context, params InsertTagParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertTag")
	cmdTag, err := q.conn.Exec(ctx, insertTagSQL, params.TagID, params.Name, params.WorkspaceID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertTag: %w", err)
	}
	return cmdTag, err
}

// InsertTagBatch implements Querier.InsertTagBatch.
func (q *DBQuerier) InsertTagBatch(batch genericBatch, params InsertTagParams) {
	batch.Queue(insertTagSQL, params.TagID, params.Name, params.WorkspaceID)
}

// InsertTagScan implements Querier.InsertTagScan.
func (q *DBQuerier) InsertTagScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertTagBatch: %w", err)
	}
	return cmdTag, err
}

const insertWorkspaceTagSQL = `INSERT INTO workspace_tags (
    tag_id,
    workspace_id
) SELECT $1, $2
  FROM workspaces w
  JOIN tags t ON (t.organization_name = w.organization_name)
  WHERE w.workspace_id = $2
  AND t.tag_id = $1
RETURNING tag_id
;`

// InsertWorkspaceTag implements Querier.InsertWorkspaceTag.
func (q *DBQuerier) InsertWorkspaceTag(ctx context.Context, tagID pgtype.Text, workspaceID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertWorkspaceTag")
	row := q.conn.QueryRow(ctx, insertWorkspaceTagSQL, tagID, workspaceID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query InsertWorkspaceTag: %w", err)
	}
	return item, nil
}

// InsertWorkspaceTagBatch implements Querier.InsertWorkspaceTagBatch.
func (q *DBQuerier) InsertWorkspaceTagBatch(batch genericBatch, tagID pgtype.Text, workspaceID pgtype.Text) {
	batch.Queue(insertWorkspaceTagSQL, tagID, workspaceID)
}

// InsertWorkspaceTagScan implements Querier.InsertWorkspaceTagScan.
func (q *DBQuerier) InsertWorkspaceTagScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan InsertWorkspaceTagBatch row: %w", err)
	}
	return item, nil
}

const insertWorkspaceTagByNameSQL = `INSERT INTO workspace_tags (
    tag_id,
    workspace_id
) SELECT t.tag_id, $1
  FROM workspaces w
  JOIN tags t ON (t.organization_name = w.organization_name)
  WHERE t.name = $2
RETURNING tag_id
;`

// InsertWorkspaceTagByName implements Querier.InsertWorkspaceTagByName.
func (q *DBQuerier) InsertWorkspaceTagByName(ctx context.Context, workspaceID pgtype.Text, tagName pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertWorkspaceTagByName")
	row := q.conn.QueryRow(ctx, insertWorkspaceTagByNameSQL, workspaceID, tagName)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query InsertWorkspaceTagByName: %w", err)
	}
	return item, nil
}

// InsertWorkspaceTagByNameBatch implements Querier.InsertWorkspaceTagByNameBatch.
func (q *DBQuerier) InsertWorkspaceTagByNameBatch(batch genericBatch, workspaceID pgtype.Text, tagName pgtype.Text) {
	batch.Queue(insertWorkspaceTagByNameSQL, workspaceID, tagName)
}

// InsertWorkspaceTagByNameScan implements Querier.InsertWorkspaceTagByNameScan.
func (q *DBQuerier) InsertWorkspaceTagByNameScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan InsertWorkspaceTagByNameBatch row: %w", err)
	}
	return item, nil
}

const findTagsSQL = `SELECT
    t.*,
    (
        SELECT count(*)
        FROM workspace_tags wt
        WHERE wt.tag_id = t.tag_id
    ) AS instance_count
FROM tags t
WHERE t.organization_name = $1
LIMIT $2
OFFSET $3
;`

type FindTagsParams struct {
	OrganizationName pgtype.Text
	Limit            int
	Offset           int
}

type FindTagsRow struct {
	TagID            pgtype.Text `json:"tag_id"`
	Name             pgtype.Text `json:"name"`
	OrganizationName pgtype.Text `json:"organization_name"`
	InstanceCount    int         `json:"instance_count"`
}

// FindTags implements Querier.FindTags.
func (q *DBQuerier) FindTags(ctx context.Context, params FindTagsParams) ([]FindTagsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindTags")
	rows, err := q.conn.Query(ctx, findTagsSQL, params.OrganizationName, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindTags: %w", err)
	}
	defer rows.Close()
	items := []FindTagsRow{}
	for rows.Next() {
		var item FindTagsRow
		if err := rows.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
			return nil, fmt.Errorf("scan FindTags row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindTags rows: %w", err)
	}
	return items, err
}

// FindTagsBatch implements Querier.FindTagsBatch.
func (q *DBQuerier) FindTagsBatch(batch genericBatch, params FindTagsParams) {
	batch.Queue(findTagsSQL, params.OrganizationName, params.Limit, params.Offset)
}

// FindTagsScan implements Querier.FindTagsScan.
func (q *DBQuerier) FindTagsScan(results pgx.BatchResults) ([]FindTagsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindTagsBatch: %w", err)
	}
	defer rows.Close()
	items := []FindTagsRow{}
	for rows.Next() {
		var item FindTagsRow
		if err := rows.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
			return nil, fmt.Errorf("scan FindTagsBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindTagsBatch rows: %w", err)
	}
	return items, err
}

const findWorkspaceTagsSQL = `SELECT
    t.*,
    (
        SELECT count(*)
        FROM workspace_tags wt
        WHERE wt.tag_id = t.tag_id
    ) AS instance_count
FROM workspace_tags wt
JOIN tags t USING (tag_id)
WHERE wt.workspace_id = $1
LIMIT $2
OFFSET $3
;`

type FindWorkspaceTagsParams struct {
	WorkspaceID pgtype.Text
	Limit       int
	Offset      int
}

type FindWorkspaceTagsRow struct {
	TagID            pgtype.Text `json:"tag_id"`
	Name             pgtype.Text `json:"name"`
	OrganizationName pgtype.Text `json:"organization_name"`
	InstanceCount    int         `json:"instance_count"`
}

// FindWorkspaceTags implements Querier.FindWorkspaceTags.
func (q *DBQuerier) FindWorkspaceTags(ctx context.Context, params FindWorkspaceTagsParams) ([]FindWorkspaceTagsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindWorkspaceTags")
	rows, err := q.conn.Query(ctx, findWorkspaceTagsSQL, params.WorkspaceID, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindWorkspaceTags: %w", err)
	}
	defer rows.Close()
	items := []FindWorkspaceTagsRow{}
	for rows.Next() {
		var item FindWorkspaceTagsRow
		if err := rows.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
			return nil, fmt.Errorf("scan FindWorkspaceTags row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindWorkspaceTags rows: %w", err)
	}
	return items, err
}

// FindWorkspaceTagsBatch implements Querier.FindWorkspaceTagsBatch.
func (q *DBQuerier) FindWorkspaceTagsBatch(batch genericBatch, params FindWorkspaceTagsParams) {
	batch.Queue(findWorkspaceTagsSQL, params.WorkspaceID, params.Limit, params.Offset)
}

// FindWorkspaceTagsScan implements Querier.FindWorkspaceTagsScan.
func (q *DBQuerier) FindWorkspaceTagsScan(results pgx.BatchResults) ([]FindWorkspaceTagsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindWorkspaceTagsBatch: %w", err)
	}
	defer rows.Close()
	items := []FindWorkspaceTagsRow{}
	for rows.Next() {
		var item FindWorkspaceTagsRow
		if err := rows.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
			return nil, fmt.Errorf("scan FindWorkspaceTagsBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindWorkspaceTagsBatch rows: %w", err)
	}
	return items, err
}

const findTagByNameSQL = `SELECT
    t.*,
    (
        SELECT count(*)
        FROM workspace_tags wt
        WHERE wt.tag_id = t.tag_id
    ) AS instance_count
FROM tags t
JOIN workspace_tags wt USING (tag_id)
WHERE t.name = $1
AND   wt.workspace_id = $2
;`

type FindTagByNameRow struct {
	TagID            pgtype.Text `json:"tag_id"`
	Name             pgtype.Text `json:"name"`
	OrganizationName pgtype.Text `json:"organization_name"`
	InstanceCount    int         `json:"instance_count"`
}

// FindTagByName implements Querier.FindTagByName.
func (q *DBQuerier) FindTagByName(ctx context.Context, name pgtype.Text, workspaceID pgtype.Text) (FindTagByNameRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindTagByName")
	row := q.conn.QueryRow(ctx, findTagByNameSQL, name, workspaceID)
	var item FindTagByNameRow
	if err := row.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
		return item, fmt.Errorf("query FindTagByName: %w", err)
	}
	return item, nil
}

// FindTagByNameBatch implements Querier.FindTagByNameBatch.
func (q *DBQuerier) FindTagByNameBatch(batch genericBatch, name pgtype.Text, workspaceID pgtype.Text) {
	batch.Queue(findTagByNameSQL, name, workspaceID)
}

// FindTagByNameScan implements Querier.FindTagByNameScan.
func (q *DBQuerier) FindTagByNameScan(results pgx.BatchResults) (FindTagByNameRow, error) {
	row := results.QueryRow()
	var item FindTagByNameRow
	if err := row.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
		return item, fmt.Errorf("scan FindTagByNameBatch row: %w", err)
	}
	return item, nil
}

const findTagByIDSQL = `SELECT
    t.*,
    (
        SELECT count(*)
        FROM workspace_tags wt
        WHERE wt.tag_id = t.tag_id
    ) AS instance_count
FROM tags t
JOIN workspace_tags wt USING (tag_id)
WHERE t.tag_id = $1
AND   wt.workspace_id = $2
;`

type FindTagByIDRow struct {
	TagID            pgtype.Text `json:"tag_id"`
	Name             pgtype.Text `json:"name"`
	OrganizationName pgtype.Text `json:"organization_name"`
	InstanceCount    int         `json:"instance_count"`
}

// FindTagByID implements Querier.FindTagByID.
func (q *DBQuerier) FindTagByID(ctx context.Context, tagID pgtype.Text, workspaceID pgtype.Text) (FindTagByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindTagByID")
	row := q.conn.QueryRow(ctx, findTagByIDSQL, tagID, workspaceID)
	var item FindTagByIDRow
	if err := row.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
		return item, fmt.Errorf("query FindTagByID: %w", err)
	}
	return item, nil
}

// FindTagByIDBatch implements Querier.FindTagByIDBatch.
func (q *DBQuerier) FindTagByIDBatch(batch genericBatch, tagID pgtype.Text, workspaceID pgtype.Text) {
	batch.Queue(findTagByIDSQL, tagID, workspaceID)
}

// FindTagByIDScan implements Querier.FindTagByIDScan.
func (q *DBQuerier) FindTagByIDScan(results pgx.BatchResults) (FindTagByIDRow, error) {
	row := results.QueryRow()
	var item FindTagByIDRow
	if err := row.Scan(&item.TagID, &item.Name, &item.OrganizationName, &item.InstanceCount); err != nil {
		return item, fmt.Errorf("scan FindTagByIDBatch row: %w", err)
	}
	return item, nil
}

const countTagsSQL = `SELECT count(*)
FROM tags t
WHERE t.organization_name = $1
;`

// CountTags implements Querier.CountTags.
func (q *DBQuerier) CountTags(ctx context.Context, organizationName pgtype.Text) (int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountTags")
	row := q.conn.QueryRow(ctx, countTagsSQL, organizationName)
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query CountTags: %w", err)
	}
	return item, nil
}

// CountTagsBatch implements Querier.CountTagsBatch.
func (q *DBQuerier) CountTagsBatch(batch genericBatch, organizationName pgtype.Text) {
	batch.Queue(countTagsSQL, organizationName)
}

// CountTagsScan implements Querier.CountTagsScan.
func (q *DBQuerier) CountTagsScan(results pgx.BatchResults) (int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan CountTagsBatch row: %w", err)
	}
	return item, nil
}

const countWorkspaceTagsSQL = `SELECT count(*)
FROM workspace_tags wt
WHERE wt.workspace_id = $1
;`

// CountWorkspaceTags implements Querier.CountWorkspaceTags.
func (q *DBQuerier) CountWorkspaceTags(ctx context.Context, workspaceID pgtype.Text) (int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountWorkspaceTags")
	row := q.conn.QueryRow(ctx, countWorkspaceTagsSQL, workspaceID)
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query CountWorkspaceTags: %w", err)
	}
	return item, nil
}

// CountWorkspaceTagsBatch implements Querier.CountWorkspaceTagsBatch.
func (q *DBQuerier) CountWorkspaceTagsBatch(batch genericBatch, workspaceID pgtype.Text) {
	batch.Queue(countWorkspaceTagsSQL, workspaceID)
}

// CountWorkspaceTagsScan implements Querier.CountWorkspaceTagsScan.
func (q *DBQuerier) CountWorkspaceTagsScan(results pgx.BatchResults) (int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan CountWorkspaceTagsBatch row: %w", err)
	}
	return item, nil
}

const deleteTagSQL = `DELETE
FROM tags
WHERE tag_id            = $1
AND   organization_name = $2
RETURNING tag_id
;`

// DeleteTag implements Querier.DeleteTag.
func (q *DBQuerier) DeleteTag(ctx context.Context, tagID pgtype.Text, organizationName pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteTag")
	row := q.conn.QueryRow(ctx, deleteTagSQL, tagID, organizationName)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteTag: %w", err)
	}
	return item, nil
}

// DeleteTagBatch implements Querier.DeleteTagBatch.
func (q *DBQuerier) DeleteTagBatch(batch genericBatch, tagID pgtype.Text, organizationName pgtype.Text) {
	batch.Queue(deleteTagSQL, tagID, organizationName)
}

// DeleteTagScan implements Querier.DeleteTagScan.
func (q *DBQuerier) DeleteTagScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteTagBatch row: %w", err)
	}
	return item, nil
}

const deleteWorkspaceTagSQL = `DELETE
FROM workspace_tags
WHERE workspace_id  = $1
AND   tag_id        = $2
RETURNING tag_id
;`

// DeleteWorkspaceTag implements Querier.DeleteWorkspaceTag.
func (q *DBQuerier) DeleteWorkspaceTag(ctx context.Context, workspaceID pgtype.Text, tagID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteWorkspaceTag")
	row := q.conn.QueryRow(ctx, deleteWorkspaceTagSQL, workspaceID, tagID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteWorkspaceTag: %w", err)
	}
	return item, nil
}

// DeleteWorkspaceTagBatch implements Querier.DeleteWorkspaceTagBatch.
func (q *DBQuerier) DeleteWorkspaceTagBatch(batch genericBatch, workspaceID pgtype.Text, tagID pgtype.Text) {
	batch.Queue(deleteWorkspaceTagSQL, workspaceID, tagID)
}

// DeleteWorkspaceTagScan implements Querier.DeleteWorkspaceTagScan.
func (q *DBQuerier) DeleteWorkspaceTagScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteWorkspaceTagBatch row: %w", err)
	}
	return item, nil
}
