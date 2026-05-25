
-- name: CreateOrGetChatRoom :one
INSERT INTO chat_rooms DEFAULT VALUES RETURNING *;

-- name: DoesChatExist :one
SELECT c.chat_room_id
FROM chat_rooms c
JOIN chat_participants p1 ON p1.chat_room_id_ref = c.chat_room_id
JOIN chat_participants p2 ON p2.chat_room_id_ref = c.chat_room_id
  AND p1.person_id = $1
  AND p2.person_id = $2;

-- name: GetLastMessagesOfChat :many
SELECT * FROM chat_messages
WHERE chat_room_id_ref = $1 AND deleted_at IS NULL
ORDER BY sent_at DESC
LIMIT 20;

-- name: GetMessagesOfChat :many
SELECT * FROM chat_messages
WHERE chat_room_id_ref = $1 AND deleted_at IS NULL
ORDER BY sent_at DESC
LIMIT 20;

-- name: GetAllChatsOfUser :many
SELECT 
    c.chat_room_id,
    c.created_at,
    c.updated_at,
    c.deleted_at,
    other_p.person_id AS other_person_id
FROM chat_rooms c
JOIN chat_participants p ON p.chat_room_id_ref = c.chat_room_id
JOIN chat_participants other_p ON other_p.chat_room_id_ref = c.chat_room_id
    AND other_p.person_id != $1
    AND other_p.left_at IS NULL
WHERE p.person_id = $1 AND p.left_at IS NULL
GROUP BY c.chat_room_id, c.created_at, c.updated_at, c.deleted_at, other_p.person_id;

-- name: InsertChatParticipants :exec
INSERT INTO chat_participants (chat_room_id_ref, person_id)
VALUES ($1, $2)
ON CONFLICT (chat_room_id_ref, person_id) DO NOTHING;

-- name: IsChatParticipant :one
SELECT * FROM chat_participants
WHERE chat_room_id_ref = $1 AND person_id = $2;


-- name: SetLastReadId :exec
UPDATE chat_participants
SET last_read_message_id = $1
WHERE chat_room_id_ref = $2 AND person_id = $3;

-- name: InsertChatMessage :one
INSERT INTO chat_messages (chat_room_id_ref, sender_person_id, content, metadata)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: MarkMessageAsRead :exec
UPDATE chat_messages
SET read_at = NOW()
WHERE message_id = $1;