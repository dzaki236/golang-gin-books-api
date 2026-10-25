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
		_fatal(fmt.Errorf("gagal memuat berkas .env: %w", err))
	}

	migrationPath := flag.String("path", "", "lokasi berkas migrasi SQL atau direktori yang berisi berkas migrasi SQL")
	databaseURL := flag.String("database-url", os.Getenv("DATABASE_URL"), "URL koneksi PostgreSQL (nilai bawaan dari DATABASE_URL)")
	flag.Parse()

	if err := _validateMigrationPath(*migrationPath); err != nil {
		_fatal(err)
	}
	if *databaseURL == "" {
		_fatal(errors.New("URL basis data wajib diisi; atur DATABASE_URL atau gunakan --database-url"))
	}

	database, err := sql.Open("postgres", *databaseURL)
	if err != nil {
		_fatal(fmt.Errorf("gagal membuka koneksi basis data: %w", err))
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		_fatal(fmt.Errorf("gagal memeriksa koneksi basis data: %w", err))
	}

	pathInfo, err := os.Stat(filepath.Clean(*migrationPath))
	if err != nil {
		_fatal(fmt.Errorf("gagal memeriksa lokasi migrasi: %w", err))
	}

	var source migrate.MigrationSource
	if pathInfo.IsDir() {
		source = &migrate.FileMigrationSource{Dir: *migrationPath}
	} else {
		file, err := os.Open(*migrationPath)
		if err != nil {
			_fatal(fmt.Errorf("gagal membuka berkas migrasi: %w", err))
		}

		migration, err := migrate.ParseMigration(filepath.Base(*migrationPath), file)
		file.Close()
		if err != nil {
			_fatal(fmt.Errorf("gagal membaca berkas migrasi: %w", err))
		}

		source = &migrate.MemoryMigrationSource{Migrations: []*migrate.Migration{migration}}
	}
	count, err := migrate.Exec(database, "postgres", source, migrate.Up)
	if err != nil {
		_fatal(fmt.Errorf("gagal menjalankan migrasi: %w", err))
	}

	fmt.Printf("%d migrasi berhasil diterapkan\n", count)
}

func _validateMigrationPath(path string) error {
	if path == "" {
		return errors.New("lokasi migrasi wajib diisi; gunakan --path=<file.sql atau direktori>")
	}

	info, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("gagal memeriksa lokasi migrasi: %w", err)
	}
	if !info.IsDir() && !strings.EqualFold(filepath.Ext(path), ".sql") {
		return fmt.Errorf("lokasi migrasi %q harus berupa direktori atau berkas .sql", path)
	}

	return nil
}

func _fatal(err error) {
	fmt.Fprintln(os.Stderr, "migrasi gagal:", err)
	os.Exit(1)
}

func main() {
	initialize()
}
