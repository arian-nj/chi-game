-- name: CreateRoom :one
INSERT INTO rooms (createdMode) VALUES ($1) RETURNING *;

-- name: CountRooms :one
Select Count(*) from rooms;

-- name: CountLastHourHub :one
SELECT COUNT(*) AS hubs_last_hour
FROM rooms
WHERE created_at >= NOW() - INTERVAL '1 hour';

-- name: CountLastDayRoomss :one
SELECT COUNT(*) AS games_today
FROM rooms
WHERE created_at >= date_trunc('day', NOW());

--
-- --name: CountLast7DaysSesion :many
-- WITH days AS (
--   SELECT generate_series(
--     date_trunc('day', NOW()) - interval '6 days',
--     date_trunc('day', NOW()),
--     interval '1 day'
--   ) AS day
-- )
-- SELECT
--   days.day,
--   COUNT(h.id) AS hubs_count
-- FROM days
-- LEFT JOIN hubs h
--   ON h.created_at >= days.day
--  AND h.created_at < days.day + interval '1 day'
-- GROUP BY days.day
-- ORDER BY days.day DESC;
