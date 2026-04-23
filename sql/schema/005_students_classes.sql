-- +goose Up
CREATE TABLE students_classes (
        id UUID PRIMARY KEY,
        joined_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        student_id UUID NOT NULL references users(id) ON DELETE CASCADE,
        class_id UUID NOT NULL references classes(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE students_classes;
