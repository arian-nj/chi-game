-- +goose Up
CREATE TABLE room_participants (
  room_id BIGINT NOT NULL REFERENCES rooms(id),
  user_id BIGINT NOT NULL REFERENCES rooms(id),
  last_read_message_id BIGINT NOT NULL REFERENCES room_messages(id),
  joined_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (room_id, user_id)
);

CREATE INDEX idx_room_participants_user_id ON room_participants(user_id);
CREATE INDEX idx_room_participants_room_id ON room_participants(room_id);
-- +goose Down
SELECT 'down SQL query';
