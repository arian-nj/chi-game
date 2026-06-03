-- +goose Up
CREATE TABLE chat_rooms (
    chat_room_id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL                -- soft delete for whole chat
);

-- +goose Down
SELECT 'down SQL query';
