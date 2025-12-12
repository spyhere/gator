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
-- WITH feed_item AS (
--   INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
--   VALUES ($1, $2, $3, $4, $5, $6)
--     ON CONFLICT (url) DO NOTHING
--   RETURNING id
-- )
-- INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
-- SELECT $1, $2, $3, $6, id FROM feed_item
-- RETURNING *;

