# Quick Start Guide

Быстрая инструкция по запуску проекта HeiCRM.

## 1. Предварительные требования

- Go 1.25.3 или выше
- Docker и Docker Compose
- Git

## 2. Установка зависимостей

```bash
# Установить goose для миграций
go install github.com/pressly/goose/v3/cmd/goose@latest

# Убедиться, что goose добавлен в PATH
# Linux/Mac: добавьте в ~/.bashrc или ~/.zshrc
# export PATH=$PATH:$(go env GOPATH)/bin

# Windows: добавьте %USERPROFILE%\go\bin в PATH
```

## 3. Настройка конфигурации

```bash
# Скопировать пример конфигурации
cp config.example.yaml config.yaml

# Отредактировать config.yaml
# Обязательно установите:
# - jwt.secret_key (длинная случайная строка)
# - database.name (например, "crm")
```

Пример `config.yaml`:
```yaml
env: dev
service_name: HeiCRM

jwt:
  secret_key: "your-super-secret-key-min-32-chars-long"
  alg: HS256

database:
  host: localhost
  port: 5432
  user: postgres
  password: root
  name: crm
  sslmode: disable

access_token_ttl: "30m"
refresh_token_ttl: "720h"

nats:
  url: "nats://localhost:4222"

redis:
  host: localhost
  port: 6379
  password: ""
```

## 4. Запуск инфраструктуры

```bash
# Запустить PostgreSQL и NATS
docker compose up -d db nats

# Проверить, что контейнеры запущены
docker ps

# Дождаться готовности PostgreSQL (обычно 5-10 секунд)
# Можно проверить логи:
docker logs database
```

## 5. Создание базы данных

```bash
# Подключиться к PostgreSQL
docker exec -it database psql -U postgres

# В консоли psql создать базу данных:
CREATE DATABASE crm;

# Выйти из psql
\q
```

## 6. Применение миграций

### Windows:
```bash
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" up
```

### Linux/Mac:
```bash
# Вариант 1: через Makefile
make migrate-up

# Вариант 2: через скрипт
./scripts/migrate.sh up

# Вариант 3: напрямую через goose
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" up
```

Проверить статус миграций:
```bash
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" status
```

Ожидаемый вывод:
```
    Applied At                  Migration
    =======================================
    <timestamp>                 01_user_tables.sql
```

## 7. Запуск сервиса авторизации

```bash
# Запустить auth service
go run ./services/auth/cmd/auth

# Или собрать и запустить бинарник
go build -o bin/auth ./services/auth/cmd/auth
./bin/auth  # Windows: bin\auth.exe
```

Сервис будет доступен на `http://localhost:8080`

## 8. Проверка работоспособности

### Тест endpoint'а login:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Ожидаемый ответ:
```json
{
  "message": "success login",
  "data": {
    "email": "test@example.com",
    "password": "password123"
  }
}
```

В заголовках ответа будет токен:
```
Authorization: Bearer <jwt_token>
```

## 9. Структура проекта

```
HeiCRM/
├── config.yaml                 # Конфигурация (не в git)
├── config.example.yaml         # Пример конфигурации
├── docker-compose.yml          # PostgreSQL + NATS
├── migrations/                 # SQL миграции
├── pkg/                        # Общие пакеты
│   ├── config/                # Загрузка конфигурации
│   ├── dbx/                   # Подключение к БД
│   ├── jwt/                   # JWT токены
│   ├── logx/                  # Логирование
│   └── password/              # Хеширование паролей
└── services/
    └── auth/                  # Сервис авторизации
        ├── cmd/auth/          # Entry point
        └── internal/          # Внутренняя логика
            ├── handler/       # HTTP handlers
            ├── middleware/    # Middleware
            └── routes/        # Роутинг
```

## 10. Полезные команды

```bash
# Логи базы данных
docker logs database

# Логи NATS
docker logs nats

# Подключиться к базе данных
docker exec -it database psql -U postgres -d crm

# Остановить инфраструктуру
docker compose down

# Остановить с удалением данных
docker compose down -v

# Форматирование кода
go fmt ./...

# Статический анализ
go vet ./...

# Запуск тестов
go test ./...
```

## Troubleshooting

### Ошибка "connection refused" при запуске сервиса
- Проверьте, что PostgreSQL запущен: `docker ps`
- Проверьте подключение: `docker exec -it database psql -U postgres -c "SELECT 1"`

### Ошибка "database does not exist"
- Создайте базу данных: `docker exec -it database psql -U postgres -c "CREATE DATABASE crm;"`

### Ошибка "goose: command not found"
- Установите goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`
- Добавьте `$(go env GOPATH)/bin` в PATH

### Ошибка "panic: global config is nil"
- Убедитесь, что файл `config.yaml` существует
- Проверьте формат YAML (пробелы, отступы)

### Ошибка при применении миграций
- Проверьте параметры подключения к БД в команде goose
- Убедитесь, что база данных `crm` создана
