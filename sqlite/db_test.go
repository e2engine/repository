package sqlite

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	repositoryerrors "github.com/e2engine/repository/errors"
	"github.com/e2engine/repository/sqlite/config"
)

func TestGetPath(t *testing.T) {
	testCases := []struct {
		name     string
		dataPath string
		dbName   string
		expected string
	}{
		{
			name:     "default",
			dataPath: "/tmp/e2engine",
			dbName:   "e2engine.sqlite",
			expected: filepath.Join("/tmp/e2engine", "e2engine.sqlite"),
		},
		{
			name:     "relative path",
			dataPath: "data",
			dbName:   "e2engine.sqlite",
			expected: filepath.Join("data", "e2engine.sqlite"),
		},
		{
			name:     "whitespace",
			dataPath: " data ",
			dbName:   " e2engine.sqlite ",
			expected: strings.TrimSpace(
				filepath.Join(" data ", " e2engine.sqlite "),
			),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := GetPath(
				testCase.dataPath,
				&config.Config{Name: testCase.dbName},
			)

			if actual != testCase.expected {
				t.Fatalf(
					"Expected path %q, got %q",
					testCase.expected,
					actual,
				)
			}
		})
	}
}

func TestOpenDBEmptyPath(t *testing.T) {
	actual, err := openDB("", &config.Config{
		Driver: "sqlite",
	})

	if !errors.Is(err, repositoryerrors.ErrEmptyDBPath) {
		t.Fatalf(
			"Expected error %v, got %v",
			repositoryerrors.ErrEmptyDBPath,
			err,
		)
	}
	if actual != nil {
		t.Fatalf("Expected database to be nil")
	}
}

func TestOpenDB(t *testing.T) {
	testCases := []struct {
		name string
		path func(string) string
	}{
		{
			name: "database in existing directory",
			path: func(dir string) string {
				return filepath.Join(dir, "test.sqlite")
			},
		},
		{
			name: "database in nested directory",
			path: func(dir string) string {
				return filepath.Join(dir, "nested", "data", "test.sqlite")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			path := testCase.path(t.TempDir())

			actual, err := openDB(path, &config.Config{
				Driver: "sqlite",
			})
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if actual == nil {
				t.Fatalf("Expected database")
			}
			t.Cleanup(func() {
				_ = actual.Close()
			})

			if err := actual.PingContext(t.Context()); err != nil {
				t.Fatalf("Expected database ping to succeed, got %v", err)
			}

			info, err := os.Stat(filepath.Dir(path))
			if err != nil {
				t.Fatalf("Expected database directory, got %v", err)
			}
			if !info.IsDir() {
				t.Fatalf("Expected %q to be a directory", filepath.Dir(path))
			}

			var foreignKeys int

			if err := actual.QueryRowContext(
				t.Context(),
				"PRAGMA foreign_keys",
			).Scan(&foreignKeys); err != nil {
				t.Fatalf("Expected foreign_keys pragma, got %v", err)
			}

			if foreignKeys != 1 {
				t.Fatalf(
					"Expected foreign_keys to be enabled, got %d",
					foreignKeys,
				)
			}

			var busyTimeout int
			if err := actual.QueryRowContext(
				t.Context(),
				"PRAGMA busy_timeout",
			).Scan(&busyTimeout); err != nil {
				t.Fatalf("Expected busy_timeout pragma, got %v", err)
			}

			if busyTimeout != 5000 {
				t.Fatalf(
					"Expected busy_timeout 5000, got %d",
					busyTimeout,
				)
			}
		})
	}
}

func TestConfigureDB(t *testing.T) {
	testCases := []struct {
		name         string
		maxOpenConns int
		maxIdleConns int
	}{
		{
			name:         "single connection",
			maxOpenConns: 1,
			maxIdleConns: 1,
		},
		{
			name:         "multiple connections",
			maxOpenConns: 4,
			maxIdleConns: 2,
		},
		{
			name:         "no idle connections",
			maxOpenConns: 8,
			maxIdleConns: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dbConn, err := openDB(
				filepath.Join(t.TempDir(), "test.sqlite"),
				&config.Config{Driver: "sqlite"},
			)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			t.Cleanup(func() {
				_ = dbConn.Close()
			})

			configureDB(dbConn, &config.Config{
				MaxOpenConns: testCase.maxOpenConns,
				MaxIdleConns: testCase.maxIdleConns,
			})

			actual := dbConn.Stats().MaxOpenConnections
			if actual != testCase.maxOpenConns {
				t.Fatalf(
					"Expected max open connections %d, got %d",
					testCase.maxOpenConns,
					actual,
				)
			}
		})
	}
}

func TestConfigureSQLite(t *testing.T) {
	dbConn, err := openDB(
		filepath.Join(t.TempDir(), "test.sqlite"),
		&config.Config{
			Driver: "sqlite",
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Cleanup(func() {
		_ = dbConn.Close()
	})

	if err := configureSQLite(t.Context(), dbConn); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var mode string

	if err := dbConn.QueryRowContext(
		t.Context(),
		"PRAGMA journal_mode",
	).Scan(&mode); err != nil {
		t.Fatalf("Expected journal_mode pragma, got %v", err)
	}

	if !strings.EqualFold(mode, "wal") {
		t.Fatalf("Expected journal mode WAL, got %q", mode)
	}
}

func TestOpen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.sqlite")

	dbConn, err := Open(
		t.Context(),
		path,
		&config.Config{
			Driver:       "sqlite",
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if dbConn == nil {
		t.Fatalf("Expected database")
	}
	t.Cleanup(func() {
		_ = dbConn.Close()
	})

	if err := dbConn.PingContext(t.Context()); err != nil {
		t.Fatalf("Expected database ping to succeed, got %v", err)
	}

	var foreignKeys int

	if err := dbConn.QueryRowContext(
		t.Context(),
		"PRAGMA foreign_keys",
	).Scan(&foreignKeys); err != nil {
		t.Fatalf("Expected foreign_keys pragma, got %v", err)
	}

	if foreignKeys != 1 {
		t.Fatalf(
			"Expected foreign_keys to be enabled, got %d",
			foreignKeys,
		)
	}

	var busyTimeout int

	if err := dbConn.QueryRowContext(
		t.Context(),
		"PRAGMA busy_timeout",
	).Scan(&busyTimeout); err != nil {
		t.Fatalf("Expected busy_timeout pragma, got %v", err)
	}

	if busyTimeout != 5000 {
		t.Fatalf(
			"Expected busy_timeout 5000, got %d",
			busyTimeout,
		)
	}

	var journalMode string

	if err := dbConn.QueryRowContext(
		t.Context(),
		"PRAGMA journal_mode",
	).Scan(&journalMode); err != nil {
		t.Fatalf("Expected journal_mode pragma, got %v", err)
	}

	if !strings.EqualFold(journalMode, "wal") {
		t.Fatalf("Expected journal mode WAL, got %q", journalMode)
	}

	for _, table := range []string{
		"resources",
		"test_executions",
		"test_suite_executions",
	} {
		t.Run(table, func(t *testing.T) {
			var exists int

			err := dbConn.QueryRowContext(
				t.Context(),
				`
					SELECT COUNT(*)
					FROM sqlite_master
					WHERE type = 'table'
					  AND name = ?
				`,
				table,
			).Scan(&exists)
			if err != nil {
				t.Fatalf("Expected table lookup to succeed, got %v", err)
			}

			if exists != 1 {
				t.Fatalf("Expected table %q to exist", table)
			}
		})
	}
}

func TestNewRepositories(t *testing.T) {
	dbConn, err := Open(
		context.Background(),
		filepath.Join(t.TempDir(), "test.sqlite"),
		&config.Config{
			Driver:       "sqlite",
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Cleanup(func() {
		_ = dbConn.Close()
	})

	repositories := NewRepositories(dbConn)

	if repositories.Environment == nil {
		t.Fatalf("Expected environment repository")
	}
	if repositories.Test == nil {
		t.Fatalf("Expected test repository")
	}
	if repositories.TestSuite == nil {
		t.Fatalf("Expected test suite repository")
	}
	if repositories.TestExecution == nil {
		t.Fatalf("Expected test execution repository")
	}
	if repositories.TestSuiteExecution == nil {
		t.Fatalf("Expected test suite execution repository")
	}
}
