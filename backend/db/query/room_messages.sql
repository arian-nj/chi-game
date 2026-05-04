-- name: CreateRoomMessage :one
INSERT INTO room_messages (room_id,player_id,message) VALUES ($1,$2,$3) RETURNING *;

-- name: GetRoomMessages :many
SELECT *
FROM room_messages
WHERE room_id = $1
ORDER BY id DESC
LIMIT $2;
