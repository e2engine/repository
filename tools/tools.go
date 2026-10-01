//go:build tools

// This file is for managing Go programs version with `go.mod`, which allows
// them to be kept up-to-date through tools like Dependabot.

package tools

import (
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
