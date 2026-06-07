-- name: InsertRoom :one
INSERT INTO rooms (code, host_person_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: InsertRoomPlayer :one
INSERT INTO room_players (room_id, person_id)
VALUES ($1, $2)
ON CONFLICT (room_id, person_id) DO NOTHING
RETURNING *;

-- name: DeleteRoomPlayer :exec
DELETE FROM room_players WHERE room_id = $1 AND person_id = $2;