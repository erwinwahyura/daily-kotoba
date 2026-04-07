package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// SeedRecord tracks what seed data has been applied
type SeedRecord struct {
	Name       string    `json:"name"`
	AppliedAt  time.Time `json:"applied_at"`
	Checksum   string    `json:"checksum"` // Optional: detect changes
	RecordCount int      `json:"record_count"`
}

// EnsureSeedsTable creates the seed tracking table
func (db *DB) EnsureSeedsTable() error {
	sql := `
CREATE TABLE IF NOT EXISTS schema_seeds (
	name TEXT PRIMARY KEY,
	applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	checksum TEXT,
	record_count INTEGER DEFAULT 0
);`
	_, err := db.Exec(sql)
	return err
}

// IsSeedApplied checks if a seed has been applied
func (db *DB) IsSeedApplied(name string) (bool, error) {
	var count int
	placeholder := db.Placeholder(1)
	err := db.QueryRow(
		"SELECT COUNT(*) FROM schema_seeds WHERE name = "+placeholder,
		name,
	).Scan(&count)
	return count > 0, err
}

// MarkSeedApplied records that a seed was applied
func (db *DB) MarkSeedApplied(name string, checksum string, count int) error {
	if db.Driver == "postgres" {
		_, err := db.Exec(
			"INSERT INTO schema_seeds (name, checksum, record_count) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET applied_at = CURRENT_TIMESTAMP, checksum = $2, record_count = $3",
			name, checksum, count,
		)
		return err
	}
	// SQLite
	_, err := db.Exec(
		"INSERT OR REPLACE INTO schema_seeds (name, applied_at, checksum, record_count) VALUES (?, CURRENT_TIMESTAMP, ?, ?)",
		name, checksum, count,
	)
	return err
}

// SeedData represents seed data to be loaded
type SeedData struct {
	Name    string
	Type    string // "vocabulary", "grammar", "placement"
	Records []map[string]interface{}
}

// LoadSeedJSON loads seed data from a JSON file
// Uses json.RawMessage to preserve JSON fields as strings for proper database insertion
func LoadSeedJSON(path string) (*SeedData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read seed file %s: %w", path, err)
	}

	// First, parse into raw messages to preserve JSON structure
	var rawRecords []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawRecords); err != nil {
		return nil, fmt.Errorf("failed to parse seed JSON %s: %w", path, err)
	}

	// Convert to map[string]interface{} with JSON fields as strings
	records := make([]map[string]interface{}, len(rawRecords))
	for i, rawRecord := range rawRecords {
		record := make(map[string]interface{})
		for key, rawValue := range rawRecord {
			// Check if this is a JSON array/object (starts with [ or {)
			trimmed := strings.TrimSpace(string(rawValue))
			if len(trimmed) > 0 && (trimmed[0] == '[' || trimmed[0] == '{') {
				// Keep as string for JSON fields
				record[key] = string(rawValue)
			} else {
				// Parse as regular value
				var value interface{}
				if err := json.Unmarshal(rawValue, &value); err != nil {
					record[key] = string(rawValue) // Fallback to string
				} else {
					record[key] = value
				}
			}
		}
		records[i] = record
	}

	// Extract seed name from filename
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Determine type from filename or content
	seedType := "unknown"
	if strings.Contains(name, "vocab") {
		seedType = "vocabulary"
	} else if strings.Contains(name, "grammar") {
		seedType = "grammar"
	} else if strings.Contains(name, "placement") {
		seedType = "placement"
	}

	return &SeedData{
		Name:    name,
		Type:    seedType,
		Records: records,
	}, nil
}

// SeedVocabulary inserts vocabulary data from seed file
func (db *DB) SeedVocabulary(seedFile string) (int, error) {
	seedData, err := LoadSeedJSON(seedFile)
	if err != nil {
		return 0, err
	}

	// Check if already seeded
	applied, err := db.IsSeedApplied(seedData.Name)
	if err != nil {
		return 0, err
	}
	if applied {
		return 0, nil // Already seeded, skip
	}

	count := 0
	for _, record := range seedData.Records {
		// Build insert query dynamically
		columns := make([]string, 0)
		placeholders := make([]string, 0)
		values := make([]interface{}, 0)

		for col, val := range record {
			columns = append(columns, col)
			placeholders = append(placeholders, db.Placeholder(len(values)+1))
			
			// Handle JSON fields
			if col == "example_sentences" || col == "related_words" {
				jsonVal, err := db.JSONValue(val)
				if err != nil {
					return count, fmt.Errorf("failed to marshal JSON for %s: %w", col, err)
				}
				values = append(values, jsonVal)
			} else {
				values = append(values, val)
			}
		}

		query := fmt.Sprintf(
			"INSERT INTO vocabulary (%s) VALUES (%s)",
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
		)

		if _, err := db.Exec(query, values...); err != nil {
			// Log error but continue (duplicate or constraint error)
			if !isDuplicateError(err, db.Driver) {
				return count, fmt.Errorf("failed to insert vocab record: %w", err)
			}
		} else {
			count++
		}
	}

	// Mark as seeded
	checksum := fmt.Sprintf("records:%d", len(seedData.Records))
	if err := db.MarkSeedApplied(seedData.Name, checksum, count); err != nil {
		return count, fmt.Errorf("failed to mark seed applied: %w", err)
	}

	return count, nil
}

// SeedGrammar inserts grammar pattern data from seed file
func (db *DB) SeedGrammar(seedFile string) (int, error) {
	seedData, err := LoadSeedJSON(seedFile)
	if err != nil {
		return 0, err
	}

	applied, err := db.IsSeedApplied(seedData.Name)
	if err != nil {
		return 0, err
	}
	if applied {
		return 0, nil
	}

	count := 0
	for _, record := range seedData.Records {
		columns := make([]string, 0)
		placeholders := make([]string, 0)
		values := make([]interface{}, 0)

		for col, val := range record {
			columns = append(columns, col)
			placeholders = append(placeholders, db.Placeholder(len(values)+1))
			
			// Handle JSON fields
			if col == "usage_examples" || col == "related_patterns" || col == "common_mistakes" {
				jsonVal, err := db.JSONValue(val)
				if err != nil {
					return count, fmt.Errorf("failed to marshal JSON for %s: %w", col, err)
				}
				values = append(values, jsonVal)
			} else {
				values = append(values, val)
			}
		}

		query := fmt.Sprintf(
			"INSERT INTO grammar_patterns (%s) VALUES (%s)",
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
		)

		if _, err := db.Exec(query, values...); err != nil {
			if !isDuplicateError(err, db.Driver) {
				return count, fmt.Errorf("failed to insert grammar record: %w", err)
			}
		} else {
			count++
		}
	}

	checksum := fmt.Sprintf("records:%d", len(seedData.Records))
	if err := db.MarkSeedApplied(seedData.Name, checksum, count); err != nil {
		return count, err
	}

	return count, nil
}

// SeedPlacement inserts placement test questions from seed file
func (db *DB) SeedPlacement(seedFile string) (int, error) {
	seedData, err := LoadSeedJSON(seedFile)
	if err != nil {
		return 0, err
	}

	applied, err := db.IsSeedApplied(seedData.Name)
	if err != nil {
		return 0, err
	}
	if applied {
		return 0, nil
	}

	count := 0
	for _, record := range seedData.Records {
		columns := make([]string, 0)
		placeholders := make([]string, 0)
		values := make([]interface{}, 0)

		for col, val := range record {
			columns = append(columns, col)
			placeholders = append(placeholders, db.Placeholder(len(values)+1))
			
			// Handle JSON fields
			if col == "options" || col == "hints" {
				jsonVal, err := db.JSONValue(val)
				if err != nil {
					return count, fmt.Errorf("failed to marshal JSON for %s: %w", col, err)
				}
				values = append(values, jsonVal)
			} else {
				values = append(values, val)
			}
		}

		query := fmt.Sprintf(
			"INSERT INTO placement_questions (%s) VALUES (%s)",
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
		)

		if _, err := db.Exec(query, values...); err != nil {
			if !isDuplicateError(err, db.Driver) {
				return count, fmt.Errorf("failed to insert placement record: %w", err)
			}
		} else {
			count++
		}
	}

	checksum := fmt.Sprintf("records:%d", len(seedData.Records))
	if err := db.MarkSeedApplied(seedData.Name, checksum, count); err != nil {
		return count, err
	}

	return count, nil
}

// RunAutoSeeding scans a directory and applies all pending seed files
func (db *DB) RunAutoSeeding(seedsDir string) error {
	// Ensure tracking table exists
	if err := db.EnsureSeedsTable(); err != nil {
		return fmt.Errorf("failed to create seeds table: %w", err)
	}

	// Get list of seed files
	files, err := os.ReadDir(seedsDir)
	if err != nil {
		return fmt.Errorf("failed to read seeds directory: %w", err)
	}

	// Sort files by name for consistent ordering
	var seedFiles []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasSuffix(name, ".json") {
			seedFiles = append(seedFiles, name)
		}
	}
	sort.Strings(seedFiles)

	log.Printf("Found %d seed files: %v", len(seedFiles), seedFiles)

	// Apply each seed file
	for _, name := range seedFiles {
		path := filepath.Join(seedsDir, name)
		
		// Skip already applied seeds (double-check)
		seedName := strings.TrimSuffix(name, ".json")
		applied, _ := db.IsSeedApplied(seedName)
		if applied {
			log.Printf("Seed %s already applied, skipping", seedName)
			continue
		}

		log.Printf("Applying seed: %s", seedName)

		var count int
		var err error

		// Determine seed type and apply
		if strings.Contains(name, "vocab") {
			count, err = db.SeedVocabulary(path)
		} else if strings.Contains(name, "grammar") {
			count, err = db.SeedGrammar(path)
		} else if strings.Contains(name, "placement") {
			count, err = db.SeedPlacement(path)
		} else {
			// Unknown type, try generic approach
			log.Printf("Unknown seed type for %s, skipping", name)
			continue
		}

		if err != nil {
			log.Printf("Failed to apply seed %s: %v", name, err)
			return fmt.Errorf("failed to apply seed %s: %w", name, err)
		}

		log.Printf("Successfully applied seed %s: %d records inserted", seedName, count)
	}

	return nil
}

// isDuplicateError checks if error is a duplicate/unique constraint violation
func isDuplicateError(err error, driver string) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	if driver == "postgres" {
		return strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "unique_violation")
	}
	// SQLite
	return strings.Contains(errStr, "UNIQUE constraint failed") || strings.Contains(errStr, "duplicate")
}
