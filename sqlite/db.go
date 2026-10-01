package sqlite

import (
	"context"
	"database/sql"
	nativeerrors "errors"
	"os"
	"path/filepath"
	"strings"

	coreapi "github.com/e2engine/core/api"
	"github.com/ygrebnov/errorc"
	_ "modernc.org/sqlite" // Register the modernc SQLite driver with database/sql.

	"github.com/e2engine/repository/errors"
	"github.com/e2engine/repository/keys"
	"github.com/e2engine/repository/sqlite/config"
	dbpkg "github.com/e2engine/repository/sqlite/internal/sqlc"
	"github.com/e2engine/repository/sqlite/migrate"
)

func GetPath(dataPath string, cfg *config.Config) string {
	// Default location stays under the config directory.
	dbPath := filepath.Join(dataPath, cfg.Name)

	return strings.TrimSpace(dbPath)
}

// Open opens a SQLite database connection, configures the connection pool, and runs migrations.
func Open(ctx context.Context, path string, cfg *config.Config) (*sql.DB, error) {
	dbConn, err := openDB(path, cfg)
	if err != nil {
		return nil, nativeerrors.Join(errors.ErrCannotOpenDBConnection, err)
	}

	configureDB(dbConn, cfg)

	if err := configureSQLite(ctx, dbConn); err != nil {
		return nil, nativeerrors.Join(
			errors.ErrCannotConfigureDB,
			err,
			dbConn.Close(),
		)
	}

	if err := migrate.RunDBMigration(ctx, dbConn); err != nil {
		return nil, nativeerrors.Join(
			errors.ErrCannotMigrateDB,
			err,
			dbConn.Close(),
		)
	}

	return dbConn, nil
}

func openDB(path string, cfg *config.Config) (*sql.DB, error) {
	if path == "" {
		return nil, errors.ErrEmptyDBPath
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	sqldb, err := sql.Open(
		cfg.Driver,
		path+"?_pragma=foreign_keys=on&_pragma=busy_timeout(5000)",
	)
	if err != nil {
		return nil, err
	}

	return sqldb, nil
}

func configureDB(sqldb *sql.DB, cfg *config.Config) {
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
}

func configureSQLite(ctx context.Context, db *sql.DB) error {
	var mode string

	if err := db.QueryRowContext(
		ctx,
		"PRAGMA journal_mode=WAL",
	).Scan(&mode); err != nil {
		return err
	}

	if !strings.EqualFold(mode, "wal") {
		return errorc.With(
			errors.ErrFailedToEnableWALMode,
			errorc.String(keys.JournalMode, mode),
		)
	}

	return nil
}

func NewRepositories(dbConn *sql.DB) coreapi.Repositories {
	// Create store.
	store := dbpkg.NewStore(dbConn)

	return coreapi.Repositories{
		Environment:        newEnvironmentRepository(store),
		Test:               newTestRepository(store),
		TestSuite:          newTestSuiteRepository(store),
		TestExecution:      newTestExecutionRepository(store),
		TestSuiteExecution: newTestSuiteExecutionRepository(store),
	}
}
