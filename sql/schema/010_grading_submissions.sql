-- +goose Up
ALTER TABLE submissions
ADD COLUMN score INTEGER;

-- +goose Down
ALTER TABLE submissions
DROP COLUMN score;
