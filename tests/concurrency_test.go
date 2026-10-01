package tests

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/e2engine/repository/sqlite"
	"github.com/e2engine/repository/sqlite/config"
)

func insertResource(
	ctx context.Context,
	exec interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	},
	id string,
) error {
	_, err := exec.ExecContext(
		ctx,
		`
			INSERT INTO resources (
				id,
				kind,
				version,
				name,
				description,
				spec,
				created_at,
				updated_at
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
		id,
		"environment",
		"1.0.0",
		id,
		"",
		[]byte(`{}`),
		time.Now().UTC(),
		time.Now().UTC(),
	)

	return err
}

func TestSQLiteConcurrentReadWrite(t *testing.T) { // WAL semantics
	path := filepath.Join(t.TempDir(), "test.sqlite")

	cfg := &config.Config{
		Driver:       "sqlite",
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}

	db1, err := sqlite.Open(t.Context(), path, cfg)
	if err != nil {
		t.Fatalf("Expected first database to open, got %v", err)
	}
	t.Cleanup(func() {
		_ = db1.Close()
	})

	db2, err := sqlite.Open(t.Context(), path, cfg)
	if err != nil {
		t.Fatalf("Expected second database to open, got %v", err)
	}
	t.Cleanup(func() {
		_ = db2.Close()
	})

	if err := insertResource(t.Context(), db1, "resource-1"); err != nil {
		t.Fatalf("Expected resource to be created, got %v", err)
	}

	tx, err := db1.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("Expected read transaction, got %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var id string

	if err := tx.QueryRowContext(
		t.Context(),
		"SELECT id FROM resources WHERE id = ?",
		"resource-1",
	).Scan(&id); err != nil {
		t.Fatalf("Expected resource to be read, got %v", err)
	}

	// The read transaction remains active here.
	if err := insertResource(t.Context(), db2, "resource-2"); err != nil {
		t.Fatalf(
			"Expected concurrent write during read transaction, got %v",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("Expected read transaction to commit, got %v", err)
	}

	var count int

	if err := db1.QueryRowContext(
		t.Context(),
		"SELECT COUNT(*) FROM resources WHERE id = ?",
		"resource-2",
	).Scan(&count); err != nil {
		t.Fatalf("Expected resource lookup to succeed, got %v", err)
	}

	if count != 1 {
		t.Fatalf("Expected second resource to exist")
	}
}

func TestSQLiteCompetingWrites(t *testing.T) { // busy_timeout
	path := filepath.Join(t.TempDir(), "test.sqlite")

	cfg := &config.Config{
		Driver:       "sqlite",
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}

	db1, err := sqlite.Open(t.Context(), path, cfg)
	if err != nil {
		t.Fatalf("Expected first database to open, got %v", err)
	}
	t.Cleanup(func() {
		_ = db1.Close()
	})

	db2, err := sqlite.Open(t.Context(), path, cfg)
	if err != nil {
		t.Fatalf("Expected second database to open, got %v", err)
	}
	t.Cleanup(func() {
		_ = db2.Close()
	})

	tx, err := db1.BeginTx(t.Context(), nil)
	if err != nil {
		t.Fatalf("Expected write transaction, got %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := insertResource(
		t.Context(),
		tx,
		"resource-1",
	); err != nil {
		t.Fatalf("Expected first resource to be created, got %v", err)
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- insertResource(
			t.Context(),
			db2,
			"resource-2",
		)
	}()

	// Keep the write lock long enough for db2 to encounter it.
	time.Sleep(100 * time.Millisecond)

	if err := tx.Commit(); err != nil {
		t.Fatalf("Expected first transaction to commit, got %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf(
				"Expected competing write to wait and succeed, got %v",
				err,
			)
		}

	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for competing write")
	}

	var count int

	if err := db1.QueryRowContext(
		t.Context(),
		`
			SELECT COUNT(*)
			FROM resources
			WHERE id IN (?, ?)
		`,
		"resource-1",
		"resource-2",
	).Scan(&count); err != nil {
		t.Fatalf("Expected resource lookup to succeed, got %v", err)
	}

	if count != 2 {
		t.Fatalf("Expected 2 resources, got %d", count)
	}
}
