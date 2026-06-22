-- +goose Up
CREATE TABLE resources (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        class_id UUID NOT NULL references classes(id) ON DELETE CASCADE,
        title TEXT NOT NULL,
        content TEXT
);

-- +goose Down
DROP TABLE resources;
