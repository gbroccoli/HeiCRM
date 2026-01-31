-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name varchar(255),
    last_name varchar(255),
    middle_name varchar(255),
    phone varchar(50),
    student_id varchar(100),
    room_id bigint,
    avatar_url varchar(500),
    date_of_birth date,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_room_id ON user_profiles(room_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_student_id ON user_profiles(student_id);

CREATE TABLE IF NOT EXISTS user_activity_log (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action varchar(255) NOT NULL,
    details jsonb,
    ip_address varchar(45),
    user_agent varchar(500),
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_activity_log_user_id ON user_activity_log(user_id);
CREATE INDEX IF NOT EXISTS idx_user_activity_log_created_at ON user_activity_log(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_activity_log;
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd
