package migration

import "embed"

// FS contains embedded SQL migration files.
//
// This allows distributing a single binary (Homebrew / GitHub Releases) without
// shipping migration files next to the executable.
//
//go:embed *.sql
var FS embed.FS
