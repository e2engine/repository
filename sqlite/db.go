package sqlite

import (
	"context"
	"database/sql"
	nativeerrors "errors"
	"os"
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

// Open opens a SQLite database connection, configures the connection pool,
// and runs migrations. The caller must prepare the database directory first.
func Open(ctx context.Context, path string, cfg *config.Config) (*sql.DB, error) {
	dbConn, err := openDB(path, cfg)
	if err != nil {
		return nil, errorc.WithCause(errors.ErrCannotOpenDBConnection, err)
	}

	configureDB(dbConn, cfg)

	if errConfigure := configureSQLite(ctx, dbConn); errConfigure != nil {
		return nil, nativeerrors.Join(
			errorc.WithCause(errors.ErrCannotConfigureDB, errConfigure),
			dbConn.Close(),
		)
	}

	if errRun := migrate.RunDBMigration(ctx, dbConn); errRun != nil {
		return nil, nativeerrors.Join(
			errorc.WithCause(errors.ErrCannotMigrateDB, errRun),
			dbConn.Close(),
		)
	}

	return dbConn, nil
}

func openDB(path string, cfg *config.Config) (*sql.DB, error) {
	if path == "" {
		return nil, errors.ErrEmptyDBPath
	}
	if strings.ContainsRune(path, '?') || strings.HasPrefix(path, "file:") {
		return nil, errors.ErrUnsafeDBPath
	}
	if cfg.Driver != "sqlite" {
		return nil, errorc.With(
			errors.ErrUnsupportedDBDriver,
			errorc.String(keys.DBDriver, cfg.Driver),
		)
	}

	// Reject an existing non-regular path before opening the connection pool.
	// This check does not prevent changes to the path after it returns.
	if info, err := os.Lstat(path); err == nil {
		if !info.Mode().IsRegular() {
			return nil, errors.ErrUnsafeDBPath
		}
	} else if !os.IsNotExist(err) {
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
