-- +goose Up
CREATE TABLE users_classes (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL references users(id) ON DELETE CASCADE,
        class_id UUID NOT NULL references classes(id) ON DELETE CASCADE,
        joined_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users_classes;
