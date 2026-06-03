-- +goose Up
CREATE TABLE friend_requests (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES persons(id),
    receiver_id BIGINT NOT NULL REFERENCES persons(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, accepted, rejected, canceled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_pending_request
ON friend_requests (sender_id, receiver_id)
WHERE status = 'pending'; 

CREATE INDEX idx_friend_requests_sender ON friend_requests(sender_id, status);
CREATE INDEX idx_friend_requests_receiver ON friend_requests(receiver_id, status);

CREATE TABLE friends (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES persons(id),
    friend_id BIGINT NOT NULL REFERENCES persons(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, friend_id),
    CHECK (user_id != friend_id)
);

-- +goose Down
SELECT 'down SQL query';
