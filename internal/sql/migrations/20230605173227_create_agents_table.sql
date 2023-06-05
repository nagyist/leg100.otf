-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS agents;
DROP TABLE IF EXISTS agent_statuses;
