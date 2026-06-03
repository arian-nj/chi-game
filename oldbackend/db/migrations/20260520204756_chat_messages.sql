-- +goose Up
CREATE TABLE chat_messages (
    message_id BIGSERIAL PRIMARY KEY,
    chat_room_id_ref BIGINT NOT NULL REFERENCES chat_rooms(chat_room_id),
    sender_person_id BIGINT NOT NULL REFERENCES persons(id), -- 0 is system messages (not from user)
    content TEXT NOT NULL,
    reply_to_message_id BIGINT NULL REFERENCES chat_messages(message_id) ON DELETE SET NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,                -- soft delete (message removed for everyone)
    metadata JSONB NULL 
    
);

CREATE TABLE message_attachments (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES chat_messages(message_id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_type VARCHAR(50),                      -- e.g., 'image', 'document'
    file_name VARCHAR(255),
    file_size BIGINT
);

CREATE INDEX idx_message_attachments_message_id ON message_attachments(message_id);

-- +goose Down
SELECT 'down SQL query';
