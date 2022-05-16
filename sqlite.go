package junkboy

import (
	"database/sql"
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dns required")
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
