-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id smallserial PRIMARY KEY,
    name varchar(50) NOT NULL UNIQUE,
    description text,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- Insert default roles
INSERT INTO roles (id, name, description) VALUES
    (0, 'user', 'Regular user with basic permissions'),
    (1, 'admin', 'Administrator with full system access'),
    (2, 'manager', 'Manager with elevated permissions')
ON CONFLICT (id) DO NOTHING;

-- Index for faster role name lookups
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles CASCADE;
-- +goose StatementEnd
