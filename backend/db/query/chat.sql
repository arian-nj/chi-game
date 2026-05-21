-- name: DoesChatExist :one
SELECT c.id
FROM chats c
JOIN chat_participants p1 ON p1.chat_id = c.id
JOIN chat_participants p2 ON p2.chat_id = c.id
WHERE c.type = 'direct'
  AND p1.user_id = $1
  AND p2.user_id = $2;

-- name: GetLastMessagesOfChat :many
SELECT * FROM chat_messages
WHERE chat_id = $1 AND deleted_at IS NULL
ORDER BY sent_at DESC
LIMIT 20;

-- name: GetAllChatsOfUser :many
SELECT c.*
FROM chats c
JOIN chat_participants p ON p.chat_id = c.id
WHERE p.user_id = 123 AND p.left_at IS NULL;

-- name: SetLastReadId :exec
UPDATE chat_participants
SET last_read_message_id = $1
WHERE chat_id = $2 AND user_id = $3;
