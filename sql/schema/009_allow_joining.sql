-- +goose Up
ALTER TABLE classes
ADD COLUMN allow_joining BOOLEAN NOT NULL
DEFAULT true;

-- +goose Down
ALTER TABLE classes
DROP COLUMN allow_joining;
