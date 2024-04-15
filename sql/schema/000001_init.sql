-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    real_name TEXT NOT NULL,
    nickname TEXT NOT NULL,
    email TEXT NOT NULL,
    profile_image TEXT
);

-- +goose Down
DROP TABLE users;
