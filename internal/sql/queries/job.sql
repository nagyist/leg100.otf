-- name: InsertJob :exec
INSERT INTO jobs (
    job_id,
    run_id,
    phase,
    agent_id,
    status
) VALUES (
    pggen.arg('job_id'),
    pggen.arg('run_id'),
    pggen.arg('phase'),
    pggen.arg('agent_id'),
    pggen.arg('status')
);

-- name: AssignJob :one
UPDATE jobs
SET
    status = pggen.arg('status'),
    agent_id = pggen.arg('agent_id')
WHERE job_id = pggen.arg('job_id')
RETURNING status
;

-- name: UpdateJobStatus :one
UPDATE jobs
SET
    status = pggen.arg('status')
WHERE job_id = pggen.arg('job_id')
RETURNING status
;

-- name: FindJobsByStatus :many
SELECT *
FROM jobs
WHERE status = pggen.arg('status')
;
