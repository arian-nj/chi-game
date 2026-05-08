-- name: GetPersonsFriends :many
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
ON CONFLICT (sender_id, receiver_id) DO NOTHING;

-- name: AcceptFriendRequest :exec
UPDATE friend_requests
SET status = 'accepted', updated_at = NOW()
WHERE id = $1 AND status = 'pending';

-- name: RejectFriendRequest :exec
UPDATE friend_requests
SET status = 'rejected', updated_at = NOW()
WHERE id = $1 AND status = 'pending';

