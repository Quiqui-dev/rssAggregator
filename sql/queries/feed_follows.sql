-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, create_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetFeedFollows :many
SELECT * FROM feed_follows ff
WHERE ff.user_id = $1;


-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;


-- name: GetPopularFeeds :many
SELECT ff.feed_id, count(ff.feed_id) as num_follows, f.url, f.name
FROM feed_follows ff
JOIN feeds f ON ff.feed_id = f.id
GROUP BY ff.feed_id, f.url, f.name
ORDER BY num_follows DESC
OFFSET $1
LIMIT $2;