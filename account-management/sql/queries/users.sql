-- name: GetUserById :one
SELECT * FROM tb_users 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM tb_users 
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetAllUsers :many
SELECT * FROM tb_users WHERE deleted_at IS NULL;

-- name: GetAllActiveUsers :many
SELECT * FROM tb_users WHERE is_active = true AND deleted_at IS NULL;

-- name: UpdateEmail :exec
UPDATE tb_users SET email = $1
WHERE id = $2;

-- name: UpdatePassword :exec
UPDATE tb_users SET password = $1
WHERE id = $2;

-- name: UpdateSuperUser :exec
UPDATE tb_users SET is_superuser = $1
WHERE id = $2;

-- name: UpdateActive :exec
UPDATE tb_users SET is_active = $1
WHERE id = $2;

-- name: DeleteUser :exec
DELETE from tb_users 
WHERE id = $1;

