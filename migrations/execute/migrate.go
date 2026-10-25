package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

// https://github.com/golang-migrate/migrate : Forked And Reference from
func initialize() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		_fatal(fmt.Errorf("load .env: %w", err))
	}

	migrationPath := flag.String("path", "", "path to a SQL migration file or directory containing SQL migration files")
	databaseURL := flag.String("database-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection URL (defaults to DATABASE_URL)")
	flag.Parse()

	if err := _validateMigrationPath(*migrationPath); err != nil {
		_fatal(err)
	}
	if *databaseURL == "" {
		_fatal(errors.New("database URL is required; set DATABASE_URL or pass --database-url"))
	}

	database, err := sql.Open("postgres", *databaseURL)
	if err != nil {
		_fatal(fmt.Errorf("open database: %w", err))
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		_fatal(fmt.Errorf("ping database: %w", err))
	}

	pathInfo, err := os.Stat(filepath.Clean(*migrationPath))
	if err != nil {
		_fatal(fmt.Errorf("stat migration path: %w", err))
	}

	var source migrate.MigrationSource
	if pathInfo.IsDir() {
		source = &migrate.FileMigrationSource{Dir: *migrationPath}
	} else {
		file, err := os.Open(*migrationPath)
		if err != nil {
			_fatal(fmt.Errorf("open migration file: %w", err))
		}

		migration, err := migrate.ParseMigration(filepath.Base(*migrationPath), file)
		file.Close()
		if err != nil {
			_fatal(fmt.Errorf("parse migration file: %w", err))
		}

		source = &migrate.MemoryMigrationSource{Migrations: []*migrate.Migration{migration}}
	}
	count, err := migrate.Exec(database, "postgres", source, migrate.Up)
	if err != nil {
		_fatal(fmt.Errorf("run migrations: %w", err))
	}

	fmt.Printf("applied %d migration(s)\n", count)
}

func _validateMigrationPath(path string) error {
	if path == "" {
		return errors.New("migration path is required; pass --path=<file.sql or directory>")
	}

	info, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("stat migration path: %w", err)
	}
	if !info.IsDir() && !strings.EqualFold(filepath.Ext(path), ".sql") {
		return fmt.Errorf("migration path %q must be a directory or a .sql file", path)
	}

	return nil
}

func _fatal(err error) {
	fmt.Fprintln(os.Stderr, "migration failed:", err)
	os.Exit(1)
}

func main() {
	initialize()
}
