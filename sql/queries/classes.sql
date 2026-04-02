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
INSERT INTO users_classes (id, joined_at, updated_at, user_id, class_id)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
)
RETURNING *;

-- name: GetClassesForUserID :many
SELECT * FROM users_classes
WHERE user_id = $1;
