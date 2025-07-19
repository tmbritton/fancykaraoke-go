package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/knaka/go-sqlite3-fts5" // Full-text search SQLite extension
	_ "github.com/mattn/go-sqlite3"      // SQLite driver
)

// SQLiteStore handles database connections and operations
type SQLiteStore struct {
	DB *sql.DB
}

var (
	conn *SQLiteStore
	once sync.Once
)

func GetConnection() *SQLiteStore {
	once.Do(func() {
		db, err := initDb()
		if err != nil {
			log.Fatal(err)
		}
		conn = db
	})
	return conn
}

func initDb() (*SQLiteStore, error) {
	dbLocation := "./fancykaraoke.db?fts5=true&foreign_keys=true"

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
	return &SQLiteStore{DB: db}, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.DB.Close()
}
