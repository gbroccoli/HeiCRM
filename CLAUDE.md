# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HeiCRM is a microservices-based CRM system written in Go, currently featuring an authentication service with JWT-based access/refresh tokens. The architecture uses a monorepo with Go workspaces, shared packages, and service-specific modules.

## Development Commands

### Infrastructure
```bash
# Start PostgreSQL and NATS
docker compose up -d db nats

# PostgreSQL is exposed on localhost:5432
# NATS is exposed on localhost:4222 (client) and localhost:8222 (monitoring)
```

### Database Migrations
Requires goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`

```bash
# Apply all migrations
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" up

# Rollback last migration
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" down

# Check migration status
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" status

# Create new migration
goose -dir migrations create migration_name sql

# Alternatively, use Makefile (Linux/Mac) or scripts/migrate.sh
make migrate-up
./scripts/migrate.sh up
```

### Auth Service
```bash
# Run the auth service (listens on :8080)
go run ./services/auth/cmd/auth

# Build auth service binary
go build -o bin/auth ./services/auth/cmd/auth
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific test
go test ./services/auth/... -run TestLogin -v

# Run tests for a specific package
go test ./pkg/jwt/... -v
```

### Code Quality
```bash
# Format all Go code
go fmt ./...

# Run Go vet for static analysis
go vet ./...
```

## Architecture

### Go Workspace Structure
The project uses Go workspaces (`go.work`) to manage multiple modules:
- Each package in `pkg/` has its own `go.mod`
- Each service in `services/` has its own `go.mod`
- This allows independent versioning while maintaining local development efficiency

### Shared Packages (`pkg/`)
- **config**: Global YAML configuration loader with singleton pattern (`config.G()`)
- **dbx**: PostgreSQL connection manager with DSN builder and singleton pattern (`dbx.G()`)
- **logx**: File-based logging initialization
- **jwt**: JWT token generation/verification for access and refresh tokens (HS512 signing)
- **password**: Password hashing with bcrypt and random password generation

### Services Architecture
Each service follows a layered structure:
```
services/{service_name}/
├── cmd/{service_name}/main.go    # Entry point
├── internal/
│   ├── handler/                  # HTTP handlers (Gin)
│   ├── routes/                   # Route registration
│   ├── middleware/               # Middleware (e.g., AuthMiddleware)
│   └── tools/                    # Service-specific utilities
└── go.mod
```

### Auth Service Flow
1. **Initialization**: `main.go` loads config, initializes logger, connects to DB, creates Gin router
2. **Handler Layer**: `handler.Handler` holds DB, JWT, and PasswordManager dependencies
3. **Routes**: Mounted at `/auth` group with:
   - `POST /auth/login` - Public endpoint
   - `POST /auth/register` - Protected by AuthMiddleware
4. **Middleware**: `AuthMiddleware` extracts Bearer token from Authorization header and verifies using JWT.VerifyAccessToken

### Configuration
- Primary config: `config.yaml` (gitignored, created from `config.example.yaml`)
- Config structure includes: database, JWT secrets, NATS, Redis, email (planned)
- Global config accessed via `config.G()` after `config.MustLoad("config.yaml")`

### Database
- PostgreSQL managed via `database/sql` (no ORM)
- Connection pooling configured in `dbx.Open()`:
  - MaxOpenConns: 15
  - MaxIdleConns: 10
  - ConnMaxLifetime: 45m
  - ConnMaxIdleTime: 5m
- Migrations managed with [goose](https://github.com/pressly/goose) in `migrations/` directory
- Migration format: SQL files with `-- +goose Up` and `-- +goose Down` sections
- Current schema: `users` table with email (unique), password (bcrypt), role, timestamps

### JWT Token Strategy
- **Access tokens**: 30-minute expiry, include email and role, type="access"
- **Refresh tokens**: 30-day expiry, email only, type="refresh"
- Both use HS512 signing algorithm
- Claims include custom fields (Email, Role, Type) plus standard JWT claims
- Token verification checks token type to prevent cross-use

## Key Patterns

### Singleton Pattern for Global State
Multiple packages use singleton pattern with `G()` getter:
- `config.G()` - global configuration
- `dbx.G()` - global database connection
- Always check for nil and panic if accessed before initialization

### Dependency Injection in Handlers
Handler constructors receive dependencies explicitly:
```go
handler.New(dbx.G(), jwt.New([]byte(config.G().Jwt.SecretKey)))
```

### Gin Framework Conventions
- Use `c.ShouldBindJSON()` for request validation
- JSON tags with `binding:"required"` or `binding:"required,email"`
- Return errors with `c.JSON(statusCode, gin.H{"error": "message"})`
- Set response headers before JSON response (e.g., Authorization header)

## Important Notes

### Security
- JWT secret keys must be set in `config.yaml` before running
- Private keys (`.pem`, `.key`, `.crt`) are gitignored
- Never commit actual credentials to `config.yaml`
- Login handler currently doesn't verify passwords against database (incomplete implementation)

### Current Limitations
- Login handler generates tokens without database verification (development stub)
- Register endpoint requires authentication (unusual pattern, may be intentional for admin-only registration)
- Error messages contain Russian text mixed with English responses

### Windows Development
- Line endings are CRLF (configured in `.editorconfig`)
- Go uses tabs, YAML/JSON use 2 spaces
- Logs written to `logs/` directory (gitignored)

## Module Import Path
All internal imports use: `github.com/gbroccoli/HeiCRM/{package}`
