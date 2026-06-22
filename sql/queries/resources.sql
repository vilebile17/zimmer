-- name: CreateResource :one
INSERT INTO resources (id, created_at, updated_at, class_id, title, content)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2,
        $3
)
RETURNING *;

-- name: GetResourcesForClass :many
SELECT id, title, created_at FROM resources
WHERE class_id = $1
ORDER BY created_at DESC;

-- name: GetResource :one
SELECT * FROM resources
WHERE id = $1;

-- name: UpdateResource :one
UPDATE resources
SET
        title = $2,
        content = $3,
        updated_at = NOW()
WHERE id = $1
RETURNING *;
