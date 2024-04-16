-- name: GetUserById :one
select * from users where id = ?;

-- name: GetUserByEmail :one
select * from users where email = ?;

-- name: CreateUser :one
insert into users (
    created_at,
    updated_at,
    real_name,
    nickname,
    email,
    profile_image
) values (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
) returning *;
