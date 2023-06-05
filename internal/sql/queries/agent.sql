-- name: InsertAgent :exec
INSERT INTO agents (
    agent_id         ,
    status           ,
    ip_address       ,
    version          ,
    name             ,
    external         ,
    last_seen        ,
    organization_name
) VALUES (
    pggen.arg('agent_id'),
    pggen.arg('status'),
    pggen.arg('ip_address'),
    pggen.arg('version'),
    pggen.arg('name'),
    pggen.arg('external'),
    pggen.arg('last_seen'),
    pggen.arg('organization_name')
);

-- name: UpdateAgentStatus :one
UPDATE agents
SET
    status = pggen.arg('status'),
    last_seen = pggen.arg('last_seen')
WHERE agent_id = pggen.arg('agent_id')
RETURNING status
;

-- name: FindAgentsByOrganization :many
SELECT *
FROM agents
WHERE organization_name = pggen.arg('organization_name')
;
