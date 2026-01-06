package config

import (
	"errors"
	"os"
	"path/filepath"
)

// Config holds the configuration for turso-migrate and Turso database connection
type Config struct {
	DatabaseURL   string
	AuthToken     string
	MigrationsDir string
}

// LoadFromEnv loads Turso configuration from environment variables
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		DatabaseURL:   os.Getenv("TURSO_DATABASE_URL"),
		AuthToken:     os.Getenv("TURSO_AUTH_TOKEN"),
		MigrationsDir: os.Getenv("MIGRATIONS_DIR"),
	}

	// Set default migrations directory
	if cfg.MigrationsDir == "" {
		cfg.MigrationsDir = "./migrations"
	}

	// Clean up the path
	cfg.MigrationsDir = filepath.Clean(cfg.MigrationsDir)

	return cfg, cfg.Validate()
}

// Validate checks if the Turso configuration is valid
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("TURSO_DATABASE_URL is required")
	}
	if c.AuthToken == "" {
		return errors.New("TURSO_AUTH_TOKEN is required")
	}
	return nil
}

// EnsureMigrationsDir creates the migrations directory if it doesn't exist
func (c *Config) EnsureMigrationsDir() error {
	return os.MkdirAll(c.MigrationsDir, 0755)
}
