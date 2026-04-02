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
