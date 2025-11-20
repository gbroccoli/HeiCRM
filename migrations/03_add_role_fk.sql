-- +goose Up
-- +goose StatementBegin
-- Add foreign key constraint from users.role_id to roles.id
ALTER TABLE users ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove foreign key constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_role;
-- +goose StatementEnd
