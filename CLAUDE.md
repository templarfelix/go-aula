# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go microservice demonstrating clean architecture patterns for a product/category management API. The project is based on [bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch).

## Tech Stack

- **Language**: Go 1.20
- **HTTP Framework**: Echo v4
- **Dependency Injection**: Uber fx
- **Database**: CockroachDB (PostgreSQL-compatible)
- **ORM**: GORM
- **Migrations**: Flyway
- **Observability**: OpenTelemetry (traces/metrics), zap (logging), pprof (profiling)
- **Testing**: testify, sqlmock, testcontainers-go
- **Validation**: go-playground/validator.v9
- **Configuration**: Viper (YAML configs)

## Project Structure

```
cmd/
  api/main.go          # Application entrypoint
  infra/               # Infrastructure setup (config, database, http, logging, otel)
domain/
  entitie/             # Domain entities (Category, Tag)
  interface/           # Interfaces and mocks for services/repositories
  handler/             # HTTP handlers (Echo routes)
  service/             # Business logic layer
  repository/          # Data access layer (GORM)
config/
  local/config.yaml    # Local development configuration
scripts/
  flyway/              # Database migration scripts
```

## Common Commands

```bash
# Start the application
make start

# Run unit tests
make test

# Run e2e tests (requires Docker for testcontainers)
make e2e

# Format code and tidy dependencies
make fmt

# Update dependencies
make update

# Generate mocks
go generate ./...
```

## API Endpoints

The server runs on `:9090` by default.

### Category
- `POST /category` - Create category
- `PUT /category` - Update category
- `GET /category/:id` - Get category by ID
- `GET /category/getAll` - List all categories
- `DELETE /category/:id` - Delete category

### Tag
- Similar CRUD endpoints at `/tag`

## Development Setup

1. Start CockroachDB and run migrations:
   ```bash
   docker-compose up -d
   ```

2. Run the application:
   ```bash
   make start
   ```

## Configuration

Configuration is loaded via Viper from `config.yaml` in:
- `./config/local/` (local development)
- `./config/cloud/` (cloud deployment)

Key settings:
- `server.address`: HTTP server address (default `:9090`)
- `database.*`: CockroachDB connection settings

## Testing

- **Unit tests**: Use sqlmock for repository tests, mockery-generated mocks for service/handler tests
- **E2E tests**: Use testcontainers-go with build tag `e2e`

Run e2e tests:
```bash
go test ./... -tags e2e
```

## Profiling (pprof)

Available at `http://localhost:9090/debug/pprof/`:
```bash
go tool pprof http://localhost:9090/debug/pprof/heap
go tool pprof http://localhost:9090/debug/pprof/profile?seconds=30
```

## Architecture Notes

- Uses fx modules for dependency injection wiring
- Clean architecture: handlers → services → repositories
- Interfaces defined in `domain/interface/` with mocks in `domain/interface/mocks/`
- GORM models embed `gorm.Model` for ID, timestamps, and soft deletes
