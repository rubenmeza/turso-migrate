package migration

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rubenmeza/turso-migrate/internal/storage"
)

// MigrationFile represents a migration file on disk
type MigrationFile struct {
	Version string
	Name    string
	Path    string
	UpSQL   string
	DownSQL string
}

// Engine handles Turso database migration operations
type Engine struct {
	storage       *storage.TursoStorage
	migrationsDir string
}

// NewEngine creates a new Turso migration engine
func NewEngine(storage *storage.TursoStorage, migrationsDir string) *Engine {
	return &Engine{
		storage:       storage,
		migrationsDir: migrationsDir,
	}
}

// Create creates a new migration file for Turso
func (e *Engine) Create(name string) error {
	// Get next version number
	version, err := e.getNextVersion()
	if err != nil {
		return fmt.Errorf("failed to get next version: %w", err)
	}

	// Sanitize name
	sanitizedName := sanitizeName(name)
	filename := fmt.Sprintf("%s_%s.sql", version, sanitizedName)
	filepath := filepath.Join(e.migrationsDir, filename)

	// Ensure migrations directory exists
	if err := os.MkdirAll(e.migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create migration file with template
	template := fmt.Sprintf(`-- Migration: %s
-- Created: %s

-- ==== UP ====


-- ==== DOWN ====

`, name, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(filepath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	fmt.Printf("Created migration: %s\n", filename)
	return nil
}

// Up applies all pending migrations
func (e *Engine) Up() error {
	// Get migration files
	files, err := e.loadMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No migrations found")
		return nil
	}

	// Get applied migrations
	applied, err := e.storage.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Build set of applied versions
	appliedSet := make(map[string]bool)
	for _, m := range applied {
		appliedSet[m.Version] = true
	}

	// Apply pending migrations
	var appliedCount int
	for _, file := range files {
		if appliedSet[file.Version] {
			continue // Skip already applied
		}

		fmt.Printf("Applying migration %s: %s\n", file.Version, file.Name)

		// Execute UP SQL
		if err := e.storage.ExecuteSQL(file.UpSQL); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Version, err)
		}

		// Record migration
		if err := e.storage.RecordMigration(file.Version, file.Name); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file.Version, err)
		}

		appliedCount++
	}

	if appliedCount == 0 {
		fmt.Println("No pending migrations")
	} else {
		fmt.Printf("Applied %d migration(s)\n", appliedCount)
	}

	return nil
}

// Down rolls back the last applied migration
func (e *Engine) Down() error {
	// Get applied migrations
	applied, err := e.storage.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Get the last applied migration
	lastMigration := applied[len(applied)-1]

	// Find the migration file
	files, err := e.loadMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	var migrationFile *MigrationFile
	for _, file := range files {
		if file.Version == lastMigration.Version {
			migrationFile = &file
			break
		}
	}

	if migrationFile == nil {
		return fmt.Errorf("migration file not found for version %s", lastMigration.Version)
	}

	if migrationFile.DownSQL == "" {
		return fmt.Errorf("no DOWN migration found for version %s", lastMigration.Version)
	}

	fmt.Printf("Rolling back migration %s: %s\n", migrationFile.Version, migrationFile.Name)

	// Execute DOWN SQL
	if err := e.storage.ExecuteSQL(migrationFile.DownSQL); err != nil {
		return fmt.Errorf("failed to execute rollback for %s: %w", migrationFile.Version, err)
	}

	// Remove migration record
	if err := e.storage.RemoveMigration(migrationFile.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", migrationFile.Version, err)
	}

	fmt.Println("Migration rolled back successfully")
	return nil
}

// Status shows the current migration status
func (e *Engine) Status() error {
	// Get migration files
	files, err := e.loadMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	// Get applied migrations
	applied, err := e.storage.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Build set of applied versions
	appliedSet := make(map[string]storage.Migration)
	for _, m := range applied {
		appliedSet[m.Version] = m
	}

	if len(files) == 0 {
		fmt.Println("No migrations found")
		return nil
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, file := range files {
		if migration, isApplied := appliedSet[file.Version]; isApplied {
			fmt.Printf("✓ %s_%s (applied: %s)\n",
				file.Version,
				file.Name,
				migration.AppliedAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("✗ %s_%s (pending)\n", file.Version, file.Name)
		}
	}

	return nil
}

// Version shows the current schema version
func (e *Engine) Version() error {
	version, err := e.storage.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if version == "" {
		fmt.Println("No migrations applied yet")
	} else {
		fmt.Printf("Current version: %s\n", version)
	}

	return nil
}

// loadMigrationFiles loads all migration files from the migrations directory
func (e *Engine) loadMigrationFiles() ([]MigrationFile, error) {
	var files []MigrationFile

	err := filepath.WalkDir(e.migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		file, err := e.parseMigrationFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		files = append(files, *file)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by version
	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	return files, nil
}

// parseMigrationFile parses a single migration file
func (e *Engine) parseMigrationFile(path string) (*MigrationFile, error) {
	// Parse filename for version and name
	filename := filepath.Base(path)
	re := regexp.MustCompile(`^(\d+)_(.+)\.sql$`)
	matches := re.FindStringSubmatch(filename)

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	version := matches[1]
	name := matches[2]

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse UP and DOWN sections
	upSQL, downSQL := parseSQL(string(content))

	return &MigrationFile{
		Version: version,
		Name:    name,
		Path:    path,
		UpSQL:   upSQL,
		DownSQL: downSQL,
	}, nil
}

// parseSQL parses UP and DOWN SQL from migration content
func parseSQL(content string) (upSQL, downSQL string) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentSection string
	var upLines, downLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "==== UP ====") {
			currentSection = "up"
			continue
		}

		if strings.Contains(line, "==== DOWN ====") {
			currentSection = "down"
			continue
		}

		switch currentSection {
		case "up":
			upLines = append(upLines, scanner.Text())
		case "down":
			downLines = append(downLines, scanner.Text())
		}
	}

	return strings.TrimSpace(strings.Join(upLines, "\n")),
		strings.TrimSpace(strings.Join(downLines, "\n"))
}

// getNextVersion returns the next migration version number
func (e *Engine) getNextVersion() (string, error) {
	files, err := e.loadMigrationFiles()
	if err != nil {
		// If directory doesn't exist, start from 001
		if os.IsNotExist(err) {
			return "001", nil
		}
		return "", err
	}

	if len(files) == 0 {
		return "001", nil
	}

	// Get the last version and increment
	lastFile := files[len(files)-1]
	lastVersion, err := strconv.Atoi(lastFile.Version)
	if err != nil {
		return "", fmt.Errorf("invalid version format: %s", lastFile.Version)
	}

	nextVersion := lastVersion + 1
	return fmt.Sprintf("%03d", nextVersion), nil
}

// sanitizeName sanitizes a migration name for use in filename
func sanitizeName(name string) string {
	// Replace spaces and special characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	sanitized := re.ReplaceAllString(name, "_")

	// Remove multiple underscores
	re = regexp.MustCompile(`_+`)
	sanitized = re.ReplaceAllString(sanitized, "_")

	// Trim underscores from start and end
	return strings.Trim(sanitized, "_")
}
