-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    real_name TEXT,
    nickname TEXT,
    email TEXT NOT NULL,
    profile_image TEXT,
    foreign key (session_id) references sessions (id)
);

-- +goose Down
DROP TABLE users;
