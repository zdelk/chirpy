-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
