-- name: CreateTodo :one
INSERT INTO todos (
    category_id,
    user_id,
    title,
    content
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListTodoByUser :many
SELECT
    t.category_id, t.user_id, t.title, t.content, t.created_at, t.updated_at,
    c.name as category_name, c.color as category_color
FROM todos t
LEFT JOIN categories c
    ON c.id = t.category_id
WHERE user_id = $1

ORDER BY created_at ASC
LIMIT $2
OFFSET $3;

-- name: UpdateTodoByUser :one
UPDATE todos
SET category_id = $2, user_id = $3, title = $4, content = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
