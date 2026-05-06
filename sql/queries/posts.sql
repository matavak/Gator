-- name: CreatePost :one
INSERT INTO posts (id,created_at,updated_at,feed_id,title,url,description,published_at)
Values(gen_random_uuid(),NOW(),NOW(),$1,$2,$3,$4,$5)
	RETURNING *;
-- name: GetPostsForUser :many
select posts.* from posts inner join feeds ON feeds.id=posts.feed_id inner join feed_follows ON feeds.id=feed_follows.feed_id WHERE feed_follows.user_id=$1 ORDER BY published_at DESC LIMIT $2;
