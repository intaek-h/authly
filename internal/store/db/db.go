package db

import "database/sql"

func MustConnect(dbUrl string) *sql.DB {
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		panic(err)
	}

	return db
}
