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
    t.id, t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at, t.date, t.color, t.is_priority, t.status,
    c.name as category_name
FROM todos t
INNER JOIN categories c
    ON c.id = t.category_id
WHERE t.user_email = $1 

ORDER BY created_at ASC
LIMIT $2
OFFSET $3;

-- name: ListTodayTodo :many
SELECT
    t.id, t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at, t.date, t.color, t.is_priority, t.status,
    c.name as category_name
FROM todos t
INNER JOIN categories c
    ON c.id = t.category_id
WHERE t.user_email = $1 
    AND date <= now() 
    AND status = FALSE 
ORDER BY is_priority DESC
LIMIT $2
OFFSET $3;

-- name: ListUpcomingTodo :many
SELECT
    t.id, t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at, t.date, t.color, t.is_priority, t.status,
    c.name as category_name
FROM todos t
INNER JOIN categories c
    ON c.id = t.category_id
WHERE t.user_email = $1 
    AND date > now() 
    AND status = FALSE 
ORDER BY is_priority DESC, date ASC
LIMIT $2
OFFSET $3;

-- name: ListDoneTodo :many
SELECT
    t.id, t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at, t.date, t.color, t.is_priority, t.status,
    c.name as category_name
FROM todos t
INNER JOIN categories c
    ON c.id = t.category_id
WHERE t.user_email = $1 
    AND status = TRUE 
ORDER BY date DESC
LIMIT $2
OFFSET $3;


-- name: GetTodo :one
SELECT
    t.id, t.category_id, t.user_email, t.title, t.content, t.created_at, t.updated_at, t.date, t.color, t.is_priority,
    c.name as category_name
FROM todos t
INNER JOIN categories c
    ON c.id = t.category_id
WHERE t.id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateTodoByUser :one
UPDATE todos
SET category_id = $2, title = $3, content = $4, updated_at = now(), date = $5, color = $6, is_priority = $7
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
