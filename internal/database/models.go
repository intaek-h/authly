// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
	"time"
)

type Session struct {
	ID        int64
	UserID    int64
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID           int64
	SessionID    sql.NullInt64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RealName     sql.NullString
	Nickname     sql.NullString
	Email        string
	ProfileImage sql.NullString
}
