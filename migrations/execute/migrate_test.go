package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateMigrationPath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantError bool
	}{
		{name: "empty path", path: "", wantError: true},
		{name: "missing path", path: filepath.Join(t.TempDir(), "missing"), wantError: true},
		{name: "sql file path", path: filepath.Join(t.TempDir(), "migration.sql"), wantError: false},
		{name: "non sql file path", path: filepath.Join(t.TempDir(), "migration.txt"), wantError: true},
		{name: "directory path", path: t.TempDir(), wantError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if filepath.Ext(tt.path) != "" {
				if err := os.WriteFile(tt.path, []byte("-- +migrate Up\nSELECT 1;\n-- +migrate Down\nSELECT 1;\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			err := validateMigrationPath(tt.path)
			if (err != nil) != tt.wantError {
				t.Fatalf("validateMigrationPath(%q) error = %v, wantError = %v", tt.path, err, tt.wantError)
			}
		})
	}
}
