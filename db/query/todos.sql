-- name: CreateTodo :one
INSERT INTO todos (
    category_id,
    user_email,
    title,
    content,
    date,
    color,
    is_priority
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListTodoByUser :many
SELECT
    t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at, t.date, t.color, t.is_priority,
    c.name as category_name
FROM todos t
LEFT JOIN categories c
    ON c.id = t.category_id
WHERE t.user_email = $1

ORDER BY created_at ASC
LIMIT $2
OFFSET $3;

-- name: GetTodo :one
SELECT
    t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at, t.date, t.color, t.is_priority,
    c.name as category_name
FROM todos t
LEFT JOIN categories c
    ON c.id = t.category_id
WHERE t.id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateTodoByUser :one
UPDATE todos
SET category_id = $2, user_email = $3, title = $4, content = $5, updated_at = now(), date = $6, color = $7, is_priority = $8
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
