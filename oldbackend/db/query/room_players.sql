-- name: CreateRoomPlayer :one
INSERT INTO room_players (room_id,person_id) VALUES ($1,$2) RETURNING *;
