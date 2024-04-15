package db

import (
	"database/sql"

	"github.com/authly/internal/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func MustConnect(dbUrl string) *database.Queries {
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)

	return dbQueries
}
