package testutil

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"goshop/pkg/dbs"
)

// MigrationsDir returns the absolute path to the repo-root `migrations/` directory.
// Resolved relative to this source file so callers don't depend on os.Getwd().
func MigrationsDir() string {
	_, here, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(here), "..", "..", "migrations")
}

// ApplyMigrations runs every up-migration in `migrations/` against the given database. Used
// by integration suites to give the test schema parity with what production applies.
func ApplyMigrations(db dbs.Database) error {
	sqlDB, err := db.GetDB().DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+MigrationsDir(), "postgres", driver)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}
