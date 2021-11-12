-- +goose Up
CREATE TABLE IF NOT EXISTS tokens (
    created_at timestamptz,
    updated_at timestamptz,
    id text,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
