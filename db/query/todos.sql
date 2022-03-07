-- name: CreateTodo :one
INSERT INTO todos (
    category_id,
    user_email,
    title,
    content
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListTodoByUser :many
SELECT
    t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at,
    c.name as category_name, c.color as category_color
FROM todos t
LEFT JOIN categories c
    ON c.id = t.category_id
WHERE t.user_email = $1

ORDER BY created_at ASC
LIMIT $2
OFFSET $3;

-- name: GetTodo :one
SELECT
    t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at,
    c.name as category_name, c.color as category_color
FROM todos t
LEFT JOIN categories c
    ON c.id = t.category_id
WHERE t.id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateTodoByUser :one
UPDATE todos
SET category_id = $2, user_email = $3, title = $4, content = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
