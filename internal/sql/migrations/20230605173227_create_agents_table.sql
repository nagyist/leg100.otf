-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS agent_statuses (
    status TEXT PRIMARY KEY
);

INSERT INTO agent_statuses (status) VALUES
	('busy'),
	('idle'),
	('unknown'),
	('errored'),
	('exited');

CREATE TABLE IF NOT EXISTS agents (
    agent_id          TEXT,
    status            TEXT REFERENCES agent_statuses NOT NULL,
    ip_address        TEXT NOT NULL,
    version           TEXT NOT NULL,
    name              TEXT,
    external          BOOLEAN NOT NULL,
    last_seen         TIMESTAMPTZ NOT NULL,
    organization_name TEXT REFERENCES organizations (name) ON UPDATE CASCADE ON DELETE CASCADE,
                      PRIMARY KEY (agent_id)
);

CREATE OR REPLACE FUNCTION agents_notify_event() RETURNS TRIGGER AS $$
DECLARE
    record RECORD;
    notification JSON;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        record = OLD;
    ELSE
        record = NEW;
    END IF;
    notification = json_build_object(
                      'table',TG_TABLE_NAME,
                      'action', TG_OP,
                      'id', record.agent_id);
    PERFORM pg_notify('events', notification::text);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_event
AFTER INSERT OR UPDATE OR DELETE ON agents
    FOR EACH ROW EXECUTE PROCEDURE agents_notify_event();

CREATE TABLE IF NOT EXISTS job_statuses (
    status TEXT PRIMARY KEY
);

INSERT INTO job_statuses (status) VALUES
	('created'),
	('assigned'),
	('running'),
	('completed'),
	('errored');

CREATE TABLE IF NOT EXISTS jobs (
    job_id      TEXT NOT NULL,
    run_id      TEXT REFERENCES runs ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    phase       TEXT REFERENCES phases ON UPDATE CASCADE NOT NULL,
    agent_id    TEXT REFERENCES agents ON UPDATE CASCADE ON DELETE CASCADE,
    status      TEXT REFERENCES job_statuses ON UPDATE CASCADE NOT NULL,
                PRIMARY KEY (job_id),
                UNIQUE (run_id, phase)
);

CREATE OR REPLACE FUNCTION jobs_notify_event() RETURNS TRIGGER AS $$
DECLARE
    record RECORD;
    notification JSON;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        record = OLD;
    ELSE
        record = NEW;
    END IF;
    notification = json_build_object(
                      'table',TG_TABLE_NAME,
                      'action', TG_OP,
                      'id', record.job_id);
    PERFORM pg_notify('events', notification::text);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_event
AFTER INSERT OR UPDATE OR DELETE ON jobs
    FOR EACH ROW EXECUTE PROCEDURE jobs_notify_event();

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS job_statuses;
DROP TABLE IF EXISTS agents;
DROP TABLE IF EXISTS agent_statuses;
