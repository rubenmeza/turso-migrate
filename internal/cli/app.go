package cli

import (
	"fmt"
	"os"

	"github.com/rubenmeza/turso-migrate/internal/migration"
	"github.com/rubenmeza/turso-migrate/internal/storage"
	"github.com/rubenmeza/turso-migrate/pkg/config"
	"github.com/urfave/cli/v2"
)

const version = "1.0.0"

// NewApp creates a new CLI application
func NewApp() *cli.App {
	return &cli.App{
		Name:    "turso-migrate",
		Usage:   "A simple database migration tool for Turso (libSQL)",
		Version: version,
		Authors: []*cli.Author{
			{
				Name: "turso-migrate contributors",
			},
		},
		Description: `turso-migrate is a simple and reliable migration tool for Turso databases.
Built specifically for libSQL with first-class Docker support and 
straightforward integration for modern applications.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "database-url",
				Aliases: []string{"d"},
				Usage:   "Turso database URL (overrides TURSO_DATABASE_URL)",
				EnvVars: []string{"TURSO_DATABASE_URL"},
			},
			&cli.StringFlag{
				Name:    "auth-token",
				Aliases: []string{"t"},
				Usage:   "Turso auth token (overrides TURSO_AUTH_TOKEN)",
				EnvVars: []string{"TURSO_AUTH_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "migrations-dir",
				Aliases: []string{"m"},
				Usage:   "Directory containing migration files",
				Value:   "./migrations",
				EnvVars: []string{"MIGRATIONS_DIR"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "create",
				Aliases:   []string{"c"},
				Usage:     "Create a new migration file for your Turso database",
				ArgsUsage: "<name>",
				Action:    createCommand,
				Description: `Create a new migration file with the given name.
The file will be created with auto-incremented version number and
pre-filled UP and DOWN sections optimized for Turso/libSQL.

Example:
  turso-migrate create add_users_table`,
			},
			{
				Name:    "up",
				Aliases: []string{"u"},
				Usage:   "Apply all pending migrations to your Turso database",
				Action:  upCommand,
				Description: `Apply all pending migrations in order to your Turso database.
Only migrations that haven't been applied yet will be executed.
Each migration runs in its own transaction for data safety.`,
			},
			{
				Name:    "down",
				Aliases: []string{"d"},
				Usage:   "Rollback the last applied migration from your Turso database",
				Action:  downCommand,
				Description: `Rollback the most recently applied migration from your Turso database.
This will execute the DOWN section of the migration file.
Use with caution in production environments.`,
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Show migration status for your Turso database",
				Action:  statusCommand,
				Description: `Show the status of all migrations for your Turso database.
Displays which migrations have been applied and which are pending.`,
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show current schema version of your Turso database",
				Action:  versionCommand,
				Description: `Show the current schema version of your Turso database.
This is the version of the last applied migration.`,
			},
		},
		Before: func(c *cli.Context) error {
			// Validate that we have required Turso configuration
			cfg := buildConfig(c)
			return cfg.Validate()
		},
	}
}

func createCommand(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("migration name is required")
	}

	name := c.Args().First()
	cfg := buildConfig(c)

	// Ensure migrations directory exists
	if err := cfg.EnsureMigrationsDir(); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create storage (we don't need it for creating files, but validate connection)
	store, err := storage.New(cfg.DatabaseURL, cfg.AuthToken)
	if err != nil {
		return fmt.Errorf("failed to connect to Turso database: %w", err)
	}
	defer store.Close()

	engine := migration.NewEngine(store, cfg.MigrationsDir)
	return engine.Create(name)
}

func upCommand(c *cli.Context) error {
	cfg := buildConfig(c)

	store, err := storage.New(cfg.DatabaseURL, cfg.AuthToken)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	engine := migration.NewEngine(store, cfg.MigrationsDir)
	return engine.Up()
}

func downCommand(c *cli.Context) error {
	cfg := buildConfig(c)

	store, err := storage.New(cfg.DatabaseURL, cfg.AuthToken)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	engine := migration.NewEngine(store, cfg.MigrationsDir)
	return engine.Down()
}

func statusCommand(c *cli.Context) error {
	cfg := buildConfig(c)

	store, err := storage.New(cfg.DatabaseURL, cfg.AuthToken)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	engine := migration.NewEngine(store, cfg.MigrationsDir)
	return engine.Status()
}

func versionCommand(c *cli.Context) error {
	cfg := buildConfig(c)

	store, err := storage.New(cfg.DatabaseURL, cfg.AuthToken)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	engine := migration.NewEngine(store, cfg.MigrationsDir)
	return engine.Version()
}

func buildConfig(c *cli.Context) *config.Config {
	cfg := &config.Config{
		DatabaseURL:   c.String("database-url"),
		AuthToken:     c.String("auth-token"),
		MigrationsDir: c.String("migrations-dir"),
	}

	// Load from environment if not provided via flags
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = os.Getenv("TURSO_DATABASE_URL")
	}
	if cfg.AuthToken == "" {
		cfg.AuthToken = os.Getenv("TURSO_AUTH_TOKEN")
	}
	if cfg.MigrationsDir == "" {
		cfg.MigrationsDir = "./migrations"
	}

	return cfg
}
