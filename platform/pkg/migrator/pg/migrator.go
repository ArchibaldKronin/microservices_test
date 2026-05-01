package pg

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

type migrator struct {
	db           *sql.DB
	migrationDir string
}

func NewMigrator(db *sql.DB, migrationDir string) *migrator {
	return &migrator{
		db:           db,
		migrationDir: migrationDir,
	}
}

func (m *migrator) Up(ctx context.Context) error {
	err := goose.UpContext(ctx, m.db, m.migrationDir)
	if err != nil {
		return err
	}

	return nil
}

func (m *migrator) Down(ctx context.Context) error {
	err := goose.DownContext(ctx, m.db, m.migrationDir)
	if err != nil {
		return err
	}

	return nil
}
