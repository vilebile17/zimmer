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

-- name: GetNumAssignmentsToDo :one
SELECT COUNT(assignments.title) FROM students_classes
INNER JOIN classes ON classes.id = students_classes.class_id
INNER JOIN assignments ON assignments.class_id = classes.id
WHERE students_classes.student_id = $1
AND NOW() < assignments.due_at;

-- name: DeleteAssignment :one
DELETE FROM assignments
WHERE id = $1
RETURNING *;

-- name: UpdateAssignment :one
UPDATE assignments
SET
        updated_at = NOW(),
        title = $2,
        instructions = $3,
        due_at = $4
WHERE id = $1
RETURNING *;
