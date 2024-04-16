-- name: CreateSession :one
insert into sessions (
    user_id,
    expires_at,
    created_at,
    updated_at
) values (
    ?,
    ?,
    ?,
    ?
) returning id;