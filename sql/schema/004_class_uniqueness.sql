-- +goose Up
ALTER TABLE classes
ADD CONSTRAINT unique_user_class_name
UNIQUE (teacher_id, name);

-- +goose Down
ALTER TABLE classes
DROP CONSTRAINT unique_user_class_name;
