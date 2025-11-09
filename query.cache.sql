-- name: GetStatsByKey :one
SELECT *
FROM cached
WHERE k = ?
AND datetime(created_at) >= datetime('now', '-1 hour');

-- name: SetStats :one
INSERT INTO cached (k, user, stars, repositories, commits, forks, top_repositories, topics, avatar_url)
VALUES (?,?,?,?,?,?,?,?,?)
RETURNING *;

-- name: CleanCache :exec
DELETE FROM cached
WHERE datetime(created_at) < datetime('now', '-1 hour');

-- name: CleanCacheTest :exec
DELETE FROM cached
WHERE datetime(created_at) < datetime('now', '-1 second');