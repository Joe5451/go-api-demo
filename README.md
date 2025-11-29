# Go API Demo

[![Go Report Card](https://goreportcard.com/badge/github.com/Joe5451/go-api-demo)](https://goreportcard.com/report/github.com/Joe5451/go-api-demo)

A simple CRUD API demo built with **Go** and the **Gin** web framework, using **PostgreSQL** for persistence and **Swagger** (swaggo) for API documentation.

## Prerequisites

- Go 1.25.1
- Docker
  - For local PostgreSQL via `docker compose`
  - Required for feature tests (uses testcontainers-go)

## Configuration

This project reads configuration from environment variables. A local `.env` file is optional (gitignored) and will be loaded if present.

Create `.env` from the provided template:

```bash
cp .env.example .env
```

`.env.example`:

```dotenv
DEBUG=true

# postgres
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DBNAME=book
POSTGRES_SCHEMA=public
```

## Project Layout

```text
.
├── cmd/
├── docs/
├── internal/
│   ├── adapter/
│   ├── application/
│   │   └── port/
│   ├── bootstrap/
│   ├── config/
│   ├── constant/
│   ├── domain/
│   ├── http/
│   └── infra/
├── mocks/
├── test/
├── Dockerfile
├── docker-compose.yml
└── init.sql
```

| Directory/File | Description |
|----------------|-------------|
| `cmd/` | Application entrypoint (Gin + Swagger route) |
| `docs/` | Generated Swagger artifacts (swaggo) |
| `internal/adapter/` | Implementations of ports (handlers, repositories) |
| `internal/application/` | Use cases (business flows) |
| `internal/application/port/` | Interfaces (in/out) for dependency inversion |
| `internal/bootstrap/` | Dependency injection wiring |
| `internal/config/` | Config loading and structs |
| `internal/constant/` | Shared error codes/constants |
| `internal/domain/` | Entities + domain rules (no dependencies on other layers) |
| `internal/http/` | HTTP routes, middleware, HTTP helpers |
| `internal/infra/` | Infrastructure (Postgres pool) |
| `mocks/` | Generated mocks (go.uber.org/mock) |
| `test/` | Feature tests (testcontainers + httptest) |
| `Dockerfile` | Container build for the app |
| `docker-compose.yml` | Postgres (and optional app) for local/dev |
| `init.sql` | DB schema init for Postgres containers |

This project follows Clean Architecture. Dependencies point inward:

- **Domain** (`internal/domain`): entities and domain errors (pure Go, no frameworks)
- **Application** (`internal/application` + `internal/application/port`): use cases depend on interfaces, not implementations
- **Adapters** (`internal/adapter`): HTTP handlers and repository implementations that satisfy ports
- **Infra** (`internal/infra`): database connection setup (pgxpool)
- **Bootstrap** (`internal/bootstrap`): wires everything together

This keeps business logic testable and decoupled from transport (HTTP) and infrastructure (Postgres).

## Run

### 1. Start PostgreSQL (via Docker)

Set env vars (e.g. in `.env`) and start Postgres:

```bash
docker compose up -d postgres
```

The database schema is initialized from `init.sql`.

### 2. Run the API server

```bash
go run ./cmd/main.go
```

Server will listen on:
- `http://localhost:8080`

Swagger UI:
- `http://localhost:8080/swagger/index.html`

## API Endpoints

Base URL: `http://localhost:8080`

Books:
- `POST /books`
- `GET /books/:id`
- `GET /books?page=1&per_page=10`
- `PUT /books/:id`
- `DELETE /books/:id`

Example (create a book):

```bash
curl -i -X POST "http://localhost:8080/books" \
  -H "Content-Type: application/json" \
  -d '{"title":"1984","author":"George Orwell"}'
```

## Test

Unit tests (focused on internal layers):

```bash
go test ./internal/...
```

All tests:

```bash
go test ./... -count=1
```

Feature tests (requires Docker; uses [testcontainers](https://github.com/testcontainers/testcontainers-go)):

```bash
go test ./test/api/... -v
```

## Swagger

Install `swag` and regenerate:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go
```
