-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    real_name TEXT,
    nickname TEXT,
    email TEXT NOT NULL,
    profile_image TEXT
);

-- +goose Down
DROP TABLE users;
