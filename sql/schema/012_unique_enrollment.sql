-- +goose Up
ALTER TABLE students_classes
ADD CONSTRAINT unique_student_class
UNIQUE (student_id, class_id);

-- +goose Down
ALTER TABLE students_classes
DROP CONSTRAINT unique_student_class;
