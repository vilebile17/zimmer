-- +goose Up
DROP TABLE resources;
CREATE TABLE class_content (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        content_type TEXT NOT NULL,
        class_id UUID NOT NULL references classes(id) ON DELETE CASCADE,
        title TEXT NOT NULL,
        content TEXT,

        CONSTRAINT check_content_type
                CHECK (content_type IN ('announcement', 'resource'))
);

-- +goose Down
CREATE TABLE resources (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        class_id UUID NOT NULL references classes(id) ON DELETE CASCADE,
        title TEXT NOT NULL,
        content TEXT
);
DROP TABLE class_content;
