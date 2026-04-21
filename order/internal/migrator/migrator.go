package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db            *sql.DB
	maigrationDir string
}

func NewMigrator(db *sql.DB, migrationDir string) *Migrator {
	return &Migrator{
		db:            db,
		maigrationDir: migrationDir,
	}
}

func (m *Migrator) Up() error {
	err := goose.Up(m.db, m.maigrationDir)
	if err != nil {
		return err
	}

	return nil
}
