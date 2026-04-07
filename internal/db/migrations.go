package db

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Migration represents a single migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
}

// EnsureMigrationsTable creates the schema_migrations table if it doesn't exist
func (db *DB) EnsureMigrationsTable() error {
	createTableSQL := `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version TEXT PRIMARY KEY,
	applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	_, err := db.Exec(createTableSQL)
	return err
}

// GetAppliedMigrations returns list of already applied migration versions
func (db *DB) GetAppliedMigrations() ([]string, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, rows.Err()
}

// MarkMigrationApplied records that a migration has been applied
func (db *DB) MarkMigrationApplied(version string) error {
	if db.Driver == "postgres" {
		_, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT (version) DO NOTHING", version)
		return err
	}
	// SQLite
	_, err := db.Exec("INSERT OR IGNORE INTO schema_migrations (version) VALUES (?)", version)
	return err
}

// LoadMigrationsFromDir reads all migration files from a directory
// For multi-driver support, it looks for driver-specific files first:
//   - 000001_name.sqlite.up.sql (SQLite-specific)
//   - 000001_name.postgres.up.sql (PostgreSQL-specific)
//   - 000001_name.up.sql (generic/default)
func LoadMigrationsFromDir(dir string, driver string) ([]Migration, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Pattern to parse: 000001_description.[driver.]up.sql
	// Driver-specific: 000001_description.sqlite.up.sql
	// Generic: 000001_description.up.sql
	pattern := regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)
	driverPattern := regexp.MustCompile(`^(\d+)_(.+)\.(sqlite|postgres)\.(up|down)\.sql$`)

	migrations := make(map[string]*Migration)

	// First pass: collect driver-specific migrations
	driverSpecific := make(map[string]map[string]string) // version -> direction -> filename

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()

		// Check for driver-specific first
		matches := driverPattern.FindStringSubmatch(name)
		if matches != nil {
			version := matches[1]
			fileDriver := matches[3]
			direction := matches[4]

			if fileDriver == driver {
				if driverSpecific[version] == nil {
					driverSpecific[version] = make(map[string]string)
				}
				driverSpecific[version][direction] = name
			}
			continue
		}

		// Check for generic
		matches = pattern.FindStringSubmatch(name)
		if matches != nil {
			version := matches[1]
			direction := matches[3]

			// Only use generic if no driver-specific version exists
			if _, hasDriver := driverSpecific[version]; !hasDriver {
				if driverSpecific[version] == nil {
					driverSpecific[version] = make(map[string]string)
				}
				// Only set if not already set (driver-specific takes priority)
				if _, exists := driverSpecific[version][direction]; !exists {
					driverSpecific[version][direction] = name
				}
			}
		}
	}

	// Second pass: load the selected files
	for version, directions := range driverSpecific {
		for direction, filename := range directions {
			content, err := os.ReadFile(filepath.Join(dir, filename))
			if err != nil {
				return nil, fmt.Errorf("failed to read migration %s: %w", filename, err)
			}

			sql := cleanSQL(string(content))

			if migrations[version] == nil {
				migrations[version] = &Migration{
					Version: version,
					Name:    filename, // simplified
				}
			}

			if direction == "up" {
				migrations[version].UpSQL = sql
			} else {
				migrations[version].DownSQL = sql
			}
		}
	}

	// Convert to slice and sort by version
	var result []Migration
	for _, m := range migrations {
		result = append(result, *m)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

// RunMigrations executes all pending migrations
func (db *DB) RunMigrations(migrationsDir string) error {
	// Ensure tracking table exists
	if err := db.EnsureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get already applied migrations
	applied, err := db.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedSet := make(map[string]bool)
	for _, v := range applied {
		appliedSet[v] = true
	}

	// Load all migrations (driver-specific)
	migrations, err := LoadMigrationsFromDir(migrationsDir, db.Driver)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		return nil
	}

	// Run pending migrations in transaction
	for _, m := range migrations {
		if appliedSet[m.Version] {
			// Already applied, skip
			continue
		}

		if m.UpSQL == "" {
			return fmt.Errorf("migration %s has no up SQL", m.Version)
		}

		// Execute migration
		if _, err := db.Exec(m.UpSQL); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", m.Version, err)
		}

		// Mark as applied
		if err := db.MarkMigrationApplied(m.Version); err != nil {
			return fmt.Errorf("failed to mark migration %s as applied: %w", m.Version, err)
		}
	}

	return nil
}

// cleanSQL removes comments and normalizes SQL
func cleanSQL(sql string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(sql))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments
		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
			continue
		}
		if line != "" {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return strings.TrimSpace(result.String())
}

// IsMigrationApplied checks if a specific migration version was applied
func (db *DB) IsMigrationApplied(version string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = "+db.Placeholder(1), version).Scan(&count)
	return count > 0, err
}
