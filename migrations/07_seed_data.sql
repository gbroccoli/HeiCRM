-- +goose Up
-- +goose NO TRANSACTION
--
-- HeiCRM Seed Data
-- Admin credentials:
--   email:    admin@heicrm.ru
--   password: Admin123!@#

-- ============================================================
-- Roles
-- ============================================================
INSERT INTO roles (id, name, description) VALUES
    (0, 'user',    'Обычный пользователь'),
    (1, 'admin',   'Администратор (полный доступ)'),
    (2, 'manager', 'Менеджер (чтение пользователей, управление задачами)')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Users (15 пользователей: 1 admin, 2 manager, 12 user)
-- Пароль для всех: Admin123!@#
-- ============================================================
INSERT INTO users (id, name, email, password, avatar, role_id, tg_send) VALUES
    (1,  'Admin',              'admin@heicrm.ru',        '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 1, false),
    (2,  'Петров Алексей',     'petrov@heicrm.ru',       '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 2, false),
    (3,  'Сидорова Мария',     'sidorova@heicrm.ru',     '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 2, false),
    (4,  'Иванов Дмитрий',     'ivanov@heicrm.ru',       '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (5,  'Козлова Анна',       'kozlova@heicrm.ru',      '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (6,  'Новиков Сергей',     'novikov@heicrm.ru',      '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (7,  'Морозова Елена',     'morozova@heicrm.ru',     '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (8,  'Волков Артём',       'volkov@heicrm.ru',       '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (9,  'Лебедева Ольга',     'lebedeva@heicrm.ru',     '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (10, 'Соколов Максим',     'sokolov@heicrm.ru',      '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (11, 'Попова Наталья',     'popova@heicrm.ru',       '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (12, 'Кузнецов Андрей',    'kuznetsov@heicrm.ru',    '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (13, 'Фёдорова Татьяна',   'fedorova@heicrm.ru',     '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (14, 'Михайлов Владимир',  'mikhailov@heicrm.ru',    '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, false),
    (15, 'Егорова Ирина',      'egorova@heicrm.ru',      '$2a$10$UmBhHPqAIfchKGrg/uQMcOTFsJP4uiFLGsjcUMk64geFB3/Pqupz2', NULL, 0, true)
ON CONFLICT (email) DO NOTHING;

SELECT setval('users_id_seq', GREATEST((SELECT MAX(id) FROM users), 15));

-- ============================================================
-- User Profiles (15 профилей)
-- ============================================================
INSERT INTO user_profiles (user_id, first_name, last_name, middle_name, phone, student_id, room_id, avatar_url, date_of_birth) VALUES
    (1,  'Администратор', 'Системный',  NULL,            '+79001000001', NULL,       NULL, NULL, '1990-01-01'),
    (2,  'Алексей',       'Петров',     'Игоревич',      '+79001000002', NULL,       NULL, NULL, '1988-03-15'),
    (3,  'Мария',         'Сидорова',   'Александровна',  '+79001000003', NULL,       NULL, NULL, '1992-07-22'),
    (4,  'Дмитрий',       'Иванов',     'Сергеевич',     '+79001000004', 'STU-001',  NULL, NULL, '2003-05-10'),
    (5,  'Анна',          'Козлова',    'Петровна',      '+79001000005', 'STU-002',  NULL, NULL, '2002-11-28'),
    (6,  'Сергей',        'Новиков',    'Дмитриевич',    '+79001000006', 'STU-003',  NULL, NULL, '2003-02-14'),
    (7,  'Елена',         'Морозова',   'Викторовна',    '+79001000007', 'STU-004',  NULL, NULL, '2001-08-05'),
    (8,  'Артём',         'Волков',     'Андреевич',     '+79001000008', 'STU-005',  NULL, NULL, '2002-12-30'),
    (9,  'Ольга',         'Лебедева',   'Николаевна',    '+79001000009', 'STU-006',  NULL, NULL, '2003-04-17'),
    (10, 'Максим',        'Соколов',    'Олегович',      '+79001000010', 'STU-007',  NULL, NULL, '2002-09-03'),
    (11, 'Наталья',       'Попова',     'Ивановна',      '+79001000011', 'STU-008',  NULL, NULL, '2001-06-21'),
    (12, 'Андрей',        'Кузнецов',   'Владимирович',  '+79001000012', 'STU-009',  NULL, NULL, '2003-01-09'),
    (13, 'Татьяна',       'Фёдорова',   'Михайловна',    '+79001000013', 'STU-010',  NULL, NULL, '2002-10-12'),
    (14, 'Владимир',      'Михайлов',   'Алексеевич',    '+79001000014', 'STU-011',  NULL, NULL, '2003-07-07'),
    (15, 'Ирина',         'Егорова',    'Сергеевна',     '+79001000015', 'STU-012',  NULL, NULL, '2001-03-25')
ON CONFLICT (user_id) DO NOTHING;

-- ============================================================
-- Buildings (3 здания)
-- ============================================================
INSERT INTO buildings (id, address, floors, description) VALUES
    (1, 'ул. Ленина, 10',      5, 'Общежитие №1 — основной корпус'),
    (2, 'ул. Пушкина, 25',     4, 'Общежитие №2 — новый корпус'),
    (3, 'пр. Мира, 8',         3, 'Общежитие №3 — малый корпус')
ON CONFLICT (id) DO NOTHING;

SELECT setval('buildings_id_seq', GREATEST((SELECT MAX(id) FROM buildings), 3));

-- ============================================================
-- Rooms (15 комнат, по 5 на здание)
-- ============================================================
INSERT INTO rooms (id, building_id, room_number, floor, capacity, room_type, status) VALUES
    (1,  1, '101', 1, 1, 'single', 'occupied'),
    (2,  1, '102', 1, 2, 'double', 'occupied'),
    (3,  1, '201', 2, 2, 'double', 'occupied'),
    (4,  1, '301', 3, 4, 'block',  'occupied'),
    (5,  1, '401', 4, 2, 'double', 'free'),
    (6,  2, '101', 1, 1, 'single', 'occupied'),
    (7,  2, '102', 1, 2, 'double', 'occupied'),
    (8,  2, '201', 2, 2, 'double', 'free'),
    (9,  2, '301', 3, 4, 'block',  'occupied'),
    (10, 2, '401', 4, 2, 'double', 'free'),
    (11, 3, '101', 1, 1, 'single', 'occupied'),
    (12, 3, '102', 1, 2, 'double', 'free'),
    (13, 3, '201', 2, 2, 'double', 'occupied'),
    (14, 3, '202', 2, 4, 'block',  'occupied'),
    (15, 3, '301', 3, 2, 'double', 'free')
ON CONFLICT (building_id, room_number) DO NOTHING;

SELECT setval('rooms_id_seq', GREATEST((SELECT MAX(id) FROM rooms), 15));

-- ============================================================
-- Residents (15 жильцов)
-- ============================================================
INSERT INTO residents (id, room_id, full_name, birth_date, email, phone, move_in_date, move_out_date) VALUES
    (1,  1,  'Иванов Дмитрий Сергеевич',     '2003-05-10', 'ivanov@heicrm.ru',     '+79001000004', '2025-09-01', NULL),
    (2,  2,  'Козлова Анна Петровна',          '2002-11-28', 'kozlova@heicrm.ru',    '+79001000005', '2025-09-01', NULL),
    (3,  2,  'Новиков Сергей Дмитриевич',      '2003-02-14', 'novikov@heicrm.ru',    '+79001000006', '2025-09-01', NULL),
    (4,  3,  'Морозова Елена Викторовна',      '2001-08-05', 'morozova@heicrm.ru',   '+79001000007', '2025-09-01', NULL),
    (5,  3,  'Волков Артём Андреевич',          '2002-12-30', 'volkov@heicrm.ru',     '+79001000008', '2025-09-01', NULL),
    (6,  4,  'Лебедева Ольга Николаевна',      '2003-04-17', 'lebedeva@heicrm.ru',   '+79001000009', '2025-09-01', NULL),
    (7,  4,  'Соколов Максим Олегович',        '2002-09-03', 'sokolov@heicrm.ru',    '+79001000010', '2025-09-01', NULL),
    (8,  6,  'Попова Наталья Ивановна',        '2001-06-21', 'popova@heicrm.ru',     '+79001000011', '2025-09-01', NULL),
    (9,  7,  'Кузнецов Андрей Владимирович',   '2003-01-09', 'kuznetsov@heicrm.ru',  '+79001000012', '2025-09-01', NULL),
    (10, 7,  'Фёдорова Татьяна Михайловна',    '2002-10-12', 'fedorova@heicrm.ru',   '+79001000013', '2025-09-01', NULL),
    (11, 9,  'Михайлов Владимир Алексеевич',   '2003-07-07', 'mikhailov@heicrm.ru',  '+79001000014', '2025-09-01', NULL),
    (12, 9,  'Егорова Ирина Сергеевна',        '2001-03-25', 'egorova@heicrm.ru',    '+79001000015', '2025-09-01', NULL),
    (13, 11, 'Смирнов Павел Олегович',         '2002-04-20', 'smirnov@example.com',  '+79001000016', '2025-09-01', NULL),
    (14, 13, 'Васильева Екатерина Дмитриевна', '2003-08-11', 'vasilieva@example.com', '+79001000017', '2025-09-01', NULL),
    (15, 14, 'Николаев Роман Викторович',      '2002-06-03', 'nikolaev@example.com',  '+79001000018', '2025-09-01', NULL)
ON CONFLICT (id) DO NOTHING;

SELECT setval('residents_id_seq', GREATEST((SELECT MAX(id) FROM residents), 15));

-- ============================================================
-- Tasks (15 задач)
-- ============================================================
INSERT INTO tasks (id, author_id, assignee_id, room_id, task_type, description, priority, status) VALUES
    (1,  4,  2,    1,  'Ремонт',        'Сломан кран в ванной',                          'high',     'assigned'),
    (2,  5,  NULL,  2,  'Уборка',        'Требуется генеральная уборка комнаты',          'low',      'new'),
    (3,  6,  2,    3,  'IT-поддержка',   'Не работает Wi-Fi роутер',                      'high',     'in_progress'),
    (4,  7,  3,    4,  'Ремонт',        'Протекает потолок после дождя',                 'critical', 'assigned'),
    (5,  8,  2,    6,  'Электрика',     'Не работает розетка у окна',                    'medium',   'completed'),
    (6,  9,  NULL,  7,  'Уборка',        'Засор в раковине на кухне',                     'medium',   'new'),
    (7,  10, 3,    9,  'Ремонт',        'Сломана дверная ручка',                         'low',      'assigned'),
    (8,  11, 2,    11, 'IT-поддержка',   'Нужна замена сетевого кабеля',                  'medium',   'in_progress'),
    (9,  12, NULL,  13, 'Ремонт',        'Трещина в оконном стекле',                      'high',     'new'),
    (10, 13, 3,    14, 'Электрика',     'Мигает свет в коридоре',                        'medium',   'assigned'),
    (11, 14, 2,    2,  'Сантехника',    'Подтекает бачок унитаза',                       'high',     'in_progress'),
    (12, 15, NULL,  3,  'Уборка',        'Плесень на стене в ванной',                     'critical', 'new'),
    (13, 4,  3,    4,  'IT-поддержка',   'Установить антивирус на ноутбук',               'low',      'completed'),
    (14, 5,  2,    7,  'Ремонт',        'Скрипит пол у входной двери',                   'medium',   'assigned'),
    (15, 6,  NULL,  9,  'Мебель',        'Нужна замена матраса — продавлен',              'medium',   'new')
ON CONFLICT (id) DO NOTHING;

SELECT setval('tasks_id_seq', GREATEST((SELECT MAX(id) FROM tasks), 15));

-- ============================================================
-- Task History (15 записей)
-- ============================================================
INSERT INTO task_history (id, task_id, previous_status, new_status, changed_by, changed_at, comment) VALUES
    (1,  1,  '',           'new',         4,  '2025-10-01 09:00:00+03', 'Задача создана'),
    (2,  1,  'new',        'assigned',    1,  '2025-10-01 10:00:00+03', 'Назначена на менеджера Петрова'),
    (3,  2,  '',           'new',         5,  '2025-10-02 08:30:00+03', 'Задача создана'),
    (4,  3,  '',           'new',         6,  '2025-10-02 11:00:00+03', 'Задача создана'),
    (5,  3,  'new',        'assigned',    1,  '2025-10-02 12:00:00+03', 'Назначена на Петрова'),
    (6,  3,  'assigned',   'in_progress', 2,  '2025-10-02 14:00:00+03', 'Начал диагностику роутера'),
    (7,  4,  '',           'new',         7,  '2025-10-03 07:00:00+03', 'Задача создана'),
    (8,  4,  'new',        'assigned',    1,  '2025-10-03 08:00:00+03', 'Срочно! Назначена на Сидорову'),
    (9,  5,  '',           'new',         8,  '2025-10-04 09:00:00+03', 'Задача создана'),
    (10, 5,  'new',        'assigned',    1,  '2025-10-04 09:30:00+03', NULL),
    (11, 5,  'assigned',   'in_progress', 2,  '2025-10-04 10:00:00+03', 'Взял в работу'),
    (12, 5,  'in_progress','completed',   2,  '2025-10-04 15:00:00+03', 'Розетка заменена'),
    (13, 13, '',           'new',         4,  '2025-10-10 08:00:00+03', 'Задача создана'),
    (14, 13, 'new',        'assigned',    1,  '2025-10-10 09:00:00+03', NULL),
    (15, 13, 'assigned',   'completed',   3,  '2025-10-10 16:00:00+03', 'Антивирус установлен')
ON CONFLICT (id) DO NOTHING;

SELECT setval('task_history_id_seq', GREATEST((SELECT MAX(id) FROM task_history), 15));

-- ============================================================
-- Task Comments (15 комментариев)
-- ============================================================
INSERT INTO task_comments (id, task_id, author_id, comment_text, created_at) VALUES
    (1,  1,  4,  'Кран течёт уже второй день, вода на полу',                '2025-10-01 09:05:00+03'),
    (2,  1,  2,  'Посмотрю сегодня после обеда',                            '2025-10-01 10:30:00+03'),
    (3,  3,  6,  'Роутер перестал работать после грозы',                     '2025-10-02 11:05:00+03'),
    (4,  3,  2,  'Скорее всего сгорел блок питания, закажу новый',          '2025-10-02 14:30:00+03'),
    (5,  4,  7,  'Протечка усиливается при дожде, нужны вёдра',            '2025-10-03 07:10:00+03'),
    (6,  4,  1,  'Вызываем аварийную службу',                               '2025-10-03 08:15:00+03'),
    (7,  4,  3,  'Аварийная бригада приедет завтра к 9:00',                 '2025-10-03 12:00:00+03'),
    (8,  5,  8,  'Розетка искрит, опасно пользоваться',                     '2025-10-04 09:05:00+03'),
    (9,  5,  2,  'Готово, заменил на новую розетку с заземлением',           '2025-10-04 15:00:00+03'),
    (10, 8,  11, 'Кабель перетёрся у стола',                                '2025-10-06 10:00:00+03'),
    (11, 8,  2,  'Принесу новый кабель завтра',                             '2025-10-06 11:00:00+03'),
    (12, 9,  12, 'Трещина появилась после сильного ветра',                  '2025-10-07 08:00:00+03'),
    (13, 11, 14, 'Бачок подтекает постоянно, расход воды большой',          '2025-10-08 09:00:00+03'),
    (14, 12, 15, 'Плесень растёт уже месяц, вентиляция не помогает',       '2025-10-09 07:30:00+03'),
    (15, 15, 6,  'Матрас совсем продавился по центру, спать невозможно',    '2025-10-11 08:00:00+03')
ON CONFLICT (id) DO NOTHING;

SELECT setval('task_comments_id_seq', GREATEST((SELECT MAX(id) FROM task_comments), 15));

-- ============================================================
-- User Activity Log (15 записей)
-- ============================================================
INSERT INTO user_activity_log (id, user_id, action, details, ip_address, user_agent, created_at) VALUES
    (1,  1,  'login',            '{"method": "password"}',                       '192.168.1.1',   'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',  '2025-10-01 08:00:00+03'),
    (2,  1,  'user_registered',  '{"target_user_id": 4, "email": "ivanov@heicrm.ru"}', '192.168.1.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', '2025-10-01 08:10:00+03'),
    (3,  4,  'login',            '{"method": "password"}',                       '192.168.1.10',  'Mozilla/5.0 (Linux; Android 13)',             '2025-10-01 08:30:00+03'),
    (4,  4,  'profile_updated',  '{"fields": ["phone", "first_name"]}',          '192.168.1.10',  'Mozilla/5.0 (Linux; Android 13)',             '2025-10-01 08:35:00+03'),
    (5,  2,  'login',            '{"method": "password"}',                       '192.168.1.5',   'Mozilla/5.0 (Macintosh; Intel Mac OS X)',     '2025-10-01 09:00:00+03'),
    (6,  5,  'login',            '{"method": "password"}',                       '192.168.1.11',  'Mozilla/5.0 (iPhone; CPU iPhone OS 17_0)',    '2025-10-02 08:00:00+03'),
    (7,  5,  'profile_updated',  '{"fields": ["last_name"]}',                    '192.168.1.11',  'Mozilla/5.0 (iPhone; CPU iPhone OS 17_0)',    '2025-10-02 08:10:00+03'),
    (8,  1,  'user_updated',     '{"target_user_id": 6, "fields": ["role_id"]}', '192.168.1.1',   'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',  '2025-10-02 10:00:00+03'),
    (9,  3,  'login',            '{"method": "password"}',                       '192.168.1.6',   'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',  '2025-10-03 07:30:00+03'),
    (10, 6,  'login',            '{"method": "password"}',                       '192.168.1.12',  'Mozilla/5.0 (Linux; Android 14)',             '2025-10-03 08:00:00+03'),
    (11, 7,  'login',            '{"method": "password"}',                       '192.168.1.13',  'Mozilla/5.0 (Macintosh; Intel Mac OS X)',     '2025-10-03 09:00:00+03'),
    (12, 1,  'token_refresh',    '{}',                                           '192.168.1.1',   'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',  '2025-10-03 12:00:00+03'),
    (13, 8,  'login',            '{"method": "password"}',                       '192.168.1.14',  'Mozilla/5.0 (X11; Linux x86_64)',             '2025-10-04 08:00:00+03'),
    (14, 1,  'user_deleted',     '{"target_user_id": 99}',                       '192.168.1.1',   'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',  '2025-10-05 16:00:00+03'),
    (15, 9,  'login',            '{"method": "password"}',                       '192.168.1.15',  'Mozilla/5.0 (iPad; CPU OS 17_0)',             '2025-10-06 07:00:00+03')
ON CONFLICT (id) DO NOTHING;

SELECT setval('user_activity_log_id_seq', GREATEST((SELECT MAX(id) FROM user_activity_log), 15));

-- +goose Down
DELETE FROM user_activity_log WHERE id BETWEEN 1 AND 15;
DELETE FROM task_comments WHERE id BETWEEN 1 AND 15;
DELETE FROM task_history WHERE id BETWEEN 1 AND 15;
DELETE FROM tasks WHERE id BETWEEN 1 AND 15;
DELETE FROM residents WHERE id BETWEEN 1 AND 15;
DELETE FROM rooms WHERE id BETWEEN 1 AND 15;
DELETE FROM buildings WHERE id BETWEEN 1 AND 15;
DELETE FROM user_profiles WHERE user_id BETWEEN 1 AND 15;
DELETE FROM users WHERE id BETWEEN 1 AND 15;
