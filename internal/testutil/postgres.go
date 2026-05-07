// Package testutil contains helpers shared across integration tests. The Postgres helper
// spins up a fresh container per test invocation via testcontainers-go so suites get a
// clean schema with no shared state.
package testutil

import (
	"context"
	"testing"
	"time"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"goshop/pkg/dbs"
)

// StartPostgres boots a throwaway Postgres 16 container, waits for it to accept connections,
// and returns a wired dbs.Database. The caller is responsible for cleaning up via t.Cleanup.
// Tests that run without a Docker daemon should call SkipIfNoDocker first.
func StartPostgres(ctx context.Context, t *testing.T) dbs.Database {
	t.Helper()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("goshop_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.BasicWaitStrategies(),
		tcpostgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	t.Cleanup(func() {
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = container.Terminate(shutdown)
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("postgres dsn: %v", err)
	}

	// Wait for psql readiness on top of testcontainers' wait-for-log to handle initdb races.
	deadline := time.Now().Add(15 * time.Second)
	var db dbs.Database
	for {
		db, err = dbs.NewDatabase(dsn)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("connect to postgres: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
	}
	return db
}
