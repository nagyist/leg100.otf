-- +goose Up
INSERT INTO tokens (
    token_id,
    created_at,
    updated_at,
    description
) VALUES (
    'at-default-token',
    now(),
    now(),
    'my token'
), (
    'at-another-token',
    now(),
    now(),
    'another token'
);

-- +goose Down
DELETE
FROM tokens
WHERE token_id = 'at-default-token';
