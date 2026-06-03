-- name: GetPersonByAuthMethod :one
SELECT p.id, p.display_name, p.username, p.coins, p.is_guest, p.merged_at, p.updated_at, p.created_at
FROM persons p
JOIN person_auth_methods pam ON p.id = pam.user_id
WHERE pam.auth_type = $1 AND pam.auth_value = $2
LIMIT 1;

-- name: InsertPerson :one
INSERT INTO persons (username, display_name, is_guest)
VALUES ($1, $2, $3)
RETURNING id, display_name, username, coins, is_guest, merged_at, updated_at, created_at;

-- name: InsertAuthMethod :one
INSERT INTO person_auth_methods (user_id, auth_type, auth_value)
VALUES ($1, $2, $3)
RETURNING id, user_id, auth_type, auth_value, created_at;

-- name: GetPersonByID :one
SELECT id, display_name, username, coins, is_guest, merged_at, updated_at, created_at
FROM persons
WHERE id = $1;

-- name: CountPersons :one
SELECT COUNT(*)::bigint FROM persons;

-- name: CountGuestPersons :one
SELECT COUNT(*)::bigint FROM persons WHERE is_guest = true;

-- name: ListRecentPersons :many
SELECT id, username, is_guest, created_at
FROM persons
ORDER BY created_at DESC
LIMIT 20;
