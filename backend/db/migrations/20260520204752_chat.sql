-- +goose Up
CREATE TYPE chat_type AS ENUM ('direct', 'group');

CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    type chat_type NOT NULL,
    name VARCHAR(255) NULL,                     -- only for group chats
    created_by BIGINT NOT NULL,                 -- user_id of creator
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL                 -- soft delete for whole chat
);

CREATE INDEX idx_chats_type ON chats(type);
CREATE INDEX idx_chats_created_by ON chats(created_by);

-- +goose Down
SELECT 'down SQL query';
