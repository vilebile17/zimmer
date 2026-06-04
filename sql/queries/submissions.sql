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

-- name: GetSubmissionsForAssignment :many
SELECT submissions.id, submissions.created_at, submissions.updated_at,
        submissions.user_id, submissions.answers, users.name
FROM submissions
INNER JOIN users
        ON submissions.user_id = users.id
WHERE assignment_id = $1;

-- name: GetSubmission :one
SELECT * FROM submissions
WHERE id = $1;

-- name: GradeSubmission :one
UPDATE submissions
SET
        score = $2,
        updated_at = NOW()
WHERE id = $1
RETURNING *;
