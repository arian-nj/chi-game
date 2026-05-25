-- +goose Up
CREATE TABLE chat_participants (
    participant_id BIGSERIAL PRIMARY KEY,
    chat_room_id_ref BIGINT NOT NULL REFERENCES chat_rooms(chat_room_id) ON DELETE CASCADE,
    person_id BIGINT NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    left_at TIMESTAMPTZ NULL,                   -- soft leave/remove
    last_read_message_id BIGINT NULL,           -- last message read by this user in this chat
    UNIQUE(chat_room_id_ref, person_id)                    -- prevent duplicate membership
);

CREATE INDEX idx_chat_participants_person_id ON chat_participants(person_id);
CREATE INDEX idx_chat_participants_last_read ON chat_participants(chat_room_id_ref);


-- +goose Down
SELECT 'down SQL query';
