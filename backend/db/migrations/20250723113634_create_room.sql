-- +goose Up
-- +goose StatementBegin

create table rooms (
  id bigserial primary key,
  name  varchar(100) not null, -- null for direct messages
  rtype varchar(20) not null, -- direct or group 
  created_by bigint not null references persons(id),
  createdmode text not null, -- private, random
  created_at timestamp not null default now()
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
