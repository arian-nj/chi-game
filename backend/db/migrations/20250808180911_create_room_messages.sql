-- +goose Up
-- +goose StatementBegin
CREATE TABLE room_messages (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
