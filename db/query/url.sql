-- name: CreateURL :one
INSERT INTO urls (
    original_url,
    short_code,
    expires_at,
    is_custom
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetURLByShortCode :one
SELECT * FROM urls
WHERE short_code = $1 AND  expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: DeleteExpiredURLs :exec
DELETE FROM urls
WHERE expires_at <= CURRENT_TIMESTAMP;

-- name: IsShortCodeAvailable :one
SELECT NOT EXISTS (
    SELECT 1 FROM urls
    WHERE short_code = $1
) AS is_available;


