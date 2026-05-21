-- +goose Up
-- +goose StatementBegin
CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    username VARCHAR(64) NOT NULL UNIQUE,
    coins INTEGER NOT NULL DEFAULT 100,
    is_guest BOOLEAN NOT NULL DEFAULT false,
    merged_at TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON persons(username);

CREATE TABLE person_auth_methods (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES persons(id),
    auth_type VARCHAR(50) NOT NULL, -- 'guest_device', 'telegram', 'phone', 'email', 'bale'
    auth_value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(auth_type, auth_value)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

-- +goose StatementEnd
