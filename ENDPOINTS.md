# HeiCRM API Endpoints

Все запросы проходят через API Gateway (`localhost:8000`).
Префикс: `/api/v1`

---

## API Gateway

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/health` | Health check gateway |

**Ответ:**
```json
{ "status": "ok", "service": "api-gateway" }
```

---

## Auth Service (порт 8080)

Gateway prefix: `/api/v1/auth`

### POST `/api/v1/auth/login`

Аутентификация пользователя. Возвращает access token в теле и refresh token в HTTPOnly cookie.

**Middleware:** нет

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

**Response 200:**
```json
{
  "code": 1000,
  "msg": "successfully logged in",
  "token": "eyJhbGciOi..."
}
```

**Cookie:** `refresh=<jwt>; Path=/; HttpOnly; SameSite=Lax|None; Secure`

---

### POST `/api/v1/auth/register`

Регистрация нового пользователя. Пароль генерируется автоматически (24 символа).

**Middleware:** нет

**Request:**
```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "role_id": 0,
  "tg_send": false
}
```

**Response 201:**
```json
{
  "code": 1001,
  "msg": "user created",
  "password": "aB3$dE..."
}
```

---

### POST `/api/v1/auth/refresh`

Ротация токенов. Старый refresh token удаляется из Redis, выдаётся новая пара.

**Middleware:** RefreshTokenMiddleware (валидация refresh cookie + Redis)

**Request:** без тела, refresh token из cookie

**Response 200:**
```json
{
  "code": 1000,
  "token": "eyJhbGciOi..."
}
```

**Cookie:** новый `refresh=<jwt>; Path=/; HttpOnly`

---

### POST `/api/v1/auth/logout`

Выход — удаление refresh token из Redis и очистка cookie.

**Middleware:** AuthMiddleware (Bearer token)

**Headers:** `Authorization: Bearer <access_token>`

**Request:** без тела

**Response 200:**
```json
{
  "code": 1000,
  "msg": "successfully logged out"
}
```

---

### GET `/api/v1/auth/me`

Текущий пользователь (из auth service).

**Middleware:** AuthMiddleware (Bearer token)

**Headers:** `Authorization: Bearer <access_token>`

**Response 200:**
```json
{
  "code": 1000,
  "id": 1,
  "name": "John Doe",
  "email": "user@example.com",
  "avatar": null,
  "role": "admin"
}
```

---

## User Service (порт 8081)

Gateway prefix: `/api/v1/users`

Все эндпоинты требуют `Authorization: Bearer <access_token>`.

### GET `/api/v1/users/me`

Профиль текущего пользователя с расширенными данными.

**Доступ:** любой авторизованный пользователь

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "user@example.com",
    "avatar": null,
    "role_id": 0,
    "role_name": "user",
    "first_name": "John",
    "last_name": "Doe",
    "middle_name": null,
    "phone": "+79001234567",
    "student_id": null,
    "room_id": null,
    "avatar_url": null,
    "date_of_birth": null,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### PUT `/api/v1/users/me`

Обновление своего профиля.

**Доступ:** любой авторизованный пользователь

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "middle_name": "M.",
  "phone": "+79001234567"
}
```

Все поля опциональны. `null` — поле не меняется.

**Response 200:**
```json
{
  "code": 1002,
  "message": "profile updated",
  "data": { ... }
}
```

---

### GET `/api/v1/users/`

Список пользователей с пагинацией.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Query params:**

| Параметр | По умолчанию | Описание |
|----------|-------------|----------|
| `page` | 1 | Номер страницы |
| `page_size` | 20 | Элементов на странице (макс. 100) |

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "items": [
      {
        "id": 1,
        "name": "John Doe",
        "email": "user@example.com",
        "role_id": 0,
        "role_name": "user",
        "first_name": "John",
        "last_name": "Doe",
        ...
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 42,
      "total_pages": 3
    }
  }
}
```

---

### GET `/api/v1/users/:id`

Получение пользователя по ID.

**Доступ:** admin, manager — любой пользователь; обычный пользователь — только свой профиль

**Response 200:**
```json
{
  "code": 1000,
  "data": { ... }
}
```

**Response 403 (чужой профиль для обычного пользователя):**
```json
{
  "code": 4000,
  "message": "insufficient permissions"
}
```

---

### PUT `/api/v1/users/:id`

Обновление пользователя администратором.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "name": "New Name",
  "role_id": 2,
  "first_name": "John",
  "last_name": "Doe",
  "middle_name": "M.",
  "phone": "+79001234567",
  "student_id": "STU-001",
  "room_id": 5,
  "avatar_url": "https://example.com/avatar.jpg",
  "date_of_birth": "2000-01-15"
}
```

Все поля опциональны. `null` — поле не меняется.
`date_of_birth` — формат `YYYY-MM-DD`.

**Response 200:**
```json
{
  "code": 1002,
  "message": "user updated",
  "data": { ... }
}
```

---

### DELETE `/api/v1/users/:id`

Удаление пользователя.

**Доступ:** admin (role_id=1)

Нельзя удалить самого себя.

**Response 200:**
```json
{
  "code": 1003,
  "message": "user deleted"
}
```

**Response 403 (попытка удалить себя):**
```json
{
  "code": 4000,
  "message": "cannot delete yourself"
}
```

**Response 404:**
```json
{
  "code": 5000,
  "message": "user not found"
}
```

---

### GET `/api/v1/users/:id/activity`

Лог активности пользователя.

**Доступ:** admin (role_id=1)

**Query params:**

| Параметр | По умолчанию | Описание |
|----------|-------------|----------|
| `page` | 1 | Номер страницы |
| `page_size` | 50 | Элементов на странице (макс. 100) |

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 5,
        "action": "profile_updated",
        "details": { "field": "phone" },
        "ip_address": "192.168.1.1",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2025-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 10,
      "total_pages": 1
    }
  }
}
```

---

## NATS Events

### Подписки (User Service слушает)

| Subject | Описание | Payload |
|---------|----------|---------|
| `user.registered` | Новый пользователь зарегистрирован | `{"user_id": 1, "email": "...", "name": "..."}` |

### Публикации (User Service отправляет)

| Subject | Описание | Payload |
|---------|----------|---------|
| `user.profile_updated` | Профиль обновлён | `{"user_id": 1, "email": "..."}` |
| `user.deactivated` | Пользователь деактивирован | `{"user_id": 1, "email": "..."}` |

---

## Роли

| role_id | Название | Описание |
|---------|----------|----------|
| 0 | user | Обычный пользователь |
| 1 | admin | Администратор (полный доступ) |
| 2 | manager | Менеджер (чтение пользователей) |

---

## Коды ответов

| Код | Константа | Описание |
|-----|-----------|----------|
| 1000 | OK | Успешная операция |
| 1001 | Created | Ресурс создан |
| 1002 | Updated | Ресурс обновлён |
| 1003 | Deleted | Ресурс удалён |
| 2000 | InvalidData | Невалидные данные |
| 2001 | InvalidFormat | Неверный формат |
| 3000 | AuthRequired | Требуется аутентификация |
| 3001 | InvalidCredentials | Неверные учётные данные |
| 3002 | InvalidToken | Невалидный токен |
| 3003 | ExpiredToken | Токен истёк |
| 4000 | AccessDenied | Доступ запрещён |
| 4001 | InsufficientRights | Недостаточно прав |
| 5000 | NotFound | Не найдено |
| 5001 | AlreadyExists | Уже существует |
| 9000 | InternalError | Внутренняя ошибка |
| 9001 | DatabaseError | Ошибка базы данных |

---

## Аутентификация

Все защищённые эндпоинты требуют заголовок:

```
Authorization: Bearer <access_token>
```

- **Access token** — JWT HS512, время жизни 45 минут
- **Refresh token** — JWT HS512, время жизни 30 дней, хранится в HTTPOnly cookie и Redis
- Ротация токенов через `POST /api/v1/auth/refresh`
