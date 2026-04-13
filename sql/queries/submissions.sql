-- name: CreateSubmission :one
INSERT INTO submissions (id, created_at, updated_at, assignment_id, user_id, answers)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2,
        $3
)
RETURNING *;
