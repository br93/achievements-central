-- name: CreateUser :one
INSERT INTO tb_users (email, password)
VALUES ($1, $2)
RETURNING *;