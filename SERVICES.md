# Архитектура сервисов HeiCRM

Детальное описание всех микросервисов системы управления заявками для общежитий.

**Дипломный проект**

---

## Оглавление

1. [Auth Service](#1-auth-service) (✅ Реализован)
2. [User Service](#2-user-service) (📋 Планируется)
3. [Building Service](#3-building-service) (📋 Планируется)
4. [Task Service](#4-task-service) (📋 Планируется)
5. [Notification Service](#5-notification-service) (📋 Планируется)
6. [Взаимодействие сервисов](#взаимодействие-сервисов)

---

## 1. Auth Service

**Статус:** ✅ Реализован (MVP)
**Порт:** 8080
**Префикс API:** `/auth`

### Назначение

Центральный сервис аутентификации и авторизации. Управляет регистрацией пользователей, выдачей JWT токенов и контролем доступа.

### Функциональность

- Регистрация новых пользователей (только администраторами)
- Аутентификация (вход в систему)
- Выдача JWT токенов (access + refresh)
- Обновление access токена через refresh token
- Выход из системы (инвалидация токенов)
- Управление ролями пользователей

### База данных

#### Таблица: `users`

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- bcrypt hash
    role VARCHAR(50) NOT NULL,       -- student, staff, admin
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

### API Endpoints

#### POST `/auth/register`
Регистрация нового пользователя (требует admin права)

**Authorization:** Bearer token (admin)

**Request:**
```json
{
  "email": "student@university.edu",
  "password": "securePassword123",
  "role": "student"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "email": "student@university.edu",
    "role": "student"
  }
}
```

#### POST `/auth/login`
Вход в систему

**Request:**
```json
{
  "email": "student@university.edu",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "email": "student@university.edu",
      "role": "student"
    },
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci..."
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
```

#### POST `/auth/refresh`
Обновление access токена

**Request:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci..."
  }
}
```

#### POST `/auth/logout`
Выход из системы

**Authorization:** Bearer token

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### JWT Token Claims

**Access Token (30 минут):**
```json
{
  "user_id": 1,
  "email": "student@university.edu",
  "role": "student",
  "type": "access",
  "exp": 1234567890,
  "iat": 1234567890
}
```

**Refresh Token (30 дней):**
```json
{
  "user_id": 1,
  "email": "student@university.edu",
  "type": "refresh",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### События NATS

**Публикует:**
- `user.registered` — новый пользователь зарегистрирован
- `user.login` — пользователь вошел в систему
- `user.logout` — пользователь вышел из системы

---

## 2. User Service

**Статус:** 📋 Планируется
**Порт:** 8081
**Префикс API:** `/users`

### Назначение

Управление профилями пользователей, расширенной информацией и привязкой к комнатам общежития.

### Функциональность

- CRUD операции для пользователей
- Расширенные профили (имя, телефон, студенческий билет)
- Загрузка аватаров
- Привязка студентов к комнатам
- История активности пользователя
- Управление ролями и правами доступа

### База данных

#### Таблица: `user_profiles`

```sql
CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,  -- связь с users.id из Auth Service
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    middle_name VARCHAR(100),
    phone VARCHAR(20),
    student_id VARCHAR(50),           -- номер студенческого билета
    room_id INTEGER,                  -- связь с rooms.id
    avatar_url VARCHAR(500),
    date_of_birth DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE SET NULL
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX idx_user_profiles_room_id ON user_profiles(room_id);
CREATE INDEX idx_user_profiles_student_id ON user_profiles(student_id);
```

#### Таблица: `user_activity_log`

```sql
CREATE TABLE user_activity_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    action VARCHAR(100) NOT NULL,     -- login, logout, task_created, profile_updated
    details JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES user_profiles(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_activity_user_id ON user_activity_log(user_id);
CREATE INDEX idx_activity_created_at ON user_activity_log(created_at DESC);
```

### API Endpoints

#### GET `/users`
Список пользователей (admin, staff)

**Authorization:** admin, staff

**Query params:**
- `page` (default: 1)
- `limit` (default: 20)
- `role` (filter: student, staff, admin)
- `room_id` (filter by room)
- `search` (search by name, email)

**Response:**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "email": "student@university.edu",
        "role": "student",
        "profile": {
          "first_name": "Иван",
          "last_name": "Петров",
          "phone": "+79001234567",
          "student_id": "ST-2024-001",
          "room_number": "305",
          "building_name": "Общежитие №1"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "pages": 8
    }
  }
}
```

#### GET `/users/:id`
Информация о пользователе

**Authorization:** admin, staff, или сам пользователь

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "student@university.edu",
    "role": "student",
    "profile": {
      "first_name": "Иван",
      "last_name": "Петров",
      "middle_name": "Сергеевич",
      "phone": "+79001234567",
      "student_id": "ST-2024-001",
      "avatar_url": "https://cdn.example.com/avatars/1.jpg",
      "date_of_birth": "2003-05-15",
      "room": {
        "id": 42,
        "number": "305",
        "floor": 3,
        "building_name": "Общежитие №1"
      }
    },
    "created_at": "2024-09-01T10:00:00Z"
  }
}
```

#### GET `/users/me`
Текущий пользователь

**Authorization:** любой авторизованный пользователь

**Response:** аналогично GET `/users/:id`

#### PUT `/users/me`
Обновление своего профиля

**Authorization:** любой авторизованный пользователь

**Request:**
```json
{
  "first_name": "Иван",
  "last_name": "Петров",
  "phone": "+79001234567"
}
```

#### PUT `/users/:id`
Обновление пользователя (admin)

**Authorization:** admin

**Request:**
```json
{
  "role": "staff",
  "is_active": true,
  "profile": {
    "first_name": "Иван",
    "room_id": 42
  }
}
```

#### POST `/users/:id/avatar`
Загрузка аватара

**Authorization:** пользователь или admin

**Request:** multipart/form-data
- `avatar` (file, max 5MB, jpg/png)

**Response:**
```json
{
  "success": true,
  "data": {
    "avatar_url": "https://cdn.example.com/avatars/1.jpg"
  }
}
```

#### DELETE `/users/:id`
Деактивация пользователя (admin)

**Authorization:** admin

**Response:**
```json
{
  "success": true,
  "message": "User deactivated successfully"
}
```

#### GET `/users/:id/activity`
История активности (admin)

**Authorization:** admin

**Response:**
```json
{
  "success": true,
  "data": {
    "activities": [
      {
        "action": "login",
        "details": {},
        "ip_address": "192.168.1.1",
        "created_at": "2024-11-21T10:00:00Z"
      },
      {
        "action": "task_created",
        "details": {"task_id": 15, "title": "Сломался кран"},
        "created_at": "2024-11-21T10:05:00Z"
      }
    ]
  }
}
```

### События NATS

**Подписывается на:**
- `user.registered` (Auth Service) — создать профиль

**Публикует:**
- `user.profile_updated` — профиль обновлен
- `user.deactivated` — пользователь деактивирован
- `user.room_assigned` — пользователь привязан к комнате

---

## 3. Building Service

**Статус:** 📋 Планируется
**Порт:** 8082
**Префикс API:** `/buildings`, `/rooms`

### Назначение

Управление зданиями общежитий, этажами, комнатами и привязкой студентов к комнатам.

### Функциональность

- Управление зданиями общежитий
- Иерархия: Здание → Этаж → Комната
- Информация о комнатах (вместимость, статус)
- Привязка студентов к комнатам
- История заселений/выселений
- Статусы комнат (свободна, занята, на ремонте)

### База данных

#### Таблица: `buildings`

```sql
CREATE TABLE buildings (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    floors_count INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_buildings_name ON buildings(name);
```

#### Таблица: `floors`

```sql
CREATE TABLE floors (
    id SERIAL PRIMARY KEY,
    building_id INTEGER NOT NULL,
    floor_number INTEGER NOT NULL,
    rooms_count INTEGER DEFAULT 0,

    FOREIGN KEY (building_id) REFERENCES buildings(id) ON DELETE CASCADE,
    UNIQUE(building_id, floor_number)
);

CREATE INDEX idx_floors_building_id ON floors(building_id);
```

#### Таблица: `rooms`

```sql
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    floor_id INTEGER NOT NULL,
    room_number VARCHAR(20) NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 2,    -- вместимость (человек)
    current_occupancy INTEGER DEFAULT 0,    -- текущее заселение
    status VARCHAR(50) DEFAULT 'available', -- available, occupied, maintenance, reserved
    room_type VARCHAR(50),                  -- single, double, triple, suite
    amenities JSONB,                        -- {wifi: true, bathroom: "shared", ...}
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (floor_id) REFERENCES floors(id) ON DELETE CASCADE
);

CREATE INDEX idx_rooms_floor_id ON rooms(floor_id);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_rooms_number ON rooms(room_number);
```

#### Таблица: `room_assignments`

```sql
CREATE TABLE room_assignments (
    id SERIAL PRIMARY KEY,
    room_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    check_in_date DATE NOT NULL,
    check_out_date DATE,
    is_active BOOLEAN DEFAULT true,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE INDEX idx_assignments_room_id ON room_assignments(room_id);
CREATE INDEX idx_assignments_user_id ON room_assignments(user_id);
CREATE INDEX idx_assignments_active ON room_assignments(is_active);
```

### API Endpoints

#### Здания

##### GET `/buildings`
Список зданий

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Общежитие №1",
      "address": "ул. Студенческая, д. 5",
      "floors_count": 9,
      "total_rooms": 180,
      "occupied_rooms": 145
    }
  ]
}
```

##### POST `/buildings`
Создание здания (admin)

**Authorization:** admin

**Request:**
```json
{
  "name": "Общежитие №1",
  "address": "ул. Студенческая, д. 5",
  "floors_count": 9,
  "description": "Главное общежитие университета"
}
```

##### GET `/buildings/:id`
Информация о здании

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Общежитие №1",
    "address": "ул. Студенческая, д. 5",
    "floors_count": 9,
    "description": "Главное общежитие университета",
    "statistics": {
      "total_rooms": 180,
      "occupied_rooms": 145,
      "available_rooms": 25,
      "maintenance_rooms": 10,
      "occupancy_rate": 80.5
    }
  }
}
```

##### GET `/buildings/:id/floors`
Этажи здания

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "floor_number": 1,
      "rooms_count": 20,
      "occupied_count": 18
    }
  ]
}
```

#### Комнаты

##### GET `/rooms`
Список комнат

**Query params:**
- `building_id` (filter)
- `floor_id` (filter)
- `status` (filter: available, occupied, maintenance)
- `room_type` (filter)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 42,
      "room_number": "305",
      "floor_number": 3,
      "building_name": "Общежитие №1",
      "capacity": 2,
      "current_occupancy": 1,
      "status": "available",
      "room_type": "double",
      "amenities": {
        "wifi": true,
        "bathroom": "shared",
        "balcony": false
      }
    }
  ]
}
```

##### GET `/rooms/:id`
Информация о комнате

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 42,
    "room_number": "305",
    "floor": {
      "floor_number": 3,
      "building_name": "Общежитие №1"
    },
    "capacity": 2,
    "current_occupancy": 1,
    "status": "available",
    "room_type": "double",
    "amenities": {
      "wifi": true,
      "bathroom": "shared",
      "balcony": false
    },
    "current_residents": [
      {
        "user_id": 5,
        "name": "Иван Петров",
        "check_in_date": "2024-09-01"
      }
    ]
  }
}
```

##### GET `/rooms/:id/residents`
Жильцы комнаты

**Response:**
```json
{
  "success": true,
  "data": {
    "current": [
      {
        "user_id": 5,
        "name": "Иван Петров",
        "email": "student@university.edu",
        "check_in_date": "2024-09-01"
      }
    ],
    "history": [
      {
        "user_id": 3,
        "name": "Петр Сидоров",
        "check_in_date": "2023-09-01",
        "check_out_date": "2024-06-30"
      }
    ]
  }
}
```

##### POST `/rooms/:id/assign`
Заселить студента (admin, staff)

**Authorization:** admin, staff

**Request:**
```json
{
  "user_id": 5,
  "check_in_date": "2024-09-01",
  "notes": "Переезд из комнаты 210"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Student assigned to room successfully",
  "data": {
    "assignment_id": 123,
    "room_number": "305",
    "user_name": "Иван Петров"
  }
}
```

##### POST `/rooms/:id/unassign`
Выселить студента (admin, staff)

**Authorization:** admin, staff

**Request:**
```json
{
  "user_id": 5,
  "check_out_date": "2024-11-21"
}
```

##### PUT `/rooms/:id`
Обновление комнаты (admin)

**Authorization:** admin

**Request:**
```json
{
  "status": "maintenance",
  "notes": "Ремонт сантехники"
}
```

### События NATS

**Публикует:**
- `room.assigned` — студент заселен в комнату
- `room.unassigned` — студент выселен из комнаты
- `room.status_changed` — статус комнаты изменен

**Подписывается на:**
- `user.deactivated` — выселить пользователя из комнаты

---

## 4. Task Service

**Статус:** 📋 Планируется
**Порт:** 8083
**Префикс API:** `/tasks`

### Назначение

Основной функционал системы — управление заявками от студентов на различные услуги.

### Функциональность

- Создание заявок студентами
- Категории: Ремонт, Уборка, Доставка, IT-поддержка, Прочее
- Статусы: Новая, В обработке, Выполнена, Отклонена
- Приоритеты: Низкий, Средний, Высокий, Критический
- Назначение исполнителя
- Комментарии и обсуждение
- Прикрепление фотографий
- История изменений
- Дедлайны

### База данных

#### Таблица: `task_categories`

```sql
CREATE TABLE task_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),                   -- emoji или название иконки
    default_assignee_role VARCHAR(50),  -- staff
    default_priority VARCHAR(50) DEFAULT 'medium',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Предустановленные категории
INSERT INTO task_categories (name, description, icon, default_assignee_role) VALUES
('Ремонт', 'Ремонт мебели, сантехники, электрики', '🔧', 'staff'),
('Уборка', 'Уборка комнаты или общих помещений', '🧹', 'staff'),
('Доставка', 'Доставка посылок, документов', '📦', 'staff'),
('IT-поддержка', 'Проблемы с интернетом, оборудованием', '💻', 'staff'),
('Прочее', 'Другие запросы', '📝', 'staff');
```

#### Таблица: `tasks`

```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'new',        -- new, in_progress, completed, rejected
    priority VARCHAR(50) DEFAULT 'medium',   -- low, medium, high, critical
    creator_id INTEGER NOT NULL,             -- студент
    assignee_id INTEGER,                     -- персонал
    room_id INTEGER,                         -- комната, где проблема
    deadline TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (category_id) REFERENCES task_categories(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE SET NULL
);

CREATE INDEX idx_tasks_creator_id ON tasks(creator_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_category_id ON tasks(category_id);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_room_id ON tasks(room_id);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);
```

#### Таблица: `task_comments`

```sql
CREATE TABLE task_comments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    comment TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT false,  -- внутренний комментарий для персонала
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_comments_task_id ON task_comments(task_id);
```

#### Таблица: `task_attachments`

```sql
CREATE TABLE task_attachments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50),              -- image/jpeg, image/png
    file_size INTEGER,                  -- bytes
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_attachments_task_id ON task_attachments(task_id);
```

#### Таблица: `task_status_history`

```sql
CREATE TABLE task_status_history (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    changed_by INTEGER NOT NULL,
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_status_history_task_id ON task_status_history(task_id);
```

### API Endpoints

#### Для студентов

##### POST `/tasks`
Создать заявку

**Authorization:** student

**Request:**
```json
{
  "title": "Сломался кран в ванной",
  "description": "Из крана капает вода, не останавливается",
  "category_id": 1,
  "priority": "high",
  "room_id": 42
}
```

**Response:**
```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {
    "id": 156,
    "title": "Сломался кран в ванной",
    "status": "new",
    "priority": "high",
    "created_at": "2024-11-21T10:00:00Z"
  }
}
```

##### GET `/tasks`
Список своих заявок

**Authorization:** student

**Query params:**
- `status` (filter)
- `category_id` (filter)
- `page`, `limit`

**Response:**
```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "id": 156,
        "title": "Сломался кран в ванной",
        "description": "Из крана капает вода...",
        "category": {
          "id": 1,
          "name": "Ремонт",
          "icon": "🔧"
        },
        "status": "in_progress",
        "priority": "high",
        "room_number": "305",
        "assignee": {
          "name": "Сергей Иванов"
        },
        "created_at": "2024-11-21T10:00:00Z",
        "updated_at": "2024-11-21T11:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 5
    }
  }
}
```

##### GET `/tasks/:id`
Детали заявки

**Authorization:** creator или staff/admin

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 156,
    "title": "Сломался кран в ванной",
    "description": "Из крана капает вода, не останавливается",
    "category": {
      "id": 1,
      "name": "Ремонт",
      "icon": "🔧"
    },
    "status": "in_progress",
    "priority": "high",
    "creator": {
      "id": 5,
      "name": "Иван Петров",
      "email": "student@university.edu"
    },
    "assignee": {
      "id": 10,
      "name": "Сергей Иванов"
    },
    "room": {
      "id": 42,
      "number": "305",
      "building": "Общежитие №1"
    },
    "deadline": "2024-11-22T18:00:00Z",
    "comments": [
      {
        "id": 1,
        "user_name": "Сергей Иванов",
        "comment": "Принял в работу, буду завтра",
        "created_at": "2024-11-21T11:00:00Z"
      }
    ],
    "attachments": [
      {
        "id": 1,
        "file_url": "https://cdn.example.com/tasks/156/photo1.jpg",
        "file_name": "кран.jpg",
        "created_at": "2024-11-21T10:01:00Z"
      }
    ],
    "history": [
      {
        "old_status": "new",
        "new_status": "in_progress",
        "changed_by": "Сергей Иванов",
        "created_at": "2024-11-21T11:00:00Z"
      }
    ],
    "created_at": "2024-11-21T10:00:00Z",
    "updated_at": "2024-11-21T11:00:00Z"
  }
}
```

##### POST `/tasks/:id/comments`
Добавить комментарий

**Authorization:** creator или staff/admin

**Request:**
```json
{
  "comment": "Фото проблемы прикрепил"
}
```

##### POST `/tasks/:id/attachments`
Прикрепить фото

**Authorization:** creator или staff/admin

**Request:** multipart/form-data
- `file` (max 10MB, jpg/png/pdf)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "file_url": "https://cdn.example.com/tasks/156/photo1.jpg",
    "file_name": "кран.jpg"
  }
}
```

#### Для персонала

##### GET `/tasks`
Список всех заявок

**Authorization:** staff, admin

**Query params:**
- `status` (filter)
- `category_id` (filter)
- `priority` (filter)
- `assignee_id` (filter)
- `room_id` (filter)
- `building_id` (filter)
- `sort` (created_at, priority, deadline)
- `page`, `limit`

**Response:** аналогично списку для студентов

##### GET `/tasks/assigned-to-me`
Мои назначенные заявки

**Authorization:** staff

**Response:**
```json
{
  "success": true,
  "data": {
    "tasks": [...],
    "statistics": {
      "total": 12,
      "new": 3,
      "in_progress": 7,
      "completed_today": 2
    }
  }
}
```

##### PUT `/tasks/:id/assign`
Назначить себе заявку

**Authorization:** staff

**Response:**
```json
{
  "success": true,
  "message": "Task assigned successfully"
}
```

##### PUT `/tasks/:id/status`
Изменить статус

**Authorization:** staff, admin

**Request:**
```json
{
  "status": "in_progress",
  "comment": "Принял в работу"
}
```

##### GET `/tasks/statistics`
Статистика

**Authorization:** staff, admin

**Query params:**
- `period` (today, week, month)
- `building_id` (filter)

**Response:**
```json
{
  "success": true,
  "data": {
    "total_tasks": 450,
    "by_status": {
      "new": 15,
      "in_progress": 32,
      "completed": 380,
      "rejected": 23
    },
    "by_category": {
      "Ремонт": 180,
      "Уборка": 120,
      "IT-поддержка": 90,
      "Доставка": 40,
      "Прочее": 20
    },
    "by_priority": {
      "low": 200,
      "medium": 180,
      "high": 60,
      "critical": 10
    },
    "avg_completion_time": "8.5 hours",
    "overdue_tasks": 5
  }
}
```

#### Для администраторов

##### GET `/tasks/categories`
Список категорий

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Ремонт",
      "description": "Ремонт мебели, сантехники, электрики",
      "icon": "🔧",
      "total_tasks": 180
    }
  ]
}
```

##### POST `/tasks/categories`
Создать категорию

**Authorization:** admin

**Request:**
```json
{
  "name": "Безопасность",
  "description": "Вопросы безопасности",
  "icon": "🔒",
  "default_assignee_role": "staff"
}
```

##### PUT `/tasks/:id/assignee`
Назначить исполнителя

**Authorization:** admin

**Request:**
```json
{
  "assignee_id": 10
}
```

##### DELETE `/tasks/:id`
Удалить заявку

**Authorization:** admin

### События NATS

**Публикует:**
- `task.created` — новая заявка создана
- `task.assigned` — заявка назначена исполнителю
- `task.status_changed` — статус изменен
- `task.completed` — заявка выполнена
- `task.comment_added` — добавлен комментарий
- `task.deadline_approaching` — приближается дедлайн (за 2 часа)
- `task.overdue` — заявка просрочена

**Подписывается на:**
- `user.deactivated` — деактивировать заявки пользователя
- `room.status_changed` (maintenance) — информация для контекста

---

## 5. Notification Service

**Статус:** 📋 Планируется
**Порт:** 8084
**Префикс API:** `/notifications`

### Назначение

Централизованная система уведомлений для информирования пользователей о событиях в системе.

### Функциональность

- Email уведомления (SMTP)
- In-app уведомления
- Push уведомления (WebSocket)
- Настройки уведомлений для пользователей
- Шаблоны сообщений
- Группировка уведомлений (дайджесты)
- Очереди через NATS

### База данных

#### Таблица: `notifications`

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(100) NOT NULL,         -- task_created, task_assigned, task_completed, etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    link VARCHAR(500),                  -- ссылка на связанный ресурс
    priority VARCHAR(50) DEFAULT 'normal', -- low, normal, high
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    metadata JSONB,                     -- дополнительные данные
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
```

#### Таблица: `notification_templates`

```sql
CREATE TABLE notification_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,          -- email, in_app, push
    subject VARCHAR(255),               -- для email
    body_template TEXT NOT NULL,        -- шаблон с переменными {{variable}}
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Примеры шаблонов
INSERT INTO notification_templates (name, type, subject, body_template) VALUES
('task_created_email', 'email', 'Новая заявка #{{task_id}}',
 'Здравствуйте!\n\nСоздана новая заявка:\n\n{{task_title}}\n\nКатегория: {{category}}\nПриоритет: {{priority}}\n\nПодробнее: {{link}}'),
('task_assigned_in_app', 'in_app', 'Заявка назначена вам',
 'Заявка "{{task_title}}" назначена вам для выполнения.');
```

#### Таблица: `user_notification_settings`

```sql
CREATE TABLE user_notification_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    email_enabled BOOLEAN DEFAULT true,
    push_enabled BOOLEAN DEFAULT true,
    email_frequency VARCHAR(50) DEFAULT 'instant', -- instant, hourly, daily
    categories JSONB DEFAULT '{}',                  -- настройки по категориям
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Пример categories:
-- {
--   "task_created": {"email": true, "push": true},
--   "task_completed": {"email": true, "push": false}
-- }
```

### API Endpoints

##### GET `/notifications`
Список уведомлений

**Authorization:** любой пользователь

**Query params:**
- `is_read` (filter: true/false)
- `type` (filter)
- `page`, `limit`

**Response:**
```json
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": 1,
        "type": "task_assigned",
        "title": "Заявка назначена вам",
        "message": "Заявка \"Сломался кран в ванной\" назначена вам для выполнения",
        "link": "/tasks/156",
        "priority": "high",
        "is_read": false,
        "created_at": "2024-11-21T11:00:00Z"
      }
    ],
    "unread_count": 5,
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 45
    }
  }
}
```

##### GET `/notifications/unread`
Непрочитанные уведомления

**Authorization:** любой пользователь

**Response:**
```json
{
  "success": true,
  "data": {
    "count": 5,
    "notifications": [...]
  }
}
```

##### PUT `/notifications/:id/read`
Отметить как прочитанное

**Authorization:** владелец уведомления

**Response:**
```json
{
  "success": true,
  "message": "Notification marked as read"
}
```

##### PUT `/notifications/read-all`
Отметить все как прочитанные

**Authorization:** любой пользователь

**Response:**
```json
{
  "success": true,
  "message": "All notifications marked as read"
}
```

##### GET `/notifications/settings`
Настройки уведомлений

**Authorization:** любой пользователь

**Response:**
```json
{
  "success": true,
  "data": {
    "email_enabled": true,
    "push_enabled": true,
    "email_frequency": "instant",
    "categories": {
      "task_created": {"email": true, "push": true},
      "task_assigned": {"email": true, "push": true},
      "task_completed": {"email": true, "push": false},
      "task_comment": {"email": false, "push": true}
    }
  }
}
```

##### PUT `/notifications/settings`
Обновление настроек

**Authorization:** любой пользователь

**Request:**
```json
{
  "email_enabled": true,
  "push_enabled": false,
  "email_frequency": "hourly",
  "categories": {
    "task_created": {"email": true, "push": false}
  }
}
```

### WebSocket

#### Подключение

```javascript
const ws = new WebSocket('ws://localhost:8084/ws?token=<access_token>');

ws.onmessage = (event) => {
  const notification = JSON.parse(event.data);
  console.log('New notification:', notification);
};
```

#### Формат сообщений

```json
{
  "type": "task_assigned",
  "notification": {
    "id": 123,
    "title": "Заявка назначена вам",
    "message": "...",
    "link": "/tasks/156",
    "priority": "high",
    "created_at": "2024-11-21T11:00:00Z"
  }
}
```

### Типы уведомлений

| Тип события | Email | In-App | Push | Кому |
|------------|-------|--------|------|------|
| `task.created` | ✅ | ✅ | ✅ | Staff (доступные для назначения) |
| `task.assigned` | ✅ | ✅ | ✅ | Assignee (исполнитель) |
| `task.status_changed` | ✅ | ✅ | ❌ | Creator (создатель заявки) |
| `task.completed` | ✅ | ✅ | ✅ | Creator |
| `task.comment_added` | ✅ | ✅ | ❌ | Creator + Assignee |
| `task.deadline_approaching` | ✅ | ✅ | ✅ | Assignee (за 2 часа до дедлайна) |
| `task.overdue` | ✅ | ✅ | ✅ | Assignee + Admin |
| `user.room_assigned` | ✅ | ✅ | ❌ | User (заселение в комнату) |

### События NATS

**Подписывается на:**
- `task.*` — все события задач
- `user.*` — события пользователей
- `room.*` — события комнат

**Публикует:**
- `notification.sent` — уведомление отправлено
- `notification.failed` — ошибка отправки

### Email шаблоны

Примеры email уведомлений:

**Новая заявка (для персонала):**
```
Тема: 🔔 Новая заявка #156: Сломался кран в ванной

Здравствуйте!

Создана новая заявка, требующая внимания:

Заявка #156
Название: Сломался кран в ванной
Категория: Ремонт 🔧
Приоритет: Высокий
Комната: 305 (Общежитие №1)
Создатель: Иван Петров (student@university.edu)

Описание:
Из крана капает вода, не останавливается

Подробнее и принять в работу:
https://heicrm.example.com/tasks/156

---
HeiCRM - Система управления общежитием
```

**Заявка выполнена (для студента):**
```
Тема: ✅ Заявка #156 выполнена

Здравствуйте, Иван!

Ваша заявка выполнена:

Заявка #156: Сломался кран в ванной
Исполнитель: Сергей Иванов
Завершена: 21.11.2024 в 15:30

Комментарий исполнителя:
Кран заменен. Проблема устранена.

Если проблема не решена, создайте новую заявку.

Подробнее: https://heicrm.example.com/tasks/156

---
HeiCRM - Система управления общежитием
```

---

## Взаимодействие сервисов

### Схема коммуникации через NATS

```
                    ┌─────────────┐
                    │    NATS     │
                    │   Message   │
                    │   Broker    │
                    └─────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐      ┌──────▼──────┐    ┌─────▼─────┐
   │  Auth   │      │    User     │    │ Building  │
   │ Service │      │   Service   │    │  Service  │
   └────┬────┘      └──────┬──────┘    └─────┬─────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
       ┌──────▼──────┐         ┌────────▼────────┐
       │    Task     │         │  Notification   │
       │   Service   │         │     Service     │
       └─────────────┘         └─────────────────┘
```

### Примеры потоков данных

#### Создание заявки студентом

1. **Студент** создает заявку через Task Service
2. **Task Service** публикует событие `task.created` в NATS
3. **Notification Service** получает событие и:
   - Создает in-app уведомление для доступного персонала
   - Отправляет email уведомление персоналу
   - Отправляет push через WebSocket

#### Заселение студента в комнату

1. **Admin** назначает студента в комнату через Building Service
2. **Building Service** публикует `room.assigned` в NATS
3. **User Service** обновляет профиль пользователя (room_id)
4. **Notification Service** отправляет уведомление студенту

#### Регистрация нового пользователя

1. **Admin** регистрирует пользователя через Auth Service
2. **Auth Service** публикует `user.registered` в NATS
3. **User Service** создает профиль пользователя
4. **Notification Service** отправляет welcome email

### Общие принципы

- **Асинхронная коммуникация** через NATS для несрочных операций
- **Синхронные HTTP вызовы** только для критичных операций
- **Event Sourcing** — все изменения логируются как события
- **Eventual Consistency** — данные синхронизируются асинхронно

---

## Общие технические требования

### Авторизация

Все сервисы используют JWT токены из Auth Service:
- Валидация токена в middleware каждого сервиса
- Проверка роли для защищенных endpoints
- Общая библиотека JWT в `pkg/jwt`

### Обработка ошибок

Единый формат ошибок для всех сервисов:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  }
}
```

### Логирование

- Структурированные логи (JSON)
- Уровни: DEBUG, INFO, WARN, ERROR
- Correlation ID для трассировки запросов между сервисами
- Логи хранятся в `logs/` директории

### Мониторинг

- Health check endpoint: `GET /health`
- Metrics endpoint: `GET /metrics` (Prometheus format)
- Ready check: `GET /ready`

---

**Версия:** 1.0
**Последнее обновление:** 21 ноября 2025
