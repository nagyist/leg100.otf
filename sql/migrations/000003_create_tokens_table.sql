-- +goose Up
CREATE TABLE IF NOT EXISTS tokens (
    token_id text,
    created_at timestamptz,
    updated_at timestamptz,
    description text,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
