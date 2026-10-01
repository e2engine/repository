package errors

import "github.com/ygrebnov/errorc"

var (
	ErrCannotOpenDBConnection = errorc.New("cannot open db connection")
	ErrCannotConfigureDB      = errorc.New("cannot configure database")
	ErrCannotMigrateDB        = errorc.New("cannot migrate database")
	ErrEmptyDBPath            = errorc.New("db path is empty")
	ErrFailedToEnableWALMode  = errorc.New("failed to enable WAL mode")
)
