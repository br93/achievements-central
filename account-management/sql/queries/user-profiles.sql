-- name: CreateUserProfile :one
INSERT INTO tb_user_profiles (user_id, steam_user, gog_user)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserProfileById :one
SELECT * FROM tb_user_profiles
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserProfileByUserId :one
SELECT * FROM tb_user_profiles
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: GetAllUserProfiles :many
SELECT * FROM tb_user_profiles
WHERE deleted_at IS NULL;

-- name: UpdateUserProfile :exec
UPDATE tb_user_profiles SET steam_user = $1, gog_user = $2
WHERE id = $3;

-- name: DeleteUserProfile :exec
DELETE FROM tb_user_profiles
WHERE id = $1;
