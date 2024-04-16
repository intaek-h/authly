// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :one
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
) returning id, created_at, updated_at, real_name, nickname, email, profile_image
`

type CreateUserParams struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RealName     sql.NullString
	Nickname     sql.NullString
	Email        string
	ProfileImage sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.RealName,
		arg.Nickname,
		arg.Email,
		arg.ProfileImage,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RealName,
		&i.Nickname,
		&i.Email,
		&i.ProfileImage,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, created_at, updated_at, real_name, nickname, email, profile_image from users where email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RealName,
		&i.Nickname,
		&i.Email,
		&i.ProfileImage,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
select id, created_at, updated_at, real_name, nickname, email, profile_image from users where id = ?
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RealName,
		&i.Nickname,
		&i.Email,
		&i.ProfileImage,
	)
	return i, err
}
