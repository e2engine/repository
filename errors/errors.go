package errors

import "github.com/ygrebnov/errorc"

var (
	ErrUnsupportedDBDriver    = errorc.New("unsupported database driver")
	ErrCannotOpenDBConnection = errorc.New("cannot open db connection")
	ErrCannotConfigureDB      = errorc.New("cannot configure database")
	ErrCannotMigrateDB        = errorc.New("cannot migrate database")
	ErrEmptyDBPath            = errorc.New("db path is empty")
	ErrUnsafeDBPath           = errorc.New("unsafe db path")
	ErrInvalidDBFilename      = errorc.New("invalid db filename")
	ErrFailedToEnableWALMode  = errorc.New("failed to enable WAL mode")
)
