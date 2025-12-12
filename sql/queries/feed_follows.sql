-- name: CreateFeedFollow :one
WITH feed_follows_ins AS (
  INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT feed_follows_ins.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM feed_follows_ins
INNER JOIN feeds
ON feeds.id = feed_follows_ins.feed_id
INNER JOIN users
ON users.id = feed_follows_ins.user_id;

-- name: GetUserFeeds :many
SELECT
  feeds.name AS feed_name,
  feeds.url,
  feed_follows.updated_at
FROM users
INNER JOIN feed_follows
ON feed_follows.user_id = users.id
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE users.id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE (feed_follows.user_id = $1 AND feed_follows.feed_id = $2);

