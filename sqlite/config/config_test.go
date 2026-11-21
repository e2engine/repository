package config

import (
	"errors"
	"strings"
	"testing"

	modellib "github.com/ygrebnov/model"

	repositoryerrors "github.com/e2engine/repository/errors"
)

func TestValidateDBFilename(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{name: "abc", valid: true},
		{name: "e2engine.sqlite", valid: true},
		{name: "A_2-backup.sqlite", valid: true},
		{name: strings.Repeat("a", 255), valid: true},
		{name: ""},
		{name: "ab"},
		{name: strings.Repeat("a", 256)},
		{name: ".db"},
		{name: "db."},
		{name: "../outside.sqlite"},
		{name: `..\outside.sqlite`},
		{name: "/tmp/outside.sqlite"},
		{name: "data.sqlite?mode=memory"},
		{name: "data.sqlite#fragment"},
		{name: "data%2esqlite"},
		{name: "data base.sqlite"},
		{name: "é.sqlite"},
		{name: "CON.sqlite"},
		{name: "aux.sqlite"},
		{name: "COM1.sqlite"},
		{name: "lpt9.sqlite"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDBFilename(tt.name)
			if tt.valid {
				if err != nil {
					t.Fatalf("validateDBFilename(%q) = %v", tt.name, err)
				}
				return
			}
			if !errors.Is(err, repositoryerrors.ErrInvalidDBFilename) {
				t.Fatalf("validateDBFilename(%q) = %v, want ErrInvalidDBFilename", tt.name, err)
			}
		})
	}
}

func TestDBFilenameValidationRule(t *testing.T) {
	binding, err := modellib.NewBinding[Config](modellib.WithRules(GetValidationRules()...))
	if err != nil {
		t.Fatal(err)
	}
	valid := Config{Name: "e2engine.sqlite", Driver: "sqlite", MaxOpenConns: 1, MaxIdleConns: 1}
	if err := binding.Validate(t.Context(), &valid); err != nil {
		t.Fatalf("valid config: %v", err)
	}
	valid.Name = "../outside.sqlite"
	if err := binding.Validate(t.Context(), &valid); err == nil {
		t.Fatal("binding accepted invalid db filename")
	}
}
