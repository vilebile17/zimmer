-- name: CreateAssignment :one
INSERT INTO assignments (id, class_id, created_at, updated_at, due_at, title, instructions, allow_late)
VALUES (
        gen_random_uuid(),
        $1,
        NOW(),
        NOW(),
        $2,
        $3,
        $4,
        $5
)
RETURNING *;
