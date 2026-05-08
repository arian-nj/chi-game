-- +goose Up
CREATE TABLE user_game_stats (
    user_id BIGINT REFERENCES persons(id) ON DELETE CASCADE,
    game_id BIGINT REFERENCES games(id) ON DELETE CASCADE,
    wins INT DEFAULT 0,
    losses INT DEFAULT 0,
    draws INT DEFAULT 0,
    elo_rating INT DEFAULT 1200, -- if competitive
    xp INT DEFAULT 0,
    level INT DEFAULT 1,
    last_played_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, game_id)
);
-- +goose Down
SELECT 'down SQL query';
