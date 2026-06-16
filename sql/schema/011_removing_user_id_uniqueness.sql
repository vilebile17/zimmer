-- +goose Up
ALTER TABLE submissions
DROP CONSTRAINT submissions_user_id_key;

-- +goose Down
ALTER TABLE submissions
ADD CONSTRAINT submissions_user_id_key
UNIQUE (user_id);
