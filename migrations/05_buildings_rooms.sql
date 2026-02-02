-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS buildings (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    address varchar(500),
    floors int NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    id bigserial PRIMARY KEY,
    building_id bigint NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    number varchar(50) NOT NULL,
    floor int NOT NULL,
    capacity int NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(building_id, number)
);

CREATE INDEX IF NOT EXISTS idx_rooms_building_id ON rooms(building_id);
CREATE INDEX IF NOT EXISTS idx_rooms_floor ON rooms(floor);

ALTER TABLE user_profiles
    ADD CONSTRAINT fk_user_profiles_room_id
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_profiles DROP CONSTRAINT IF EXISTS fk_user_profiles_room_id;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS buildings;
-- +goose StatementEnd
