-- name: CreateRoomGame :one
INSERT INTO room_games (room_id,game_type) VALUES ($1,$2) RETURNING *;
