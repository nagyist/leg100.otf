-- +goose Up
CREATE TABLE IF NOT EXISTS tokens (
    token_id text,
    created_at timestamptz,
    updated_at timestamptz,
    description text,
    hash bytea,
    PRIMARY KEY (token_id)
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
