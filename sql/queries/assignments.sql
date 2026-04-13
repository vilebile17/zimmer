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

-- name: GetAssignmentsForAClass :many
SELECT * FROM assignments
WHERE class_id = $1
ORDER BY created_at DESC;

-- name: GetAssignmentFromID :one
SELECT * FROM assignments
WHERE id = $1;
