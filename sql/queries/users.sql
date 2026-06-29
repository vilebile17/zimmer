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
SELECT id,name,created_at FROM users
WHERE id IN (
        SELECT student_id
        FROM students_classes
        WHERE class_id = $1
)
ORDER BY name;

-- name: UpdateUserImportant :one
UPDATE users
SET
        email = $2,
        hashed_password = $3,
        updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetTotalUserCount :one
SELECT COUNT(*) FROM users;

-- name: UpdateUserLessImportant :one
UPDATE users
SET
        name = $2,
        bio = $3,
        updated_at = NOW()
WHERE id = $1
RETURNING *;
