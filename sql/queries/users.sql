-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, hashed_password)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2,
        $3
)
RETURNING *;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetStudentsForClass :many
SELECT id,name FROM users
WHERE id IN (
        SELECT student_id
        FROM students_classes
        WHERE class_id = $1
);

-- name: UpdateUser :one
UPDATE users
SET
        name = $2,
        email = $3,
        hashed_password = $4,
        updated_at = NOW()
WHERE id = $1
RETURNING *;
