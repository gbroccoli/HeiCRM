# Quick Start Guide

Пошаговая инструкция для запуска HeiCRM в development режиме.

## Предварительные требования

- Go 1.25+
- Docker Desktop (для PostgreSQL и Redis)
- goose (для миграций): `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Шаг 1: Запуск инфраструктуры (PostgreSQL + Redis)

```bash
# Запустить Docker контейнеры
docker compose up -d db redis

# Проверить что контейнеры запущены
docker compose ps
```

Ожидаемый вывод:
```
NAME      IMAGE          STATUS        PORTS
db        postgres:16    Up 5 seconds  0.0.0.0:5432->5432/tcp
redis     redis:7        Up 5 seconds  0.0.0.0:6379->6379/tcp
```

## Шаг 2: Применение миграций БД

```bash
# Применить все миграции
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" up

# Проверить статус
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" status
```

## Шаг 3: Запуск Auth Service

**Открыть новый терминал (Terminal 1)**

```bash
# Перейти в корень проекта
cd C:\Users\zerat\Desktop\HeiCRM

# Запустить auth service
go run ./services/auth/cmd/auth
```

Ожидаемый вывод:
```
starting auth service
PID=12345
Connecting to database
Connecting to Redis
starting http server
[GIN-debug] Listening and serving HTTP on :8080
```

✅ Auth Service слушает на **http://localhost:8080**

## Шаг 4: Запуск API Gateway

**Открыть новый терминал (Terminal 2)**

```bash
# Перейти в корень проекта
cd C:\Users\zerat\Desktop\HeiCRM

# Запустить API Gateway
go run ./services/apigateway/cmd/gateway
```

Ожидаемый вывод:
```
Starting API Gateway
PID=12346
API Gateway listening on :8000
Proxying auth service: http://localhost:8080
[GIN-debug] Listening and serving HTTP on :8000
```

✅ API Gateway слушает на **http://localhost:8000**

## Шаг 5: Проверка работоспособности

### Health Check API Gateway
```bash
curl http://localhost:8000/api/v1/health
```

Ожидаемый ответ:
```json
{"status":"ok","service":"api-gateway"}
```

### Проверка Auth Service через Gateway
```bash
curl -X POST http://localhost:8000/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
```

## Архитектура запросов

```
Frontend (localhost:3000)
    ↓
API Gateway (localhost:8000) ← YOUR FRONTEND SHOULD CALL THIS
    ↓ proxy
Auth Service (localhost:8080) ← Internal
    ↓
PostgreSQL (localhost:5432)
Redis (localhost:6379)
```

## Остановка сервисов

### Остановить Go сервисы
Нажмите `Ctrl+C` в терминалах

### Остановить Docker
```bash
docker compose down
```
