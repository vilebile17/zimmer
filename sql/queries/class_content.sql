-- name: CreateClassContent :one
INSERT INTO class_content (id, created_at, updated_at, content_type, class_id, title, content)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2,
        $3,
        $4
)
RETURNING *;

-- name: GetClassContentForClass :many
SELECT id, title, created_at FROM class_content
WHERE class_id = $2 AND content_type = $1
ORDER BY created_at DESC;

-- name: GetClassContent :one
SELECT * FROM class_content
WHERE id = $1;

-- name: UpdateClassContent :one
UPDATE class_content
SET
        title = $2,
        content = $3,
        updated_at = NOW()
WHERE id = $1
RETURNING *;
