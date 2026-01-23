# CLAUDE.md

This file provides guidance for Claude when working with the go-aula codebase.

## Project Overview

Go-aula is a Go microservice demonstrating clean architecture principles. It provides a RESTful API for managing Tags and Categories (similar to a product/commerce system).

## Build & Test Commands

```bash
make start        # Run the application: go run ./cmd/api/main.go
make test         # Run unit tests: go test ./...
make e2e          # Run e2e tests: go test ./... -tags e2e
make fmt          # Format code: go fmt ./...
make tidy         # Download dependencies: go mod tidy
make update       # Update libs: go get -u ./...
```

## Architecture

This project follows **Clean Architecture** with 3 layers:

1. **Handler** (`domain/handler/`) - HTTP request/response, validation
2. **Service** (`domain/service/`) - Business logic
3. **Repository** (`domain/repository/`) - Data access (GORM)

Interfaces are defined in `domain/interface/` and implementations are in their respective packages.

## Key Directories

```
cmd/api/main.go          # Entry point (Uber Fx DI)
cmd/infra/               # Infrastructure (HTTP, DB, config, logging, telemetry)
domain/entitie/          # Data models (Tag, Category)
domain/interface/        # Contracts and custom errors
domain/interface/mocks/  # Auto-generated mocks (mockery)
domain/handler/          # HTTP handlers
domain/service/          # Business logic
domain/repository/       # Data access
config/local/            # Local dev configuration
```

## Tech Stack

- **Go 1.20**, **Echo v4** (web framework), **GORM** (ORM)
- **Uber Fx** (dependency injection), **Viper** (configuration)
- **OpenTelemetry** (tracing/metrics), **Prometheus**, **Zap** (logging)
- **CockroachDB/PostgreSQL** (database), **Flyway** (migrations)
- **Testify**, **Mockery**, **go-sqlmock**, **TestContainers** (testing)

## Code Conventions

### Naming
- Interfaces: `CategoryService`, `CategoryRepository` (in `domain/interface/`)
- Implementations: `categoryService struct` (lowercase, unexported)
- DI providers: `Provide<Type>Name()` functions
- Route registration: `RegisterHooks()` functions

### Dependency Injection (Uber Fx)
- Use `fx.Module()`, `fx.Provide()`, `fx.Invoke()`
- Lifecycle hooks: `fx.Hook{OnStart, OnStop}` for graceful shutdown
- Constructor injection through function parameters

### Testing
- Table-driven tests with `t.Run()` subtests
- Generate mocks: `//go:generate mockery --name <InterfaceName>`
- Use `mockObj.AssertExpectations(t)` for mock verification
- Use Faker for test data generation
- E2E tests use TestContainers (tag with `e2e`)

### Error Handling
- Custom errors defined in `domain/interface/errors.go`
- Map errors to HTTP status codes in `domain/handler/helper/handler.go`

### Entities
- Use GORM model embedding (`gorm.Model`)
- Add JSON tags for serialization
- Add validation tags: `validate:"required"`

## API Endpoints

- **Tag**: `PUT/POST /tag`, `GET /tag/{id}`, `DELETE /tag/{id}`, `GET /tag/getAll`
- **Category**: `PUT/POST /category`, `GET /category/{id}`, `DELETE /category/{id}`, `GET /category/getAll`
- **Metrics**: `/metrics` (Prometheus)
- **Profiling**: `/debug/pprof/*`

## Local Development

1. Start database: `docker-compose up`
2. Run application: `make start` (serves on `:9090`)

## Adding New Features

When adding a new domain entity (e.g., Product):

1. Create entity in `domain/entitie/product.go`
2. Define interfaces in `domain/interface/product.go`
3. Implement repository in `domain/repository/product/`
4. Implement service in `domain/service/product/`
5. Implement handler in `domain/handler/product/`
6. Register in Fx modules in `cmd/api/main.go`
7. Add migration in `scripts/flyway/sql/`
8. Update `openapi.yaml` with new endpoints
