package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func RunMigrations(db *sql.DB) error {
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return err
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		return err
	}

	return nil
}
