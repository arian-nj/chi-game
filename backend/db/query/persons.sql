-- name: GetPersonByAuthMethod :one
SELECT p.* FROM persons p
JOIN person_auth_methods pam ON p.id = pam.user_id
WHERE pam.auth_type = $1 AND pam.auth_value = $2
LIMIT 1;

-- name: InsertPerson :one
INSERT INTO persons (username,display_name, is_guest)
VALUES ($1,$2, $3) RETURNING *;

-- name: InsertAuthMethod :one
INSERT INTO person_auth_methods (user_id,auth_type,auth_value)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetPersonByID :one
SELECT * FROM persons
WHERE id = $1;

-- name: SearchPersonByUsername :many
SELECT id,username,display_name FROM persons
WHERE username LIKE $1 || '%' ORDER BY username LIMIT 5;  

-- name: GetUserMethod :one
SELECT user_id FROM person_auth_methods 
WHERE auth_type = $1 AND auth_value = $2;

-- -- name: GetTgUserByTgID :one
-- SELECT * FROM persons
-- WHERE tg_id = $1;


-- name: GetAllTgUsers :many
SELECT * FROM persons;

-- name: CountAllTgUsers :one
SELECT COUNT(*) FROM persons;

-- -- name: CreateTgUser :one
-- INSERT INTO persons (tg_id,name)
-- VALUES ($1,$2)
-- ON CONFLICT (tg_id) DO UPDATE
-- SET updated_at = NOW(),
--     is_active = TRUE
-- RETURNING *;



-- name: CountActiveTgUsers :one
SELECT COUNT(*) FROM persons
WHERE is_active = TRUE;

-- name: CountUsersTgCreatedBetween :one
SELECT COUNT(*) FROM persons
WHERE created_at >= $1 AND created_at <= $2;

--
-- -- name: ActiveTgUser :exec
-- UPDATE persons SET is_active = TRUE
-- WHERE (
--   tg_id = $1
-- );
--
-- -- name: DiactiveTgUser :exec
-- UPDATE persons SET is_active = FALSE
-- WHERE (
--   tg_id = $1
-- );
--

-- name: UpdateMixedTgUserStatuses :exec
UPDATE persons
SET is_active = CASE
  WHEN id = ANY ($1::bigint[]) THEN TRUE
  ELSE FALSE  
END;


-- name: GetTgUsersStatic :one
SELECT
  COUNT(*) FILTER (WHERE is_active = TRUE) AS active_users,
  COUNT(*) AS total_users,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '1 day') AS users_created_last_24_hours,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '7 days') AS users_created_last_week,
  COUNT(*) FILTER ( WHERE created_at >= NOW() - INTERVAL '1 month') AS users_created_last_month
FROM persons;
