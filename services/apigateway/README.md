# API Gateway

API Gateway for HeiCRM microservices architecture. Provides centralized routing, CORS handling, and request proxying to backend services.

## Architecture

```
Frontend (localhost:3000)
    ↓
API Gateway (localhost:8000)
    ↓
┌─────────────────┐
│ Auth Service    │ :8080
│ User Service    │ :8081 (future)
│ Ticket Service  │ :8082 (future)
└─────────────────┘
```

## Running the Gateway

### Development
```bash
# Start API Gateway (default port: 8000)
go run ./services/apigateway/cmd/gateway

# With custom configuration
AUTH_SERVICE_URL=http://localhost:8080 GATEWAY_PORT=8000 go run ./services/apigateway/cmd/gateway
```

### Build
```bash
cd services/apigateway
go build -o bin/gateway ./cmd/gateway
./bin/gateway
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GATEWAY_PORT` | `8000` | API Gateway listening port |
| `AUTH_SERVICE_URL` | `http://localhost:8080` | Auth microservice URL |

## API Routes

### Health Check
```
GET /api/v1/health
```
Response:
```json
{
  "status": "ok",
  "service": "api-gateway"
}
```

### Auth Service (proxied)
All requests to `/api/v1/auth/*` are proxied to the auth service:

```
POST   /api/v1/auth/login     → http://localhost:8080/auth/login
POST   /api/v1/auth/register  → http://localhost:8080/auth/register
POST   /api/v1/auth/refresh   → http://localhost:8080/auth/refresh
POST   /api/v1/auth/logout    → http://localhost:8080/auth/logout
GET    /api/v1/auth/me        → http://localhost:8080/auth/me
```

## CORS Configuration

CORS is centralized at the API Gateway level. Allowed origins:
- `http://localhost:3000` (React/Next.js)
- `http://localhost:5173` (Vite)
- `http://localhost:4200` (Angular)
- `http://localhost:8081`

**For production**: Add your production domain in `cmd/gateway/main.go`:
```go
AllowOrigins: []string{
    "http://localhost:3000",
    "https://crm.yourdomain.com", // Add this
}
```

## Adding New Services

To add a new microservice:

1. **Update `routes/routes.go`**:
```go
type ServiceConfig struct {
    AuthServiceURL string
    UserServiceURL string  // Add new service
}

func Mount(r *gin.Engine, config ServiceConfig) {
    api := r.Group("/api/v1")
    {
        // Add new service route
        users := api.Group("/users")
        {
            users.Any("/*path", proxy.ReverseProxy(config.UserServiceURL))
        }
    }
}
```

2. **Update `cmd/gateway/main.go`**:
```go
serviceConfig := routes.ServiceConfig{
    AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
    UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8081"), // Add
}
```

3. **Set environment variable** (optional):
```bash
USER_SERVICE_URL=http://localhost:8081 go run ./cmd/gateway
```

## Development Workflow

### Running all services

```bash
# Terminal 1: Start PostgreSQL & Redis
docker compose up -d db redis

# Terminal 2: Start Auth Service
go run ./services/auth/cmd/auth

# Terminal 3: Start API Gateway
go run ./services/apigateway/cmd/gateway

# Terminal 4: Frontend
cd frontend && npm run dev
```

### Testing API Gateway

```bash
# Health check
curl http://localhost:8000/api/v1/health

# Login through gateway
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

## Security Features

- ✅ Centralized CORS handling
- ✅ HTTPOnly cookie support (AllowCredentials)
- ✅ Request/Response logging
- ✅ Error handling for unavailable services
- 🔄 Rate limiting (planned)
- 🔄 JWT validation at gateway level (planned)
- 🔄 Request tracing (planned)

## Future Improvements

- [ ] Rate limiting middleware
- [ ] JWT validation at gateway (optional)
- [ ] Request/Response transformation
- [ ] Caching layer
- [ ] Load balancing for multiple service instances
- [ ] Circuit breaker pattern
- [ ] Metrics and monitoring (Prometheus)
- [ ] Distributed tracing (Jaeger/OpenTelemetry)
