-- name: CreateFeed :one
INSERT INTO feeds (id,created_at,updated_at,name,url,user_id)
Values(gen_random_uuid(),NOW(),NOW(),$1,$2,$3)
	RETURNING *;
-- name: GetAllFeeds :many
SELECT * FROM feeds ORDER BY created_at desc;
-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url=$1 ORDER BY created_at desc;
-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at=NOW() , updated_at=NOW() WHERE id=$1;
-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS first LIMIT 1;
