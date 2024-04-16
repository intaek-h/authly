-- +goose Up
CREATE TABLE sessions (
    id integer primary key autoincrement,
    user_id integer not null,
    expires_at timestamp not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    foreign key (user_id) references users (id)
);

-- +goose Down
DROP TABLE sessions;
