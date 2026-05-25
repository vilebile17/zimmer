-- name: CreateClass :one
INSERT INTO classes (id, created_at, updated_at, name, teacher_id, allow_joining)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2,
        true
)
RETURNING *;

-- name: JoinClass :one
INSERT INTO students_classes (id, joined_at, updated_at, student_id, class_id)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
)
RETURNING *;

-- name: GetClass :one
SELECT classes.id, classes.created_at, classes.updated_at,
        classes.name, classes.teacher_id, classes.allow_joining, users.name as teacher_name
FROM classes
INNER JOIN users
        ON classes.teacher_id = users.id
WHERE classes.id = $1;

-- name: GetClassesAsStudent :many
SELECT * FROM classes
WHERE id IN (
        SELECT class_id
        FROM students_classes
        WHERE student_id = $1
);

-- name: GetClassesAsTeacher :many
SELECT * FROM classes
WHERE teacher_id = $1;

-- name: GetClassFromClassID :one
SELECT * FROM classes
WHERE id = $1;

-- name: UpdateClass :one
UPDATE classes
SET
        name = $2,
        teacher_id = $3,
        allow_joining = $4,
        updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: RemoveUserFromClass :exec
DELETE FROM students_classes
WHERE student_id = $1 AND class_id = $2;
