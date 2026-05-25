-- +goose Up
CREATE TABLE chat_messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_user_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    reply_to_message_id BIGINT NULL REFERENCES chat_messages(id) ON DELETE SET NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL                -- soft delete (message removed for everyone)
    
);

CREATE INDEX idx_chat_messages_chat_id_sent_at ON chat_messages(chat_id, sent_at);
CREATE INDEX idx_chat_messages_sender_user_id ON chat_messages(sender_user_id);
CREATE INDEX idx_chat_messages_reply_to ON chat_messages(reply_to_message_id);

CREATE TABLE message_attachments (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_type VARCHAR(50),                      -- e.g., 'image', 'document'
    file_name VARCHAR(255),
    file_size BIGINT
);

CREATE INDEX idx_message_attachments_message_id ON message_attachments(message_id);

-- +goose Down
SELECT 'down SQL query';
