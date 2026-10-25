package configs

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func OpenDatabase() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL wajib diisi")
	}
	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("gagal membuka koneksi basis data: %w", err)
	}
	if err := database.Ping(); err != nil {
		database.Close()
		return fmt.Errorf("gagal memeriksa koneksi basis data: %w", err)
	}
	DB = database
	return nil
}
