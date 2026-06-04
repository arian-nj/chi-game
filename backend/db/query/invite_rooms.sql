-- name: InsertInviteRoom :one
INSERT INTO invite_rooms (invite_code, game_key, host_person_id, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, invite_code, game_key, host_person_id, created_at, expires_at;

-- name: GetInviteRoomByCode :one
SELECT id, invite_code, game_key, host_person_id, created_at, expires_at
FROM invite_rooms
WHERE invite_code = $1 AND expires_at > NOW();

-- name: InsertInviteRoomPlayer :exec
INSERT INTO invite_room_players (room_id, person_id)
VALUES ($1, $2)
ON CONFLICT (room_id, person_id) DO NOTHING;

-- name: CountInviteRoomPlayers :one
SELECT COUNT(*)::bigint
FROM invite_room_players
WHERE room_id = $1;

-- name: ListInviteRoomPlayers :many
SELECT p.id, p.display_name, p.username, p.coins, p.is_guest, p.merged_at, p.updated_at, p.created_at
FROM invite_room_players irp
JOIN persons p ON p.id = irp.person_id
WHERE irp.room_id = $1
ORDER BY irp.joined_at ASC;

-- name: GetInviteRoomForPerson :one
SELECT r.id, r.invite_code, r.game_key, r.host_person_id, r.created_at, r.expires_at
FROM invite_room_players irp
JOIN invite_rooms r ON r.id = irp.room_id
WHERE irp.person_id = $1 AND r.expires_at > NOW()
LIMIT 1;

-- name: DeleteInviteRoomPlayer :exec
DELETE FROM invite_room_players
WHERE room_id = $1 AND person_id = $2;

-- name: DeleteInviteRoom :exec
DELETE FROM invite_rooms
WHERE id = $1;
