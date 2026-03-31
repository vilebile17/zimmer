-- name: CreateClass :one
INSERT INTO classes (id, created_at, updated_at, name, teacher_id)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
)
RETURNING *;
