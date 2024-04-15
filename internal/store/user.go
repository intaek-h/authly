package store

import (
	"github.com/authly/internal/database"
)

type UserStore struct {
	db *database.Queries
}

type UserStoreParams struct {
	DB *database.Queries
}

func NewUserStore(params UserStoreParams) *UserStore {
	return &UserStore{
		db: params.DB,
	}
}
