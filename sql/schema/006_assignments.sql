-- +goose Up
CREATE TABLE assignments (
        id UUID PRIMARY KEY,
        class_id UUID NOT NULL references classes(id) ON DELETE CASCADE,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        due_at TIMESTAMP,
        title TEXT NOT NULL,
        instructions TEXT,
        allow_late BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE assignments;
