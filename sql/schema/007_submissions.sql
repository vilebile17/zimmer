-- +goose Up
CREATE TABLE submissions (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        assignment_id UUID NOT NULL references assignments(id) ON DELETE CASCADE,
        user_id UUID NOT NULL UNIQUE references users(id) ON DELETE CASCADE,
        answers TEXT
);

-- +goose Down
DROP TABLE submissions;
