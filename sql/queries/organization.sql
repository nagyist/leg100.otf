-- name: FindOrganizationNameByWorkspaceID :one
SELECT organization_name
FROM workspaces
WHERE workspace_id = pggen.arg('workspace_id')
;

-- name: FindOrganizationByName :one
SELECT * FROM organizations WHERE name = pggen.arg('name');

-- name: FindOrganizationByID :one
SELECT * FROM organizations WHERE organization_id = pggen.arg('organization_id');

-- name: FindOrganizationByNameForUpdate :one
SELECT *
FROM organizations
WHERE name = pggen.arg('name')
FOR UPDATE
;

-- name: FindOrganizations :many
SELECT *
FROM organizations
ORDER BY updated_at DESC
LIMIT pggen.arg('limit') OFFSET pggen.arg('offset');

-- name: CountOrganizations :one
SELECT count(*)
FROM organizations;

-- name: FindOrganizationsByUserID :many
SELECT o.*
FROM organizations o
JOIN organization_memberships om ON o.name = om.organization_name
WHERE om.user_id = pggen.arg('user_id')
;

-- name: InsertOrganization :exec
INSERT INTO organizations (
    organization_id,
    created_at,
    updated_at,
    name,
    session_remember,
    session_timeout
) VALUES (
    pggen.arg('id'),
    pggen.arg('created_at'),
    pggen.arg('updated_at'),
    pggen.arg('name'),
    pggen.arg('session_remember'),
    pggen.arg('session_timeout')
);

-- name: UpdateOrganizationByName :one
UPDATE organizations
SET
    name = pggen.arg('new_name'),
    session_remember = pggen.arg('session_remember'),
    session_timeout = pggen.arg('session_timeout'),
    updated_at = pggen.arg('updated_at')
WHERE name = pggen.arg('name')
RETURNING organization_id;

-- name: DeleteOrganization :one
DELETE
FROM organizations
WHERE name = pggen.arg('name')
RETURNING organization_id;
