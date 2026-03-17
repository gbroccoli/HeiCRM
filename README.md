# HeiCRM

Система управления заявками для общежитий и кампусов.

**Дипломный проект**

## О проекте

**HeiCRM** (Higher Education Institution CRM — CRM для высших учебных заведений) — это микросервисная платформа для управления общежитиями и студенческими кампусами. Система позволяет жильцам (студентам) подавать заявки на различные услуги, а администрации общежития — эффективно обрабатывать обращения и управлять инфраструктурой.

### Для кого этот проект?

- **Общежития университетов и колледжей**
- **Студенческие кампусы**
- **Общежития предприятий**
- **Арендное жильё с централизованным управлением**

### Основные возможности

- **Подача заявок студентами** — ремонт, уборка, доставка, IT-поддержка и другое
- **Отслеживание статуса** — от создания до выполнения
- **Уведомления в реальном времени** — email, push, in-app
- **Управление зданиями** — структура общежитий, комнаты, этажи
- **Ролевая система** — студенты, персонал, администраторы
- **История обращений** — полный аудит всех заявок

## Архитектура

Проект построен на микросервисной архитектуре с использованием Go Workspaces, PostgreSQL, NATS для межсервисной коммуникации и JWT для авторизации.

### Сервисы

Система состоит из 5 микросервисов:

1. **Auth Service** (✅ Реализован) — аутентификация и авторизация
2. **User Service** (📋 Планируется) — управление профилями пользователей
3. **Building Service** (📋 Планируется) — управление зданиями и комнатами
4. **Task Service** (📋 Планируется) — система заявок
5. **Notification Service** (📋 Планируется) — уведомления (Email, Push, In-app)

📖 **Детальное описание архитектуры:** [SERVICES.md](SERVICES.md)

## Текущее состояние

На данный момент реализован базовый **Auth Service** с поддержкой:

- ✅ Регистрация пользователей с указанием роли
- ✅ Вход в систему с выдачей JWT токенов
- ✅ Access tokens (30 минут) и Refresh tokens (30 дней)
- ✅ Middleware для защиты endpoints
- ✅ Хеширование паролей с bcrypt

## Быстрый старт

### Требования

- **Go** 1.25.3 или выше
- **Docker** и Docker Compose
- **Git**
- **goose** (для миграций базы данных)

### Установка

```bash
# 1. Клонировать репозиторий
git clone https://github.com/gbroccoli/HeiCRM.git
cd HeiCRM

# 2. Установить goose для миграций
go install github.com/pressly/goose/v3/cmd/goose@latest

# 3. Создать конфигурацию из примера
cp config.example.yaml config.yaml

# 4. Отредактировать config.yaml
# ВАЖНО: Установите jwt.secret_key (минимум 32 символа)
```

### Конфигурация

Пример `config.yaml`:

```yaml
env: dev
service_name: HeiCRM

jwt:
  secret_key: "ваш-секретный-ключ-минимум-32-символа"
  alg: HS256

database:
  host: localhost
  port: 5432
  user: postgres
  password: root
  name: crm
  sslmode: disable

access_token_ttl: "30m"
refresh_token_ttl: "720h"  # 30 дней

nats:
  url: "nats://localhost:4222"
```

### Запуск

```bash
# 1. Запустить инфраструктуру (PostgreSQL + NATS)
docker compose up -d db nats

# 2. Создать базу данных
docker exec -it database psql -U postgres -c "CREATE DATABASE crm;"

# 3. Применить миграции
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" up

# 4. Запустить Auth Service
go run ./services/auth/cmd/auth
```

Сервис будет доступен на `http://localhost:8080`

Подробная инструкция: [QUICKSTART.md](QUICKSTART.md)

## API

### Auth Service

#### POST /auth/login
Вход в систему

**Request:**
```json
{
  "email": "student@university.edu",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "message": "success login",
  "data": {
    "email": "student@university.edu"
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
```

#### POST /auth/register
Регистрация нового пользователя (требует авторизацию администратора)

**Headers:**
```
Authorization: Bearer <admin_access_token>
```

**Request:**
```json
{
  "email": "newstudent@university.edu",
  "password": "securepassword",
  "role": "student"
}
```

**Доступные роли:**
- `student` — студент (жилец)
- `staff` — персонал общежития
- `admin` — администратор

## Структура проекта

```
HeiCRM/
├── config.yaml                 # Конфигурация (не в git)
├── docker-compose.yml          # PostgreSQL + NATS
├── migrations/                 # SQL миграции базы данных
├── pkg/                        # Общие пакеты (config, dbx, jwt, logx, password)
└── services/                   # Микросервисы
    ├── auth/                  # ✅ Сервис авторизации
    ├── user/                  # 📋 Сервис пользователей (планируется)
    ├── building/              # 📋 Сервис зданий (планируется)
    ├── task/                  # 📋 Сервис заявок (планируется)
    └── notification/          # 📋 Сервис уведомлений (планируется)
```

## Технологический стек

- **Язык:** Go 1.25+
- **Web Framework:** Gin
- **База данных:** PostgreSQL
- **Миграции:** goose
- **Брокер сообщений:** NATS
- **Авторизация:** JWT (HS512)
- **Хеширование:** bcrypt

## Типичные сценарии использования

### Для студентов

1. **Регистрация** — создание аккаунта (через администратора)
2. **Подача заявки** — например, "сломался кран в комнате 305"
3. **Отслеживание статуса** — получение уведомлений о продвижении заявки
4. **История заявок** — просмотр всех предыдущих обращений

### Для персонала

1. **Получение заявок** — список новых заявок по категориям
2. **Принятие в работу** — смена статуса на "в обработке"
3. **Обновление статуса** — промежуточные этапы выполнения
4. **Завершение заявки** — отметка о выполнении

### Для администраторов

1. **Управление пользователями** — создание, редактирование, блокировка
2. **Управление зданиями** — структура общежитий, комнаты, жильцы
3. **Аналитика** — статистика по заявкам, времени выполнения
4. **Настройка системы** — категории заявок, шаблоны уведомлений

## Разработка

### Команды для разработки

```bash
# Форматирование кода
go fmt ./...

# Статический анализ
go vet ./...

# Запуск тестов
go test ./...

# Тесты с покрытием
go test ./... -cover

# Создание новой миграции
goose -dir migrations create migration_name sql
```

### Управление миграциями

```bash
# Применить все миграции
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" up

# Откатить последнюю миграцию
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" down

# Статус миграций
goose -dir migrations postgres "postgres://postgres:root@localhost:5432/crm?sslmode=disable" status
```

## Roadmap

Проект находится в активной разработке с целью создать полнофункциональную систему управления общежитиями к концу января 2026 года.

**Текущий статус:** Auth Service (MVP) завершен, работа над User Service и Task Service.

📅 **Детальный план разработки:** [ROADMAP.md](ROADMAP.md)

## Документация

- **[QUICKSTART.md](QUICKSTART.md)** — быстрый старт и запуск проекта
- **[ROADMAP.md](ROADMAP.md)** — детальный план развития
- **[CLAUDE.md](CLAUDE.md)** — инструкции для разработки с Claude Code
- **[AGENTS.md](AGENTS.md)** — автоматизация и агенты

## Вклад в проект

Мы приветствуем ваш вклад в развитие HeiCRM!

1. Изучите [ROADMAP.md](ROADMAP.md) и выберите задачу
2. Создайте issue для обсуждения
3. Форкните репозиторий
4. Создайте feature branch
5. Реализуйте функционал с тестами
6. Создайте Pull Request

## Контакты

- GitHub: [gbroccoli](https://github.com/gbroccoli)

---

**Версия:** 1.5
**Статус:** В разработке
**Последнее обновление:** Февраль 2026

⚠️ **Внимание:** Проект находится на ранней стадии разработки. Не рекомендуется для production использования.
