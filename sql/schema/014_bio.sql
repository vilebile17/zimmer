-- +goose Up
ALTER TABLE users
ADD COLUMN bio TEXT NOT NULL
DEFAULT 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec non tempus diam. Integer bibendum odio nec tristique varius. Praesent euismod tempus urna, eget hendrerit mi dapibus ac. Quisque sodales porttitor.';

-- +goose Down
ALTER TABLE users
DROP COLUMN bio;
