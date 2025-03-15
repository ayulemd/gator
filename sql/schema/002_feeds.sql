-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT,
    url TEXT UNIQUE,
    user_id UUID,
    CONSTRAINT fk_user
    FOREIGN KEY (user_id)
    REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
