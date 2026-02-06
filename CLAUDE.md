# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HeiCRM is a microservices-based CRM system written in Go, featuring authentication, user management, and an API gateway. The architecture uses a monorepo with Go workspaces, shared packages, and service-specific modules.

## Development Commands

### Infrastructure
```bash
# Start PostgreSQL, NATS, and Redis
docker compose up -d db nats redis

# PostgreSQL is exposed on localhost:5432
# NATS is exposed on localhost:4222 (client) and localhost:8222 (monitoring)
# Redis is on the internal heicrm network (port 6379)
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
```

### Running Services
```bash
# API Gateway (port 8000) — entry point for all requests
go run ./services/apigateway/cmd/gateway

# Auth Service (port 8080)
go run ./services/auth/cmd/auth

# User Service (port 8081)
go run ./services/users/cmd/users

# Housing Service (port 8082)
go run ./services/housing/cmd/housing

# Tasks Service (port 8083)
go run ./services/tasks/cmd/tasks
```

### Building
```bash
# Build individual services
go build -o bin/auth ./services/auth/cmd/auth
go build -o bin/users ./services/users/cmd/users
go build -o bin/gateway ./services/apigateway/cmd/gateway

# Build all (from workspace root)
go build ./services/auth/cmd/auth
go build ./services/users/cmd/users
go build ./services/apigateway/cmd/gateway
```

### Testing
```bash
# Run tests for a specific service
go test ./services/auth/... -v
go test ./services/users/... -v

# Run tests for a specific package
go test ./pkg/jwt/... -v

# Run tests with coverage
go test ./services/auth/... -cover
```

### Code Quality
```bash
# Vet per module (workspace doesn't support ./...)
go vet ./services/auth/...
go vet ./services/users/...
go vet ./services/apigateway/...
go vet ./pkg/natsx/...
go vet ./pkg/models/...

# Format per module
cd services/users && go fmt ./...
cd pkg/natsx && go fmt ./...

# Sync workspace dependencies
go work sync
```

## Architecture

### Go Workspace Structure
The project uses Go workspaces (`go.work`) to manage multiple modules:
- Each package in `pkg/` has its own `go.mod`
- Each service in `services/` has its own `go.mod`
- Local packages are resolved via `go.work`, not published to remote
- New local packages require `go work sync` after adding to `go.work`

Current workspace modules:
```
pkg/config, pkg/dbx, pkg/jwt, pkg/logx, pkg/models,
pkg/natsx, pkg/password, pkg/redisx, pkg/response
services/apigateway, services/auth, services/users, services/housing, services/tasks
```

### Shared Packages (`pkg/`)
- **config**: Global YAML configuration loader with singleton pattern (`config.G()`)
- **dbx**: PostgreSQL connection manager with DSN builder and singleton pattern (`dbx.G()`)
- **logx**: File-based logging initialization
- **jwt**: JWT token generation/verification for access and refresh tokens (HS512 signing)
- **password**: Password hashing with bcrypt and random password generation
- **redisx**: Redis client with singleton pattern (`redisx.G()`)
- **natsx**: NATS connection manager with singleton pattern (`natsx.G()`)
- **models**: Shared data structures (User, Role, UserProfile, UserWithProfile, ActivityLog, pagination)
- **response**: Standard API response helpers with application codes (1xxx-9xxx)

### Services

#### API Gateway (port 8000)
- Reverse proxy to all microservices
- Centralized CORS configuration
- Routes: `/api/v1/auth/*` → Auth Service, `/api/v1/users/*` → User Service
- Health check: `GET /api/v1/health`

#### Auth Service (port 8080)
- JWT-based authentication with access/refresh tokens
- Refresh token rotation with Redis-backed sessions
- HTTPOnly cookie for refresh tokens
- Endpoints: login, register, refresh, logout, me
- Dependencies: DB, JWT, Redis, Password

#### User Service (port 8081)
- CRUD operations for user profiles
- Role-based access control (admin, manager, user)
- NATS integration for event-driven profile creation
- Endpoints: get/update own profile, list/get/update/delete users, activity log
- Dependencies: DB, JWT, NATS

#### Housing Service (port 8082)
- Management of buildings, rooms, and residents
- Room types: single, double, block
- Room status: free, occupied (auto-updated based on occupancy)
- Resident management with move-in/move-out dates
- Endpoints: buildings CRUD, rooms CRUD, residents assignment
- Dependencies: DB, JWT

#### Tasks Service (port 8083)
- Task/request management system
- Task types: custom (e.g., Ремонт, Уборка, IT-поддержка)
- Priorities: low, medium, high, critical
- Statuses: new, assigned, in_progress, completed, closed
- Task assignment, comments, status history
- Endpoints: tasks CRUD, status updates, assign/take, comments, history
- Dependencies: DB, JWT

### Service Architecture Pattern
Each service follows a layered structure:
```
services/{service_name}/
├── cmd/{service_name}/main.go    # Entry point
├── internal/
│   ├── handler/                  # HTTP handlers (Gin)
│   │   ├── handler.go            # Handler struct + New()
│   │   └── ...                   # One file per endpoint
│   ├── routes/                   # Route registration (Mount function)
│   ├── middleware/               # Auth, role middleware
│   └── nats/                     # NATS subscribers (if applicable)
└── go.mod
```

### Main Entry Point Pattern
Every service follows the same initialization sequence:
1. `logx.Init("logs/<service>.log")` — file logging
2. `config.MustLoad("config.yaml")` — load configuration
3. `gin.Default()` with Recovery and Logger
4. `dbx.Open()` — connect to PostgreSQL
5. Service-specific connections (redisx, natsx)
6. `jwt.New([]byte(config.G().Jwt.SecretKey))` — init JWT
7. `handler.New(...)` — create handler with dependencies
8. NATS subscribers (if applicable)
9. `routes.Mount(...)` — register routes
10. `g.Run(":port")` — start HTTP server

### Configuration
- Primary config: `config.yaml` (gitignored, created from `config.example.yaml`)
- Config structure: database, JWT, cookie, Redis, NATS, email (planned), serves (service URLs)
- Global config accessed via `config.G()` after `config.MustLoad("config.yaml")`
- Service URLs configured in `serves` section:
  ```yaml
  serves:
    auth: "http://localhost:8080"
    users: "http://localhost:8081"
    housing: "http://localhost:8082"
    tasks: "http://localhost:8083"
  ```

### Database
- PostgreSQL managed via `database/sql` (no ORM)
- Connection pooling configured in `dbx.Open()`:
  - MaxOpenConns: 15
  - MaxIdleConns: 10
  - ConnMaxLifetime: 45m
  - ConnMaxIdleTime: 5m
- Migrations managed with [goose](https://github.com/pressly/goose) in `migrations/` directory
- Migration format: SQL files with `-- +goose Up` and `-- +goose Down` sections
- Current schema:
  - `users` — id, name, email (unique), password (bcrypt), avatar, role_id (FK), tg_send, timestamps
  - `roles` — id, name (unique), description; defaults: 0=user, 1=admin, 2=manager
  - `user_profiles` — user_id (unique FK), first_name, last_name, middle_name, phone, student_id, room_id, avatar_url, date_of_birth, timestamps
  - `user_activity_log` — user_id (FK), action, details (JSONB), ip_address, user_agent, created_at
  - `buildings` — id, address, floors, description, timestamps
  - `rooms` — id, building_id (FK), room_number, floor, capacity, room_type (single/double/block), status (free/occupied), timestamps
  - `residents` — id, room_id (FK), full_name, birth_date, passport_series, passport_number, email, phone, move_in_date, move_out_date, timestamps
  - `tasks` — id, author_id, assignee_id, room_id, task_type, description, priority, status, timestamps
  - `task_history` — id, task_id (FK), previous_status, new_status, changed_by, changed_at, comment
  - `task_comments` — id, task_id (FK), author_id, comment_text, created_at
  - `task_attachments` — id, task_id (FK), file_name, file_path, file_size, uploaded_by, uploaded_at

### JWT Token Strategy
- **Access tokens**: 45-minute expiry, include email (Subject) and role, type="access"
- **Refresh tokens**: 30-day expiry, email and role, type="refresh"
- Both use HS512 signing algorithm
- `claims.ID` is a random UUID (jti), **not** the database user ID
- `claims.Subject` contains the user's email — use this to identify users
- Token verification checks token type to prevent cross-use
- Refresh tokens stored in Redis with SHA256 hash as key

### NATS Events
- User Service subscribes to `user.registered` — creates profile on new registration
- User Service publishes `user.profile_updated` and `user.deactivated`
- Connection via `natsx.G()` singleton

### Roles
| role_id | Name    | Access Level |
|---------|---------|-------------|
| 0       | user    | Own profile only |
| 1       | admin   | Full CRUD, activity logs, delete users |
| 2       | manager | Read-only access to user list |

## Key Patterns

### Singleton Pattern for Global State
Multiple packages use singleton pattern with `G()` getter:
- `config.G()` — global configuration
- `dbx.G()` — global database connection
- `redisx.G()` — global Redis client
- `natsx.G()` — global NATS connection
- Always check for nil and panic if accessed before initialization

### Dependency Injection in Handlers
Handler constructors receive dependencies explicitly:
```go
// Auth service — needs DB, JWT, Redis
handler.New(dbx.G(), jwt.New([]byte(config.G().Jwt.SecretKey)), redisx.G())

// User service — needs DB, JWT only
handler.New(dbx.G(), jwt.New([]byte(config.G().Jwt.SecretKey)))
```

### User Identity Resolution
JWT tokens do NOT contain the database user ID. The `claims.ID` field is a random UUID (jti).
To get the database user ID, always resolve via email:
```go
email, _ := c.Get("email")            // from JWT Subject
userID, err := getUserIDByEmail(db, email.(string))  // DB lookup
```
The auth service's `/me` handler uses the same pattern.

### Gin Framework Conventions
- Use `c.ShouldBindJSON()` for request validation
- JSON tags with `binding:"required"` or `binding:"required,email"`
- Standardized responses via `pkg/response` helpers
- Role checks via middleware or in-handler logic
- `/me` routes registered before `/:id` routes to prevent Gin treating "me" as a parameter

### Response Format
All responses use the standard envelope:
```json
{
  "code": 1000,
  "message": "Success",
  "data": { ... },
  "error": "..."
}
```
Application codes: 1xxx (success), 2xxx (validation), 3xxx (auth), 4xxx (access), 5xxx (resources), 9xxx (system).

## Important Notes

### Security
- JWT secret keys must be set in `config.yaml` before running
- Private keys (`.pem`, `.key`, `.crt`) are gitignored
- Never commit actual credentials to `config.yaml`
- Refresh tokens use SHA256 hashing for Redis keys (never stored as plain JWT)
- HTTPOnly cookies for refresh tokens (not accessible via JavaScript)
- HTTPS detection via `X-Forwarded-Proto` header for tunnel/proxy support

### Current Limitations
- Register endpoint is currently public (no auth required)
- NATS event publishing (`user.profile_updated`, `user.deactivated`) is defined but not yet called from handlers
- Auth Service needs to publish `user.registered` event after registration for User Service profile auto-creation
- `go mod tidy` fails for modules using `pkg/natsx` (not published to remote) — use `go work sync` instead
- Error messages mix Russian and English in some places

### Windows Development
- Line endings are CRLF (configured in `.editorconfig`)
- Go uses tabs, YAML/JSON use 2 spaces
- Logs written to `logs/` directory (gitignored)

## API Documentation

See [ENDPOINTS.md](ENDPOINTS.md) for full API endpoint documentation with request/response examples.

## Module Import Path
All internal imports use: `github.com/gbroccoli/HeiCRM/{package}`
