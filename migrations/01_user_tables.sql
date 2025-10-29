CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    name varchar not null,
    email varchar not null,
    password varchar(255) not null,
    tg_send boolean default false,
    created_at timestamptz,
    updated_at timestamptz
)