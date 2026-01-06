# turso-migrate

A **simple, Go-native database migration tool for Turso (libSQL)**. Easy to use, reliable, and built specifically for modern applications using Turso databases.

## Overview

`turso-migrate` is a **straightforward migration solution for Turso** designed to make database schema management simple and reliable. Built from the ground up for **libSQL databases** with **first-class Docker support**, **seamless CI/CD integration**, and **production-ready reliability**.

### Why turso-migrate?

- **🎯 Turso-focused**: Built specifically for libSQL/Turso databases
- **🐳 Docker-ready**: Easy containerized workflows  
- **⚡ Simple & fast**: Single compiled binary, zero dependencies
- **📦 CI/CD friendly**: Seamless integration with deployment pipelines
- **🔒 Reliable**: Transaction safety and comprehensive error handling

### Inspiration

`turso-migrate` draws inspiration from the excellent [shmig](https://github.com/naquad/shmig) migration tool, adapting its proven workflow concepts for the Turso ecosystem. While shmig provided great foundations for database migrations, turso-migrate is built specifically for Turso's libSQL architecture and modern Go development workflows.

---

## Quick Start

### 1. Install

```bash
# Download binary (recommended)
go install github.com/rubenmeza/turso-migrate/cmd/turso-migrate@latest

# Or build from source
git clone https://github.com/rubenmeza/turso-migrate
cd turso-migrate && make install-migrate
```

### 2. Configure

Set your Turso credentials:

```bash
export TURSO_DATABASE_URL="libsql://your-db.turso.io"
export TURSO_AUTH_TOKEN="your-auth-token"
export MIGRATIONS_DIR="./migrations"  # Optional, defaults to ./migrations
```

### 3. Create your first migration

```bash
turso-migrate create create_users_table
```

This creates `001_create_users_table.sql`:

```sql
-- Migration: create_users_table
-- Created: 2024-01-05 15:04:05

-- ==== UP ====


-- ==== DOWN ====

```

### 4. Add your SQL

Edit the migration file:

```sql
-- ==== UP ====
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==== DOWN ====
DROP TABLE users;
```

### 5. Apply migrations

```bash
turso-migrate up
```

---

## Migration File Format

Each migration contains **both UP and DOWN** sections in a single file:

```sql
-- Migration: Description of what this does
-- Created: 2024-01-05 15:04:05

-- ==== UP ====
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT
);

-- ==== DOWN ====
DROP TABLE posts;
```

### File Naming Convention

```
001_create_users.sql
002_add_posts_table.sql
003_add_user_indexes.sql
```

- **Auto-incremented numbers** define execution order
- **Descriptive names** help with organization
- **Single `.sql` extension** keeps it simple

---

## CLI Reference

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `create <name>` | Create new migration file | `turso-migrate create add_users` |
| `up` | Apply all pending migrations | `turso-migrate up` |
| `down` | Rollback last migration | `turso-migrate down` |
| `status` | Show migration status | `turso-migrate status` |
| `version` | Show current schema version | `turso-migrate version` |

### Global Flags

| Flag | Short | Environment | Default | Description |
|------|-------|-------------|---------|-------------|
| `--database-url` | `-d` | `TURSO_DATABASE_URL` | - | Turso database URL |
| `--auth-token` | `-t` | `TURSO_AUTH_TOKEN` | - | Turso auth token |
| `--migrations-dir` | `-m` | `MIGRATIONS_DIR` | `./migrations` | Migration files directory |

### Examples

```bash
# Create migration with custom directory
turso-migrate --migrations-dir ./db/migrations create add_indexes

# Apply migrations with inline credentials
turso-migrate --database-url libsql://mydb.turso.io --auth-token token123 up

# Check status
turso-migrate status
```

---

## Makefile Integration

### Complete Makefile Example

Integrate turso-migrate into your Go project workflow:

```makefile
# Turso Migration Configuration
MIGRATION_NAME ?= 
TURSO_DATABASE_URL ?= 
TURSO_AUTH_TOKEN ?= 
MIGRATIONS_DIR ?= ./migrations

# Use Docker in CI environments
USE_DOCKER ?= false
DOCKER_IMAGE = turso-migrate:latest

# Binary selection
ifeq ($(USE_DOCKER),true)
    TURSO_MIGRATE = docker run --rm \\
        -e TURSO_DATABASE_URL=$(TURSO_DATABASE_URL) \\
        -e TURSO_AUTH_TOKEN=$(TURSO_AUTH_TOKEN) \\
        -e MIGRATIONS_DIR=/migrations \\
        -v $(PWD)/$(MIGRATIONS_DIR):/migrations \\
        $(DOCKER_IMAGE)
else
    TURSO_MIGRATE = turso-migrate
endif

.PHONY: migration migrate migrate-up migrate-down migrate-status

# Create new migration
migration:
	@if [ -z "$(MIGRATION_NAME)" ]; then \\
		echo "Usage: make migration MIGRATION_NAME=create_users"; \\
		exit 1; \\
	fi
	$(TURSO_MIGRATE) create $(MIGRATION_NAME)

# Apply migrations
migrate: migrate-up
migrate-up:
	$(TURSO_MIGRATE) up

# Rollback last migration  
migrate-down:
	$(TURSO_MIGRATE) down

# Show status
migrate-status:
	$(TURSO_MIGRATE) status

# CI-friendly migration
ci-migrate:
	@echo "Running migrations in CI..."
	@$(MAKE) migrate-status
	@$(MAKE) migrate-up
```

### Usage Examples

```bash
# Create migration
make migration MIGRATION_NAME=create_users_table

# Apply migrations (local)
make migrate

# Apply migrations (Docker/CI)
make migrate USE_DOCKER=true

# Check status
make migrate-status

# Rollback
make migrate-down
```

---

## Docker Support

### Build the image

```bash
docker build -t turso-migrate .
```

### Run migrations

```bash
docker run --rm \\
  -e TURSO_DATABASE_URL=libsql://your-db.turso.io \\
  -e TURSO_AUTH_TOKEN=your-token \\
  -e MIGRATIONS_DIR=/migrations \\
  -v $(pwd)/migrations:/migrations \\
  turso-migrate up
```

### Docker Compose Example

```yaml
version: '3.8'

services:
  migrate:
    build: .
    environment:
      - TURSO_DATABASE_URL=${TURSO_DATABASE_URL}
      - TURSO_AUTH_TOKEN=${TURSO_AUTH_TOKEN}
      - MIGRATIONS_DIR=/migrations
    volumes:
      - ./migrations:/migrations
    command: up
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Database Migrations

on:
  push:
    branches: [main]

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Database Migrations
        run: |
          docker run --rm \\
            -e TURSO_DATABASE_URL=${{ secrets.TURSO_DATABASE_URL }} \\
            -e TURSO_AUTH_TOKEN=${{ secrets.TURSO_AUTH_TOKEN }} \\
            -e MIGRATIONS_DIR=/migrations \\
            -v $(pwd)/migrations:/migrations \\
            turso-migrate up
```

### GitLab CI

```yaml
migrate:
  stage: deploy
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t turso-migrate .
    - docker run --rm 
        -e TURSO_DATABASE_URL=$TURSO_DATABASE_URL
        -e TURSO_AUTH_TOKEN=$TURSO_AUTH_TOKEN
        -e MIGRATIONS_DIR=/migrations
        -v $(pwd)/migrations:/migrations
        turso-migrate up
  only:
    - main
```

---

## Migration Tracking

turso-migrate automatically creates and manages a `schema_migrations` table:

```sql
CREATE TABLE schema_migrations (
    version TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Query migration status

```sql
-- See all applied migrations
SELECT version, name, applied_at 
FROM schema_migrations 
ORDER BY version;

-- Get current version
SELECT version 
FROM schema_migrations 
ORDER BY version DESC 
LIMIT 1;
```

---

## Migration Tracking

turso-migrate automatically creates and manages a `schema_migrations` table:

```sql
CREATE TABLE schema_migrations (
    version TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Query migration status

```sql
-- See all applied migrations
SELECT version, name, applied_at 
FROM schema_migrations 
ORDER BY version;

-- Get current version
SELECT version 
FROM schema_migrations 
ORDER BY version DESC 
LIMIT 1;
```

---

## Architecture & Performance

### Design Principles

| Principle | Implementation |
|-----------|----------------|
| **Simplicity** | Easy to understand commands and clear migration format |
| **Turso-focused** | Built specifically for libSQL, optimized for Turso's features |
| **Reliability** | Each migration runs in its own transaction for safety |
| **Portability** | Single compiled binary, works anywhere |
| **Modern workflows** | Docker support and CI/CD integration out of the box |

### Migration Features

| Feature | Description |
|---------|-------------|
| **Auto-versioning** | Automatic sequential numbering (001, 002, 003...) |
| **UP/DOWN Support** | Both directions in single file |
| **Transaction Isolation** | Each migration in separate transaction |
| **Error Handling** | Comprehensive error reporting and rollback |
| **Idempotency** | Safe to run multiple times |

---

## Examples

### Example Migration Files

See the [examples/migrations](examples/migrations/) directory for complete examples:

- [001_create_users.sql](examples/migrations/001_create_users.sql) - Basic table creation
- [002_create_posts.sql](examples/migrations/002_create_posts.sql) - Foreign keys and indexes  
- [003_add_updated_at_triggers.sql](examples/migrations/003_add_updated_at_triggers.sql) - Triggers and advanced SQL

### Example Project Structure

```
your-go-project/
├── cmd/
├── internal/
├── migrations/           # Your migration files
│   ├── 001_create_users.sql
│   ├── 002_create_posts.sql
│   └── 003_add_indexes.sql
├── Makefile             # Migration commands
├── docker-compose.yml   # Include migration service
└── .github/
    └── workflows/
        └── migrate.yml  # CI migration workflow
```

---

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `TURSO_DATABASE_URL` | Your Turso database URL | ✅ | - |
| `TURSO_AUTH_TOKEN` | Your Turso authentication token | ✅ | - |
| `MIGRATIONS_DIR` | Directory containing migration files | ❌ | `./migrations` |

### Example .env file

```bash
# Turso Configuration
TURSO_DATABASE_URL=libsql://my-awesome-app-db.turso.io
TURSO_AUTH_TOKEN=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9...
MIGRATIONS_DIR=./db/migrations
```

---

## Design Decisions & Limitations

### Design Decisions

1. **Single file per migration**: Both UP and DOWN SQL in one file for simplicity
2. **Auto-incremented versioning**: Eliminates merge conflicts from manual numbering  
3. **Transaction per migration**: Each migration runs in its own transaction
4. **Fail-fast approach**: Stop on first migration error to maintain consistency
5. **Docker-first**: Designed for containerized deployment from day one

### Current Capabilities

1. **Turso-optimized**: Built specifically for libSQL databases and Turso's infrastructure
2. **Simple migrations**: Ordered execution with automatic version management
3. **Safe transactions**: Each migration runs in isolation for data integrity
4. **Container-ready**: Works great in Docker environments
5. **Go-friendly**: Easy integration with Go applications and toolchains

### Roadmap

- Multiple database rollback (`down --count N`)
- Migration dependency system
- Dry-run mode with validation
- Schema drift detection
- Turso-specific optimizations (edge locations, sync strategies)

---

## Contributing

We welcome contributions from the community! turso-migrate is an open-source project and we'd love your help making it better for everyone using Turso databases.

### How to Contribute

#### 🐛 Bug Reports
- Search [existing issues](https://github.com/rubenmeza/turso-migrate/issues) first
- Use our [bug report template](https://github.com/rubenmeza/turso-migrate/issues/new?template=bug_report.md)
- Include reproduction steps and environment details
- Provide Turso database configuration (without sensitive tokens)

#### 💡 Feature Requests
- Check [existing feature requests](https://github.com/rubenmeza/turso-migrate/labels/enhancement)
- Use our [feature request template](https://github.com/rubenmeza/turso-migrate/issues/new?template=feature_request.md)
- Explain the use case and how it benefits Turso users
- Consider implementation complexity and maintenance overhead

#### 🔧 Code Contributions

1. **Fork & Clone**
   ```bash
   git clone https://github.com/your-username/turso-migrate
   cd turso-migrate
   ```

2. **Setup Development Environment**
   ```bash
   go mod tidy
   make dev-setup
   ```

3. **Create Feature Branch**
   ```bash
   git checkout -b feature/amazing-turso-feature
   ```

4. **Make Your Changes**
   - Follow Go best practices and project conventions
   - Add tests for new functionality
   - Update documentation if needed
   - Ensure compatibility with Turso's libSQL

5. **Test Your Changes**
   ```bash
   # Run unit tests
   go test ./...
   
   # Test with real Turso database
   export TURSO_DATABASE_URL=your-test-db
   export TURSO_AUTH_TOKEN=your-test-token
   make ci-migrate
   ```

6. **Commit & Push**
   ```bash
   git add .
   git commit -m "feat: add amazing Turso feature"
   git push origin feature/amazing-turso-feature
   ```

7. **Create Pull Request**
   - Use our [PR template](https://github.com/rubenmeza/turso-migrate/blob/main/.github/pull_request_template.md)
   - Link related issues
   - Explain what changed and why
   - Include testing steps
   - Add screenshots for UI changes

### Development Guidelines

#### Code Style
- Follow standard Go formatting (`gofmt`, `goimports`)
- Use meaningful variable and function names
- Add godoc comments for public APIs
- Keep functions focused and testable

#### Testing
- Write unit tests for new functionality
- Include integration tests for Turso-specific features
- Test Docker builds and CI workflows
- Verify compatibility with different Turso configurations

#### Documentation
- Update README.md for user-facing changes
- Add godoc comments for new APIs
- Include examples in `/examples` directory
- Update CLI help text and descriptions

#### Turso-Specific Considerations
- Ensure compatibility with libSQL features
- Test with Turso's edge replication
- Consider Turso's authentication mechanisms
- Optimize for Turso's performance characteristics

### Community Guidelines

- **Be helpful**: We're all here to make Turso migrations easier
- **Be constructive**: Provide useful feedback and suggestions  
- **Be patient**: Maintainers review contributions when they can
- **Be clear**: Well-explained issues and PRs are easier to help with

### Recognition

Contributors are recognized in several ways:
- Listed in [CONTRIBUTORS.md](CONTRIBUTORS.md)
- Mentioned in release notes for significant contributions
- Given credit in documentation and examples they create
- Invited to join the maintainer team for sustained contributions

### Getting Help

- **Questions**: Use [GitHub Discussions](https://github.com/rubenmeza/turso-migrate/discussions)
- **Real-time chat**: Join our community Discord (link coming soon)
- **Documentation**: Check this README and inline help (`turso-migrate --help`)
- **Examples**: Browse the `/examples` directory

### Development Roadmap

Want to contribute but not sure where to start? Check out these areas:

#### 🔥 High Priority
- Multiple database rollback support (`down --count N`)
- Dry-run mode for migration validation
- Enhanced error messages and helpful suggestions
- Performance optimizations for Turso

#### 🌟 Medium Priority
- Migration dependency resolution system
- Schema drift detection and reporting
- Integration with Turso's branching features
- Simple web UI for migration management

#### 💡 Future Ideas
- Migration templates and generators
- Automated schema documentation
- Integration with Turso's analytics
- Multi-region deployment strategies

### License

By contributing to turso-migrate, you agree that your contributions will be licensed under the same [MIT License](LICENSE) that covers the project.

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Support

- **Issues**: [GitHub Issues](https://github.com/rubenmeza/turso-migrate/issues)
- **Discussions**: [GitHub Discussions](https://github.com/rubenmeza/turso-migrate/discussions)
- **Documentation**: This README and inline help (`turso-migrate --help`)

---

**Ready to get started with simple Turso migrations?** 🚀

Start with: `make migration MIGRATION_NAME=your_first_migration`