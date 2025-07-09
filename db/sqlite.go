package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLiteStore handles database connections and operations
type SQLiteStore struct {
	db *sql.DB
}

func InitDb() (*SQLiteStore, error) {
	dbLocation := "./fancykaraoke.db"

	// Open database connection
	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &SQLiteStore{db: db}, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
