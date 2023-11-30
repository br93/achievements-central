-- name: CreateUser :one
INSERT INTO tb_users (email, password)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM tb_users 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM tb_users 
WHERE email = $1 AND deleted_at IS NULL;