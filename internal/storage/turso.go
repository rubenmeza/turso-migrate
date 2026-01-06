package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// TursoStorage handles database operations for Turso migrations
type TursoStorage struct {
	db *sql.DB
}

// Migration represents a single migration record
type Migration struct {
	Version   string
	Name      string
	AppliedAt time.Time
}

// New creates a new TursoStorage instance
func New(databaseURL, authToken string) (*TursoStorage, error) {
	// Construct the connection string with auth token
	connStr := databaseURL
	if authToken != "" {
		connStr = fmt.Sprintf("%s?authToken=%s", databaseURL, authToken)
	}

	db, err := sql.Open("libsql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &TursoStorage{db: db}

	// Initialize schema migrations table
	if err := storage.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return storage, nil
}

// InitSchema creates the schema_migrations table if it doesn't exist
func (s *TursoStorage) InitSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := s.db.Exec(query)
	return err
}

// RecordMigration records a migration as applied
func (s *TursoStorage) RecordMigration(version, name string) error {
	query := `
		INSERT INTO schema_migrations (version, name, applied_at)
		VALUES (?, ?, ?)
	`
	_, err := s.db.Exec(query, version, name, time.Now())
	return err
}

// RemoveMigration removes a migration record
func (s *TursoStorage) RemoveMigration(version string) error {
	query := `DELETE FROM schema_migrations WHERE version = ?`
	_, err := s.db.Exec(query, version)
	return err
}

// GetAppliedMigrations returns all applied migrations ordered by version
func (s *TursoStorage) GetAppliedMigrations() ([]Migration, error) {
	query := `
		SELECT version, name, applied_at 
		FROM schema_migrations 
		ORDER BY version ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Version, &m.Name, &m.AppliedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}

	return migrations, rows.Err()
}

// IsMigrationApplied checks if a migration has been applied
func (s *TursoStorage) IsMigrationApplied(version string) (bool, error) {
	query := `SELECT COUNT(*) FROM schema_migrations WHERE version = ?`
	var count int
	err := s.db.QueryRow(query, version).Scan(&count)
	return count > 0, err
}

// ExecuteSQL executes a SQL statement in a transaction
func (s *TursoStorage) ExecuteSQL(sql string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(sql)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetCurrentVersion returns the latest applied migration version
func (s *TursoStorage) GetCurrentVersion() (string, error) {
	query := `
		SELECT version 
		FROM schema_migrations 
		ORDER BY version DESC 
		LIMIT 1
	`

	var version string
	err := s.db.QueryRow(query).Scan(&version)
	if err == sql.ErrNoRows {
		return "", nil // No migrations applied
	}
	return version, err
}

// Close closes the database connection
func (s *TursoStorage) Close() error {
	return s.db.Close()
}
