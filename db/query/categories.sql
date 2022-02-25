-- name: CreateCategory :one
INSERT INTO categories (
    name,
    color
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, color = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;
