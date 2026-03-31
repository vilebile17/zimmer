-- +goose Up
CREATE TABLE classes (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        name TEXT NOT NULL,
        teacher_id UUID NOT NULL references users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE classes;
