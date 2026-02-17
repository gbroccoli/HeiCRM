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

## Housing Service (порт 8082)

Gateway prefix: `/api/v1/housing`

Все эндпоинты требуют `Authorization: Bearer <access_token>`.

### Buildings

#### GET `/api/v1/housing/`

Список зданий с пагинацией.

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
        "address": "ул. Ленина, 10",
        "floors": 5,
        "description": "Общежитие №1",
        "room_count": 40,
        "resident_count": 35,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

#### POST `/api/v1/housing/`

Создание нового здания.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "address": "ул. Ленина, 10",
  "floors": 5,
  "description": "Общежитие №1"
}
```

`description` — опционально.

**Response 201:**
```json
{
  "code": 1001,
  "data": {
    "id": 1,
    "address": "ул. Ленина, 10",
    "floors": 5,
    "description": "Общежитие №1",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### GET `/api/v1/housing/:id`

Получение здания по ID с количеством комнат и жильцов.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "id": 1,
    "address": "ул. Ленина, 10",
    "floors": 5,
    "description": "Общежитие №1",
    "room_count": 40,
    "resident_count": 35,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### PUT `/api/v1/housing/:id`

Обновление здания.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "address": "ул. Ленина, 10А",
  "floors": 6,
  "description": "Обновлённое описание"
}
```

Все поля опциональны.

**Response 200:**
```json
{
  "code": 1002,
  "data": { ... }
}
```

---

#### DELETE `/api/v1/housing/:id`

Удаление здания. Каскадно удаляет все комнаты и жильцов.

**Доступ:** admin (role_id=1)

**Response 200:**
```json
{
  "code": 1003,
  "message": "building deleted"
}
```

---

### Rooms

#### GET `/api/v1/housing/:id/rooms`

Список комнат в здании с пагинацией и фильтрацией.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Query params:**

| Параметр | По умолчанию | Описание |
|----------|-------------|----------|
| `page` | 1 | Номер страницы |
| `page_size` | 20 | Элементов на странице (макс. 100) |
| `status` | — | Фильтр по статусу: `free`, `occupied` |
| `floor` | — | Фильтр по этажу |

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "items": [
      {
        "id": 1,
        "building_id": 1,
        "room_number": "101",
        "floor": 1,
        "capacity": 2,
        "room_type": "double",
        "status": "occupied",
        "occupancy": 1,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 40,
      "total_pages": 2
    }
  }
}
```

---

#### POST `/api/v1/housing/:id/rooms`

Создание комнаты в здании.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "room_number": "101",
  "floor": 1,
  "capacity": 2,
  "room_type": "double"
}
```

`room_type`: `single`, `double`, `block`.
Валидация: `floor` не может превышать количество этажей здания.

**Response 201:**
```json
{
  "code": 1001,
  "data": {
    "id": 1,
    "building_id": 1,
    "room_number": "101",
    "floor": 1,
    "capacity": 2,
    "room_type": "double",
    "status": "free",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### GET `/api/v1/housing/:id/rooms/:roomId`

Получение комнаты с активными жильцами.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "id": 1,
    "building_id": 1,
    "room_number": "101",
    "floor": 1,
    "capacity": 2,
    "room_type": "double",
    "status": "occupied",
    "residents": [
      {
        "id": 1,
        "full_name": "Иванов Иван Иванович",
        "birth_date": "2000-05-15",
        "email": "ivanov@example.com",
        "phone": "+79001234567",
        "move_in_date": "2025-09-01",
        "move_out_date": null
      }
    ],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### PUT `/api/v1/housing/:id/rooms/:roomId`

Обновление комнаты.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "room_number": "101A",
  "floor": 1,
  "capacity": 3,
  "room_type": "block",
  "status": "free"
}
```

Все поля опциональны.

**Response 200:**
```json
{
  "code": 1002,
  "data": { ... }
}
```

---

#### DELETE `/api/v1/housing/:id/rooms/:roomId`

Удаление комнаты. Каскадно удаляет жильцов.

**Доступ:** admin (role_id=1)

**Response 200:**
```json
{
  "code": 1003,
  "message": "room deleted"
}
```

---

### Residents

#### POST `/api/v1/housing/:id/rooms/:roomId/residents`

Заселение жильца в комнату.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "full_name": "Иванов Иван Иванович",
  "birth_date": "2000-05-15",
  "email": "ivanov@example.com",
  "phone": "+79001234567",
  "move_in_date": "2025-09-01"
}
```

`email`, `phone` — опциональны.
Валидация: количество активных жильцов не может превышать `capacity` комнаты.

**Response 201:**
```json
{
  "code": 1001,
  "data": {
    "id": 1,
    "room_id": 1,
    "full_name": "Иванов Иван Иванович",
    "birth_date": "2000-05-15",
    "email": "ivanov@example.com",
    "phone": "+79001234567",
    "move_in_date": "2025-09-01",
    "move_out_date": null,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### GET `/api/v1/housing/:id/rooms/:roomId/residents`

Список жильцов в комнате.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Query params:**

| Параметр | По умолчанию | Описание |
|----------|-------------|----------|
| `page` | 1 | Номер страницы |
| `page_size` | 20 | Элементов на странице (макс. 100) |
| `include_moved_out` | false | Включить выселенных жильцов |

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "items": [ ... ],
    "pagination": { ... }
  }
}
```

---

#### GET `/api/v1/housing/:id/rooms/:roomId/residents/:residentId`

Получение жильца по ID.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Response 200:**
```json
{
  "code": 1000,
  "data": { ... }
}
```

---

#### PUT `/api/v1/housing/:id/rooms/:roomId/residents`

Обновление данных жильца. ID жильца передаётся в теле запроса.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "resident_id": 1,
  "full_name": "Иванов Иван Петрович",
  "birth_date": "2000-05-15",
  "email": "ivanov_new@example.com",
  "phone": "+79009876543",
  "move_out_date": "2026-06-30"
}
```

`resident_id` — обязательное. Остальные поля опциональны.

**Response 200:**
```json
{
  "code": 1002,
  "data": { ... }
}
```

---

#### DELETE `/api/v1/housing/:id/rooms/:roomId/residents`

Выселение жильца (soft-delete — устанавливает `move_out_date` на сегодня). ID жильца передаётся в теле запроса.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "resident_id": 1
}
```

**Response 200:**
```json
{
  "code": 1003,
  "data": {
    "id": 1,
    "move_out_date": "2025-06-15",
    ...
  }
}
```

---

#### POST `/api/v1/housing/:id/rooms/:roomId/residents/transfer`

Перевод жильца в другую комнату. Операция транзакционная: выселяет из текущей комнаты и заселяет в новую. ID жильца передаётся в теле запроса.

**Доступ:** admin (role_id=1)

**Request:**
```json
{
  "resident_id": 1,
  "new_building_id": 2,
  "new_room_id": 5
}
```

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "id": 15,
    "room_id": 5,
    "full_name": "Иванов Иван Иванович",
    "move_in_date": "2025-06-15",
    "move_out_date": null,
    ...
  }
}
```

---

## Tasks Service (порт 8083)

Gateway prefix: `/api/v1/tasks`

Все эндпоинты требуют `Authorization: Bearer <access_token>`.

### Tasks CRUD

#### GET `/api/v1/tasks/`

Список задач с пагинацией и фильтрацией.

**Доступ:** все авторизованные пользователи (обычные пользователи видят только свои задачи)

**Query params:**

| Параметр | По умолчанию | Описание |
|----------|-------------|----------|
| `page` | 1 | Номер страницы |
| `page_size` | 20 | Элементов на странице (макс. 100) |
| `status` | — | Фильтр: `new`, `assigned`, `in_progress`, `completed`, `closed` |
| `priority` | — | Фильтр: `low`, `medium`, `high`, `critical` |
| `assignee` | — | Фильтр: `me`, `unassigned`, или ID пользователя |

Сортировка по приоритету (critical > high > medium > low).

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "items": [
      {
        "id": 1,
        "author_id": 1,
        "author_name": "Admin",
        "assignee_id": 2,
        "assignee_name": "Manager",
        "room_id": 5,
        "room_number": "101",
        "building_id": 1,
        "task_type": "Ремонт",
        "description": "Сломан кран в ванной",
        "priority": "high",
        "status": "assigned",
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 15,
      "total_pages": 1
    }
  }
}
```

---

#### POST `/api/v1/tasks/`

Создание новой задачи. Автоматически создаёт запись в истории.

**Доступ:** все авторизованные пользователи

**Request:**
```json
{
  "room_id": 5,
  "task_type": "Ремонт",
  "description": "Сломан кран в ванной",
  "priority": "high"
}
```

`task_type`: произвольная строка (напр. `Ремонт`, `Уборка`, `IT-поддержка`).
`priority`: `low`, `medium`, `high`, `critical`.

**Response 201:**
```json
{
  "code": 1001,
  "data": {
    "id": 1,
    "author_id": 1,
    "assignee_id": null,
    "room_id": 5,
    "task_type": "Ремонт",
    "description": "Сломан кран в ванной",
    "priority": "high",
    "status": "new",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### GET `/api/v1/tasks/:id`

Получение задачи по ID.

**Доступ:** автор или исполнитель (обычные пользователи); все задачи (admin, manager)

**Response 200:**
```json
{
  "code": 1000,
  "data": {
    "id": 1,
    "author_id": 1,
    "author_name": "Admin",
    "assignee_id": 2,
    "assignee_name": "Manager",
    "room_id": 5,
    "room_number": "101",
    "building_id": 1,
    "task_type": "Ремонт",
    "description": "Сломан кран в ванной",
    "priority": "high",
    "status": "assigned",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### PUT `/api/v1/tasks/:id`

Обновление задачи (тип, описание, приоритет).

**Доступ:** автор задачи (обычные пользователи); любая задача (admin, manager)

Нельзя обновлять задачу со статусом `completed` или `closed`.

**Request:**
```json
{
  "task_type": "IT-поддержка",
  "description": "Обновлённое описание",
  "priority": "critical"
}
```

Все поля опциональны.

**Response 200:**
```json
{
  "code": 1002,
  "data": { ... }
}
```

---

#### DELETE `/api/v1/tasks/:id`

Удаление задачи. Каскадно удаляет комментарии, историю и вложения.

**Доступ:** admin (role_id=1)

**Response 200:**
```json
{
  "code": 1003,
  "message": "task deleted"
}
```

---

### Task Status & Assignment

#### PUT `/api/v1/tasks/:id/status`

Изменение статуса задачи. Валидирует допустимые переходы.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Request:**
```json
{
  "status": "in_progress",
  "comment": "Начал работу"
}
```

`comment` — опционально, сохраняется в истории.

Допустимые статусы: `new`, `assigned`, `in_progress`, `completed`, `closed`.

Публикует NATS-событие `task.status_changed`.

**Response 200:**
```json
{
  "code": 1002,
  "data": { ... }
}
```

---

#### PUT `/api/v1/tasks/:id/assign`

Назначение исполнителя задачи.

**Доступ:** admin (role_id=1), manager (role_id=2)

**Request:**
```json
{
  "assignee_id": 3
}
```

Автоматически переводит статус с `new` на `assigned`.
Публикует NATS-событие `task.assigned`.

**Response 200:**
```json
{
  "code": 1002,
  "data": { ... }
}
```

---

#### POST `/api/v1/tasks/:id/take`

Взять задачу на себя (назначить себя исполнителем).

**Доступ:** admin (role_id=1), manager (role_id=2)

**Request:** без тела

Публикует NATS-событие `task.assigned`.

**Response 200:**
```json
{
  "code": 1002,
  "data": { ... }
}
```

---

### Comments

#### GET `/api/v1/tasks/:id/comments`

Список комментариев к задаче.

**Доступ:** все авторизованные пользователи

**Response 200:**
```json
{
  "code": 1000,
  "data": [
    {
      "id": 1,
      "task_id": 1,
      "author_id": 2,
      "author_name": "Manager",
      "comment_text": "Займусь сегодня",
      "created_at": "2025-01-01T12:00:00Z"
    }
  ]
}
```

---

#### POST `/api/v1/tasks/:id/comments`

Добавление комментария к задаче.

**Доступ:** все авторизованные пользователи

**Request:**
```json
{
  "comment_text": "Займусь сегодня"
}
```

**Response 201:**
```json
{
  "code": 1001,
  "data": {
    "id": 1,
    "task_id": 1,
    "author_id": 2,
    "comment_text": "Займусь сегодня",
    "created_at": "2025-01-01T12:00:00Z"
  }
}
```

---

### History

#### GET `/api/v1/tasks/:id/history`

История изменения статуса задачи.

**Доступ:** все авторизованные пользователи

**Response 200:**
```json
{
  "code": 1000,
  "data": [
    {
      "id": 1,
      "task_id": 1,
      "previous_status": "new",
      "new_status": "assigned",
      "changed_by": 1,
      "changed_by_name": "Admin",
      "changed_at": "2025-01-01T10:00:00Z",
      "comment": "Назначил на менеджера"
    }
  ]
}
```

---

### Attachments

#### POST `/api/v1/tasks/:id/attachments`

Загрузка файла-вложения к задаче.

**Доступ:** все авторизованные пользователи (с доступом к задаче)

**Request:** `multipart/form-data`, поле `file` (макс. 10 МБ)

**Response 201:**
```json
{
  "code": 1001,
  "data": {
    "id": 1,
    "task_id": 1,
    "file_name": "photo.jpg",
    "file_path": "./uploads/tasks/1/uuid-photo.jpg",
    "file_size": 245760,
    "uploaded_by": 3,
    "uploaded_at": "2025-01-01T12:00:00Z"
  }
}
```

---

#### GET `/api/v1/tasks/:id/attachments`

Список вложений задачи.

**Доступ:** все авторизованные пользователи (с доступом к задаче)

**Response 200:**
```json
{
  "code": 1000,
  "data": [
    {
      "id": 1,
      "task_id": 1,
      "file_name": "photo.jpg",
      "file_size": 245760,
      "uploaded_by": 3,
      "uploader_name": "User",
      "uploaded_at": "2025-01-01T12:00:00Z"
    }
  ]
}
```

---

#### GET `/api/v1/tasks/:id/attachments/:attachmentId/download`

Скачивание файла-вложения.

**Доступ:** все авторизованные пользователи (с доступом к задаче)

**Response:** бинарный файл с заголовком `Content-Disposition: attachment; filename="photo.jpg"`

---

#### DELETE `/api/v1/tasks/:id/attachments/:attachmentId`

Удаление вложения (файл удаляется с диска и из БД).

**Доступ:** загрузивший файл пользователь или admin (role_id=1)

**Response 200:**
```json
{
  "code": 1003,
  "message": "attachment deleted"
}
```

---

## Notification Service

Сервис уведомлений. Не имеет HTTP-эндпоинтов. Слушает NATS-события и отправляет email-уведомления.

### Обрабатываемые события

| Subject | Действие |
|---------|----------|
| `user.registered` | Отправка приветственного email с логином и паролем |
| `task.assigned` | Уведомление исполнителю о назначении задачи |
| `task.status_changed` | Уведомление автору задачи об изменении статуса |

---

## NATS Events

### User Events

| Subject | Издатель | Подписчики | Payload |
|---------|----------|------------|---------|
| `user.registered` | Auth Service | User Service, Notification Service | `{"user_id": 1, "email": "...", "name": "...", "password": "..."}` |
| `user.profile_updated` | User Service | — | `{"user_id": 1, "email": "..."}` |
| `user.deactivated` | User Service | — | `{"user_id": 1, "email": "..."}` |

### Task Events

| Subject | Издатель | Подписчики | Payload |
|---------|----------|------------|---------|
| `task.assigned` | Tasks Service | Notification Service | `{"task_id": 1, "assignee_id": 2, "author_id": 1, "task_type": "...", "description": "...", "priority": "..."}` |
| `task.status_changed` | Tasks Service | Notification Service | `{"task_id": 1, "author_id": 1, "assignee_id": 2, "previous_status": "new", "new_status": "assigned", "changed_by": 1, "task_type": "...", "description": "..."}` |

---

## Роли

| role_id | Название | Описание |
|---------|----------|----------|
| 0 | user | Обычный пользователь |
| 1 | admin | Администратор (полный доступ) |
| 2 | manager | Менеджер (чтение пользователей, управление задачами) |

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
