-- name: InsertRoomMessage :one
INSERT INTO room_messages (room_id, person_id, message)
VALUES ($1, $2, $3)
RETURNING id, room_id, person_id, message, created_at;

-- name: GetRoomMessages :many
SELECT id, room_id, person_id, message, created_at
FROM room_messages
WHERE room_id = $1
ORDER BY created_at ASC LIMIT $2 OFFSET $3;