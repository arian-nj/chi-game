-- +goose Up
-- +goose StatementBegin
CREATE TABLE invite_rooms (
    id BIGSERIAL PRIMARY KEY,
    invite_code VARCHAR(8) NOT NULL UNIQUE,
    game_key VARCHAR(32) NOT NULL,
    host_person_id BIGINT NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_invite_rooms_code ON invite_rooms(invite_code);
CREATE INDEX idx_invite_rooms_expires ON invite_rooms(expires_at);

CREATE TABLE invite_room_players (
    room_id BIGINT NOT NULL REFERENCES invite_rooms(id) ON DELETE CASCADE,
    person_id BIGINT NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (room_id, person_id)
);

CREATE INDEX idx_invite_room_players_person ON invite_room_players(person_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invite_room_players;
DROP TABLE IF EXISTS invite_rooms;
-- +goose StatementEnd
