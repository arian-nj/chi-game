-- +goose Up
CREATE TABLE chat_participants (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- 'admin', 'member'
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    left_at TIMESTAMPTZ NULL,                   -- soft leave/remove
    last_read_message_id BIGINT NULL,           -- last message read by this user in this chat
    UNIQUE(chat_id, user_id)                    -- prevent duplicate membership
);

CREATE INDEX idx_chat_participants_chat_id ON chat_participants(chat_id);
CREATE INDEX idx_chat_participants_user_id ON chat_participants(user_id);
CREATE INDEX idx_chat_participants_last_read ON chat_participants(chat_id, last_read_message_id);
-- +goose Down
SELECT 'down SQL query';
