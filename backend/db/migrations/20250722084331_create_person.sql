-- +goose Up
-- +goose StatementBegin
CREATE TABLE persons(
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(32) UNIQUE,
  is_guest BOOLEAN DEFAULT FALSE,
  updated_at timestamp NOT NULL DEFAULT NOW(),
  created_at timestamp NOT NULL DEFAULT NOW()
  );

CREATE TABLE guest_user(

);

-- CREATE TABLE user_platforms (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
--     platform VARCHAR(20) NOT NULL, -- 'telegram', 'bale', 'web'
--     platform_user_id VARCHAR(100) NOT NULL, -- platform's user_id
--     guest_identifier VARCHAR(255), -- for web guests (stored in localStorage)
--     created_at TIMESTAMP DEFAULT NOW(),
--     UNIQUE(platform, platform_user_id)
-- );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

-- +goose StatementEnd
