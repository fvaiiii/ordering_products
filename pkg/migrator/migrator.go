package migrator

import (
	"database/sql"

	"github.com/pressly/goose"
)

var globalMigrator *migrator

type migrator struct {
	db            *sql.DB
	migrationsDir string
}

func Migrator() *migrator {
	return globalMigrator
}

func Init(db *sql.DB, migrationsDir string) {
	if globalMigrator == nil {
		globalMigrator = &migrator{
			db:            db,
			migrationsDir: migrationsDir,
		}
	}
}

func (m *migrator) Up() error {
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}
	return nil
}

func (m *migrator) Down() error {
	err := goose.Down(m.db, m.migrationsDir)
	if err != nil {
		return err
	}
	return nil
}
