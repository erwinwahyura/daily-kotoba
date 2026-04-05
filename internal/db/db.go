// Package db provides database abstraction for PostgreSQL and SQLite
db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DB wraps sql.DB with driver-aware methods
type DB struct {
	*sql.DB
	Driver string // "postgres" or "sqlite"
}

// New wraps an existing sql.DB
func New(db *sql.DB, driver string) *DB {
	return &DB{DB: db, Driver: driver}
}

// GenerateUUID returns a UUID string appropriate for the database
func (db *DB) GenerateUUID() string {
	return uuid.New().String()
}

// JSONValue converts a value to database-appropriate JSON
func (db *DB) JSONValue(v interface{}) (interface{}, error) {
	if v == nil {
		if db.Driver == "postgres" {
			return nil, nil
		}
		return "[]", nil
	}
	
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	
	if db.Driver == "postgres" {
		return b, nil // PostgreSQL accepts []byte for JSONB
	}
	return string(b), nil // SQLite stores as TEXT
}

// ScanJSON scans a JSON value from the database
func (db *DB) ScanJSON(src interface{}, dest interface{}) error {
	if src == nil {
		return nil
	}
	
	var bytes []byte
	switch v := src.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into JSON", src)
	}
	
	if len(bytes) == 0 {
		return nil
	}
	
	return json.Unmarshal(bytes, dest)
}

// Now returns current time in appropriate format
func (db *DB) Now() interface{} {
	now := time.Now()
	if db.Driver == "postgres" {
		return now
	}
	return now.Format(time.RFC3339)
}

// Date returns a date in appropriate format
func (db *DB) Date(t time.Time) interface{} {
	if db.Driver == "postgres" {
		return t
	}
	return t.Format("2006-01-02")
}

// Placeholder returns the appropriate placeholder for the driver
// $1, $2... for PostgreSQL, ?, ?... for SQLite
func (db *DB) Placeholder(n int) string {
	if db.Driver == "postgres" {
		return fmt.Sprintf("$%d", n)
	}
	return "?"
}

// Placeholders returns a slice of placeholders for the given count
func (db *DB) Placeholders(count int) []string {
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = db.Placeholder(i + 1)
	}
	return result
}

// AutoIncrement returns the auto-increment keyword
func (db *DB) AutoIncrement() string {
	if db.Driver == "postgres" {
		return "SERIAL"
	}
	return "INTEGER PRIMARY KEY AUTOINCREMENT"
}

// CurrentTimestamp returns the current timestamp SQL
func (db *DB) CurrentTimestamp() string {
	if db.Driver == "postgres" {
		return "CURRENT_TIMESTAMP"
	}
	return "CURRENT_TIMESTAMP"
}

// LimitOffset returns LIMIT/OFFSET clause
func (db *DB) LimitOffset(limit, offset int) string {
	if db.Driver == "postgres" {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

// InitializeSQLite sets up SQLite with required PRAGMAs
func (db *DB) InitializeSQLite() error {
	if db.Driver != "sqlite" {
		return nil
	}
	
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
	`)
	return err
}
