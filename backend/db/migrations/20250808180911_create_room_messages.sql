-- +goose Up
-- +goose StatementBegin
CREATE TABLE room_messages (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL REFERENCES rooms(id),
    person_id BIGINT NOT NULL REFERENCES persons(id),
    reply_to BIGINT REFERENCES room_messages(id),
    content TEXT NOT NULL,
    is_edited BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_messages_room_id_created_at ON room_messages(room_id, created_at DESC);
CREATE INDEX idx_messages_user_id ON room_messages(person_id);
CREATE INDEX idx_messages_parent_id ON room_messages(reply_to);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
