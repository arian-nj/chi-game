-- +goose Up
-- +goose StatementBegin
CREATE TABLE room_players (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL REFERENCES rooms(id),
    person_id BIGINT NOT NULL REFERENCES persons(id),
    UNIQUE(room_id, person_id) -- a player can't be added twice
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
