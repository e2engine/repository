package config

import (
	"strings"

	"github.com/ygrebnov/errorc"
	modellib "github.com/ygrebnov/model"

	"github.com/e2engine/repository/errors"
	"github.com/e2engine/repository/keys"
)

type Config struct {
	Name         string `yaml:"name" default:"e2engine.sqlite" validate:"dbfilename"`
	Driver       string `yaml:"driver" default:"sqlite" validate:"oneof(sqlite)"`
	MaxOpenConns int    `yaml:"max_open_conns" default:"1" validate:"min(1)"`
	MaxIdleConns int    `yaml:"max_idle_conns" default:"1" validate:"min(1)"`
}

// GetValidationRules returns a fresh set of database config rules.
func GetValidationRules() []modellib.Rule {
	return []modellib.Rule{getDBFilenameRule()}
}

func getDBFilenameRule() modellib.Rule {
	return modellib.NewRule[string]("dbfilename", func(name string, _ ...string) error {
		return validateDBFilename(name)
	})
}

// validateDBFilename keeps the database filename within its configured data directory
// and prevents filename characters from being interpreted as SQLite DSN options.
func validateDBFilename(name string) error {
	if len(name) < 3 || len(name) > 255 {
		return errorc.With(
			errors.ErrInvalidDBFilename,
			errorc.String(keys.Filename, name),
		)
	}

	if strings.HasSuffix(name, ".") {
		return errorc.With(
			errors.ErrInvalidDBFilename,
			errorc.String(keys.Filename, name),
		)
	}

	for i, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || (i > 0 && (r == '.' || r == '_' || r == '-')) {
			continue
		}
		return errorc.With(
			errors.ErrInvalidDBFilename,
			errorc.String(keys.Filename, name),
		)
	}

	base := strings.ToUpper(strings.SplitN(name, ".", 2)[0])
	switch base {
	case "CON", "PRN", "AUX", "NUL":
		return errorc.With(
			errors.ErrInvalidDBFilename,
			errorc.String(keys.Filename, name),
		)
	}

	if len(base) == 4 &&
		(strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) &&
		base[3] >= '1' && base[3] <= '9' {
		return errorc.With(
			errors.ErrInvalidDBFilename,
			errorc.String(keys.Filename, name),
		)
	}
	return nil
}
