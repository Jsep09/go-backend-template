

-- name: GetExample :one
SELECT id, name, description, user_id, created_at, updated_at
FROM examples
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListExamples :many
SELECT id, name, description, user_id, created_at, updated_at
FROM examples
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateExample :one
INSERT INTO examples (name, description, user_id)
VALUES ($1, $2, $3)
RETURNING id, name, description, user_id, created_at, updated_at;

-- name: UpdateExample :one
UPDATE examples
SET name = $1,
    description = $2
WHERE id = $3 AND user_id = $4
RETURNING id, name, description, user_id, created_at, updated_at;

-- name: DeleteExample :exec
DELETE FROM examples
WHERE id = $1 AND user_id = $2;