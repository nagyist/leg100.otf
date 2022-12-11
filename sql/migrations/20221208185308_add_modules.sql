-- +goose Up
CREATE TABLE IF NOT EXISTS modules (
    module_id       TEXT,
    created_at      TIMESTAMPTZ NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL,
    name            TEXT        NOT NULL,
    provider        TEXT        NOT NULL,
    organization_id TEXT REFERENCES organizations ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
                    PRIMARY KEY (module_id)
);

CREATE TABLE IF NOT EXISTS module_versions (
    module_version_id TEXT,
    version           TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL,
    module_id         TEXT REFERENCES modules ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
                      PRIMARY KEY (module_version_id)
);

CREATE TABLE IF NOT EXISTS module_tarballs (
    tarball           BYTEA NOT NULL,
    module_version_id TEXT REFERENCES module_versions ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
                    UNIQUE (module_version_id)
);

CREATE TABLE IF NOT EXISTS module_repos (
    -- do not cascade deletes because the otfd code relies on getting an error
    -- when attempting to delete a webhook, to determine whether there are any
    -- module repos referencing it; only when no more module repos are referencing
    -- a webhook do we delete it.
    webhook_id        UUID REFERENCES webhooks ON UPDATE CASCADE NOT NULL,
    vcs_provider_id   TEXT REFERENCES vcs_providers ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    module_id         TEXT REFERENCES modules ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
                      UNIQUE (module_id)
);

-- +goose Down
DROP TABLE IF EXISTS module_repos;
DROP TABLE IF EXISTS module_tarballs;
DROP TABLE IF EXISTS module_versions;
DROP TABLE IF EXISTS modules;
