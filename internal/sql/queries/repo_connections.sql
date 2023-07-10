-- name: InsertRepoConnection :exec
INSERT INTO repo_connections (
    webhook_id,
    workspace_id,
    module_id
) VALUES (
    pggen.arg('webhook_id'),
    pggen.arg('workspace_id'),
    pggen.arg('module_id')
);

-- name: DeleteWorkspaceConnectionByID :one
DELETE
FROM repo_connections
WHERE workspace_id = pggen.arg('workspace_id')
RETURNING webhook_id;

-- name: DeleteModuleConnectionByID :one
DELETE
FROM repo_connections
WHERE module_id = pggen.arg('module_id')
RETURNING webhook_id;