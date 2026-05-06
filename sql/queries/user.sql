-- name: CreateUser :one
INSERT INTO users (id,created_at,updated_at,name)
Values(gen_random_uuid(),NOW(),NOW(),$1)
	RETURNING *;
-- name: GetUser :one
SELECT * FROM users WHERE name=$1;
-- name: ResetUsers :exec
DELETE FROM users;
-- name: GetAllUsers :many
SELECT * FROM users ORDER BY created_at DESC;
-- name: GetUserByID :one
SELECT * FROM users WHERE id=$1;
