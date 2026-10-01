package migrate

import (
	"context"
	"database/sql"
	"errors"

	migratepkg "github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/e2engine/repository/sqlite/migration"
)

// RunDBMigration executes migrations that have not yet been applied.
func RunDBMigration(
	ctx context.Context,
	db *sql.DB,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	source, err := iofs.New(
		migration.FS,
		".",
	)
	if err != nil {
		return err
	}

	driver, err := migratedb.WithInstance(
		db,
		&migratedb.Config{},
	)
	if err != nil {
		return err
	}

	m, err := migratepkg.NewWithInstance(
		"iofs",
		source,
		"sqlite",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil &&
		!errors.Is(err, migratepkg.ErrNoChange) {
		return err
	}

	return nil
}
