-- name: CreateSession :one
insert into
    sessions (user_id, expires_at, created_at, updated_at)
values
    (?, ?, ?, ?) returning *;

-- name: GetSession :one
select
    *
from
    sessions
where
    id = ?;

-- name: GetUserFromSession :one
select
    users.*,
    sessions.expires_at as session_expires_at
from
    sessions
    join users on sessions.user_id = users.id
where
    sessions.id = ?;