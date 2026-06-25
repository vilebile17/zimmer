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
WHERE assignment_id = $1
ORDER BY submissions.updated_at ASC;

-- name: GetSubmission :one
SELECT submissions.id, submissions.created_at, submissions.updated_at,
        submissions.user_id as user_id, submissions.answers, users.name as user_name,
        assignments.title as assignment_title, classes.id as class_id,
        assignments.id as assignment_id, submissions.score as grade FROM submissions
INNER JOIN assignments ON submissions.assignment_id = assignments.id
INNER JOIN users ON submissions.user_id = users.id
INNER JOIN classes ON assignments.class_id = classes.id
WHERE submissions.id = $1;

-- name: GradeSubmission :one
UPDATE submissions
SET
        score = $2,
        updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateSubmission :one
UPDATE submissions
SET
        answers = $2,
        updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetSubmissionForUser :one
SELECT * FROM submissions
WHERE assignment_id = $1 AND user_id = $2;
