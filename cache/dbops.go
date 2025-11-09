package cache

import (
	"context"
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

func CreateNewDb(dbPath string) (*sql.DB, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return db, nil
}
