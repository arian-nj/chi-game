-- +goose Up
CREATE TABLE games (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    min_players INT,
    max_players INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
SELECT 'down SQL query';
