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
		return fmt.Errorf("DATABASE_URL is required")
	}
	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("gagak database: %w", err)
	}
	if err := database.Ping(); err != nil {
		database.Close()
		return fmt.Errorf("ping database: %w", err)
	}
	DB = database
	return nil
}
