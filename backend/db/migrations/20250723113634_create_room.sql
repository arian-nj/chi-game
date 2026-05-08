-- +goose Up
-- +goose StatementBegin

CREATE TABLE rooms (
  id BIGSERIAL PRIMARY KEY,
  name  varchar(100) NOT NULL, -- null for direct messages
  rtype varchar(20) NOT NULL, -- direct or group 
  created_by BIGINT NOT NULL REFERENCES person(id),
  createdMode TEXT NOT NULL, -- private, random
  created_at timestamp NOT NULL DEFAULT NOW()
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
