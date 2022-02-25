-- name: CreateUser :one
INSERT INTO users (
    name,
    address,
    pic,
    hashed_password,
    email
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, address = $3, pic = $4, email = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
