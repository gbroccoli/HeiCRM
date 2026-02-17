# Database Migrations

Этот проект использует [goose](https://github.com/pressly/goose) для управления миграциями базы данных.

## Установка goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Использование

### Вариант 1: Через Makefile (Linux/Mac)

```bash
# Применить все миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down

# Посмотреть статус миграций
make migrate-status

# Создать новую миграцию
make migrate-create NAME=add_companies_table
```

### Вариант 2: Через скрипт (Linux/Mac/Git Bash)

```bash
# Дать права на выполнение (только первый раз)
chmod +x scripts/migrate.sh

# Применить все миграции
./scripts/migrate.sh up

# Откатить последнюю миграцию
./scripts/migrate.sh down

# Посмотреть статус миграций
./scripts/migrate.sh status

# Создать новую миграцию
./scripts/migrate.sh create add_companies_table
```

### Вариант 3: Прямой вызов goose (Windows/универсально)

```bash
# Применить все миграции
goose -dir migrations postgres 'postgres://postgres:root@localhost:5432/crm?sslmode=disable' up

# Откатить последнюю миграцию
goose -dir migrations postgres 'postgres://postgres:root@localhost:5432/crm?sslmode=disable' down

# Посмотреть статус миграций
goose -dir migrations postgres 'postgres://postgres:root@localhost:5432/crm?sslmode=disable' status

# Создать новую миграцию
goose -dir migrations create add_companies_table sql
```

## Настройка подключения к БД

По умолчанию используются следующие параметры:
- Host: localhost
- Port: 5432
- User: postgres
- Password: root
- Database: crm
- SSLMode: disable

Для изменения параметров используйте переменные окружения:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASS=root
export DB_NAME=crm
export DB_SSLMODE=disable

# Теперь запускайте миграции
./scripts/migrate.sh up
```

## Структура миграций

Каждая миграция состоит из двух частей:

```sql
-- +goose Up
-- Код для применения миграции
CREATE TABLE ...;

-- +goose Down
-- Код для отката миграции
DROP TABLE ...;
```

## Первая миграция

Создает таблицу `users` с полями:
- `id` - уникальный идентификатор (автоинкремент)
- `name` - имя пользователя
- `email` - email (уникальный, с индексом)
- `password` - хеш пароля (bcrypt)
- `role` - роль пользователя (0 = обычный пользователь)
- `tg_send` - флаг отправки в Telegram
- `created_at` - дата создания
- `updated_at` - дата обновления

## Пример создания новой миграции

1. Создайте файл миграции:
```bash
goose -dir migrations create add_companies_table sql
```

2. Отредактируйте созданный файл:
```sql
-- +goose Up
CREATE TABLE companies (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS companies CASCADE;
```

3. Примените миграцию:
```bash
goose -dir migrations postgres 'postgres://...' up
```

## Важно

- **Никогда не редактируйте** уже примененные миграции
- Всегда создавайте новую миграцию для изменений
- Проверяйте миграции на тестовой БД перед применением на продакшене
- Сохраняйте бэкапы перед применением миграций на продакшене
