package migrate

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRunDBMigration(t *testing.T) {
	db, err := sql.Open(
		"sqlite",
		filepath.Join(t.TempDir(), "test.sqlite"),
	)
	if err != nil {
		t.Fatalf("Expected database to open, got %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := RunDBMigration(t.Context(), db); err != nil {
		t.Fatalf("Expected migration to succeed, got %v", err)
	}

	for _, table := range []string{
		"resources",
		"test_executions",
		"test_suite_executions",
	} {
		t.Run(table, func(t *testing.T) {
			var exists int

			if err := db.QueryRowContext(
				t.Context(),
				`
					SELECT COUNT(*)
					FROM sqlite_master
					WHERE type = 'table'
					  AND name = ?
				`,
				table,
			).Scan(&exists); err != nil {
				t.Fatalf("Expected table lookup to succeed, got %v", err)
			}

			if exists != 1 {
				t.Fatalf("Expected table %q to exist", table)
			}
		})
	}

	// Running migrations again should result in ErrNoChange internally,
	// which RunDBMigration treats as success.
	if err := RunDBMigration(t.Context(), db); err != nil {
		t.Fatalf("Expected repeated migration to succeed, got %v", err)
	}
}

func TestRunDBMigrationCanceledContext(t *testing.T) {
	db, err := sql.Open(
		"sqlite",
		filepath.Join(t.TempDir(), "test.sqlite"),
	)
	if err != nil {
		t.Fatalf("Expected database to open, got %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err = RunDBMigration(ctx, db)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled, got %v", err)
	}
}
