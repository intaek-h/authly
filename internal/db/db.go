package db

import (
	"database/sql"

	"github.com/authly/internal/database"
	"github.com/authly/internal/env"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func MustConnect(env *env.Env) *database.Queries {
	if env.Environment == "development" {
		return MustConnectLocal(env.DatabaseUrl)
	}

	db, err := sql.Open("libsql", env.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)

	return dbQueries
}

func MustConnectLocal(dbUrl string) *database.Queries {
	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)

	return dbQueries
}
