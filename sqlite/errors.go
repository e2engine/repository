package sqlite

import (
	"database/sql"
	"errors"

	"github.com/e2engine/core/repository"
)

func mapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrNotFound
	}

	var sqliteErr sqliteErrorCoder
	if !errors.As(err, &sqliteErr) {
		return err
	}

	//nolint:gocritic // SQLite symbolic names document the numeric extended error codes.
	switch sqliteErr.Code() {
	case 1555, 2067:
		// SQLITE_CONSTRAINT_PRIMARYKEY
		// SQLITE_CONSTRAINT_UNIQUE
		return repository.ErrAlreadyExists

	case 275:
		// SQLITE_CONSTRAINT_CHECK
		return repository.ErrInvalidPayload

	case 787:
		// SQLITE_CONSTRAINT_FOREIGNKEY
		return repository.ErrConflict

	default:
		return err
	}
}

type sqliteErrorCoder interface {
	Code() int
}
