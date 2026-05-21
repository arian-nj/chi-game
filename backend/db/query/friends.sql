-- name: GetPersonsFriendsRequest :many
SELECT 
    CASE 
        WHEN fr.sender_id = $1 THEN fr.receiver_id
        ELSE fr.sender_id
    END AS friend_id
FROM friend_requests fr
WHERE (fr.sender_id = $1 OR fr.receiver_id = $1)
  AND fr.status = 'accepted';

-- name: InsertFriendRequest :exec
INSERT INTO friend_requests (sender_id, receiver_id)
VALUES ($1, $2)
ON CONFLICT (sender_id, receiver_id) 
WHERE status = 'pending' 
DO NOTHING;

-- name: AcceptFriendRequest :exec
UPDATE friend_requests
SET status = 'accepted', updated_at = NOW()
WHERE sender_id = $1 AND receiver_id=$2 AND status = 'pending';

-- name: RejectFriendRequest :exec
UPDATE friend_requests
SET status = 'rejected', updated_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: GetPersonFriends :many
SELECT f1.friend_id FROM friends f1 WHERE f1.user_id=$1
UNION 
SELECT f2.user_id FROM friends f2 WHERE f2.friend_id=$1;

-- name: CancelFriendRequest :one
UPDATE friend_requests
SET status = 'canceled' WHERE sender_id=$1 AND receiver_id=$2 RETURNING *;


-- name: CheckFriendship :one
SELECT EXISTS (
    SELECT 1 FROM friends 
    WHERE (user_id = $1 AND friend_id = $2)
       OR (user_id = $2 AND friend_id = $1)
) AS is_friend;


-- name: GetFriendshipStatus :one
SELECT 
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM friends 
            WHERE (user_id = $1 AND friend_id = $2)
               OR (user_id = $2 AND friend_id = $1)
        ) THEN 'friends'
        WHEN EXISTS (
            SELECT 1 FROM friend_requests 
            WHERE sender_id = $1 AND receiver_id = $2 AND status = 'pending'
        ) THEN 'request_sent'
        WHEN EXISTS (
            SELECT 1 FROM friend_requests 
            WHERE sender_id = $2 AND receiver_id = $1 AND status = 'pending'
        ) THEN 'request_received'
        ELSE 'not_connected'
    END AS status;



-- name: InsertFriend :exec
INSERT INTO friends (user_id, friend_id)
VALUES ($1, $2)
ON CONFLICT (user_id, friend_id) DO NOTHING;

