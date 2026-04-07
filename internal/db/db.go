// Package db provides database abstraction for PostgreSQL and SQLite
package db

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
// Handles both raw JSON strings (from seed files) and parsed Go types
func (db *DB) JSONValue(v interface{}) (interface{}, error) {
	if v == nil {
		if db.Driver == "postgres" {
			return nil, nil
		}
		return "[]", nil
	}
	
	// If it's already a string, assume it's pre-formatted JSON from seed files
	// Validate it's valid JSON and return as-is
	if s, ok := v.(string); ok {
		var test interface{}
		if err := json.Unmarshal([]byte(s), &test); err != nil {
			// Invalid JSON, wrap it in an array
			return fmt.Sprintf("[%q]", s), nil
		}
		// Valid JSON - return as string for SQLite, []byte for Postgres
		if db.Driver == "postgres" {
			return []byte(s), nil
		}
		return s, nil
	}
	
	// For other types (slices, maps from parsed JSON), marshal them
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

// InitializeSQLite sets up SQLite with performance PRAGMAs
// Based on: https://phiresky.github.io/blog/2020/sqlite-performance-tuning/
func (db *DB) InitializeSQLite() error {
	if db.Driver != "sqlite" {
		return nil
	}

	// Essential pragmas for every connection
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",           // Enable WAL mode for concurrent readers
		"PRAGMA synchronous = NORMAL",           // Safe with WAL, faster than FULL
		"PRAGMA temp_store = memory",          // Store temp tables/indexes in memory
		"PRAGMA mmap_size = 268435456",        // 256MB memory map (adjust based on RAM)
		"PRAGMA wal_autocheckpoint = 1000",    // Checkpoint every 1000 pages (prevents WAL growth)
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute %s: %w", pragma, err)
		}
	}

	return nil
}

// Optimize runs PRAGMA optimize before closing database connection
// Recommended for long-running applications
func (db *DB) Optimize() error {
	if db.Driver != "sqlite" {
		return nil
	}
	_, err := db.Exec("PRAGMA optimize")
	return err
}

// CheckpointWAL manually checkpoints the WAL file
// Use this periodically to prevent WAL from growing too large
// 'passive' = don't block readers, 'full' = block but complete, 'truncate' = reset WAL
func (db *DB) CheckpointWAL(mode string) error {
	if db.Driver != "sqlite" {
		return nil
	}
	if mode == "" {
		mode = "passive"
	}
	_, err := db.Exec(fmt.Sprintf("PRAGMA wal_checkpoint(%s)", mode))
	return err
}

// Vacuum reclaims storage and defragments database
// Run this during low-traffic periods (expensive operation)
func (db *DB) Vacuum() error {
	if db.Driver != "sqlite" {
		return nil
	}
	_, err := db.Exec("VACUUM")
	return err
}

// EnableIncrementalVacuum sets up auto-vacuum for databases that shrink regularly
// Run once: db.EnableIncrementalVacuum() on fresh database
func (db *DB) EnableIncrementalVacuum() error {
	if db.Driver != "sqlite" {
		return nil
	}
	_, err := db.Exec("PRAGMA auto_vacuum = incremental")
	return err
}

// IncrementalVacuum moves freelist pages to end and truncates
// Run periodically if you delete data regularly
func (db *DB) IncrementalVacuum() error {
	if db.Driver != "sqlite" {
		return nil
	}
	_, err := db.Exec("PRAGMA incremental_vacuum")
	return err
}
