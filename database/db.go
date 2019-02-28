package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes a new SQLite DB
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		return nil, err
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db, nil
}