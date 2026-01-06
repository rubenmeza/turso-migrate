# Turso Migration Tools - Simple migration workflow for Turso databases
# Integrate this into your Go project for easy database migrations

# Configuration
MIGRATION_NAME ?= 
TURSO_DATABASE_URL ?= 
TURSO_AUTH_TOKEN ?= 
MIGRATIONS_DIR ?= ./migrations

# Docker image (replace with your registry if publishing)
DOCKER_IMAGE = turso-migrate:latest

# Check if Docker should be used (useful for CI environments)
USE_DOCKER ?= false

# Binary selection (local vs Docker)
ifeq ($(USE_DOCKER),true)
    TURSO_MIGRATE = docker run --rm \
        -e TURSO_DATABASE_URL=$(TURSO_DATABASE_URL) \
        -e TURSO_AUTH_TOKEN=$(TURSO_AUTH_TOKEN) \
        -e MIGRATIONS_DIR=/migrations \
        -v $(PWD)/$(MIGRATIONS_DIR):/migrations \
        $(DOCKER_IMAGE)
else
    TURSO_MIGRATE = ./bin/turso-migrate
endif

.PHONY: help migration migrate migrate-up migrate-down migrate-status migrate-version \
        migrate-docker build-migrate install-migrate

help: ## Show this help message
	@echo "Turso Database Migration Commands"
	@echo "================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Migration creation
migration: ## Create a new migration (use: make migration MIGRATION_NAME=name)
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Error: MIGRATION_NAME is required"; \
		echo "Usage: make migration MIGRATION_NAME=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration: $(MIGRATION_NAME)"
	$(TURSO_MIGRATE) create $(MIGRATION_NAME)

# Apply all pending migrations
migrate: migrate-up ## Alias for migrate-up

migrate-up: ## Apply all pending migrations to Turso database
	@echo "Applying migrations to Turso database..."
	$(TURSO_MIGRATE) up

# Rollback last migration
migrate-down: ## Rollback the last applied migration
	@echo "Rolling back last migration..."
	$(TURSO_MIGRATE) down

# Show migration status
migrate-status: ## Show current migration status
	$(TURSO_MIGRATE) status

# Show current version
migrate-version: ## Show current schema version
	$(TURSO_MIGRATE) version

# Docker-specific commands
migrate-docker: ## Run migrations using Docker
	@$(MAKE) migrate USE_DOCKER=true

build-migrate: ## Build Docker image for turso-migrate
	docker build -t $(DOCKER_IMAGE) .

# Development helpers
install-migrate: ## Install turso-migrate binary locally
	@echo "Building turso-migrate for Turso databases..."
	@mkdir -p ./bin
	go build -o ./bin/turso-migrate ./cmd/turso-migrate
	@echo "turso-migrate installed to ./bin/turso-migrate"

# Environment validation
check-env: ## Check if required Turso environment variables are set
	@echo "Checking Turso configuration..."
	@if [ -z "$(TURSO_DATABASE_URL)" ]; then \
		echo "Error: TURSO_DATABASE_URL is not set"; \
		echo "Get your database URL from: https://app.turso.tech/"; \
		exit 1; \
	fi
	@if [ -z "$(TURSO_AUTH_TOKEN)" ]; then \
		echo "Error: TURSO_AUTH_TOKEN is not set"; \
		echo "Create an auth token at: https://app.turso.tech/"; \
		exit 1; \
	fi
	@echo "Turso configuration is valid ✓"

# Combined commands for common workflows
setup-migrations: install-migrate check-env ## Setup migrations (install binary + check env)
	@echo "Turso migration setup complete ✓"

# CI/CD friendly commands
ci-migrate: check-env ## Run migrations in CI environment
	@echo "Running Turso migrations in CI mode..."
	@$(MAKE) migrate-status
	@$(MAKE) migrate-up

# Development workflow
dev-setup: ## Setup development environment for Turso migrations
	@echo "Setting up Turso development environment..."
	@$(MAKE) install-migrate
	@mkdir -p $(MIGRATIONS_DIR)
	@echo "Development environment ready ✓"
	@echo ""
	@echo "Next steps:"
	@echo "1. Set your Turso credentials in .env or environment"
	@echo "2. Create your first migration: make migration MIGRATION_NAME=initial_schema"
	@echo "3. Apply migrations: make migrate"

# Turso-specific helpers
turso-info: ## Show Turso database information
	@echo "Turso Database Information:"
	@echo "=========================="
	@echo "Database URL: $(TURSO_DATABASE_URL)"
	@echo "Migrations Dir: $(MIGRATIONS_DIR)"
	@if [ -f "$(MIGRATIONS_DIR)" ]; then \
		echo "Migration files: $$(ls -1 $(MIGRATIONS_DIR)/*.sql 2>/dev/null | wc -l)"; \
	else \
		echo "Migration files: 0 (directory not found)"; \
	fi