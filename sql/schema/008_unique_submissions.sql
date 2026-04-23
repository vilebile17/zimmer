-- +goose Up
ALTER TABLE submissions
ADD CONSTRAINT unique_student_submission
UNIQUE (assignment_id, user_id);

-- +goose Down
ALTER TABLE submissions
DROP CONSTRAINT unique_student_submission;
