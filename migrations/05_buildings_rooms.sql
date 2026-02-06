-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS buildings (
    id bigserial PRIMARY KEY,
    address varchar(255) NOT NULL,
    floors integer NOT NULL,
    description text,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    id bigserial PRIMARY KEY,
    building_id bigint NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    room_number varchar(10) NOT NULL,
    floor integer NOT NULL,
    capacity integer NOT NULL,
    room_type varchar(20) NOT NULL CHECK (room_type IN ('single', 'double', 'block')),
    status varchar(20) NOT NULL DEFAULT 'free' CHECK (status IN ('free', 'occupied')),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(building_id, room_number)
);

CREATE INDEX IF NOT EXISTS idx_rooms_building_id ON rooms(building_id);
CREATE INDEX IF NOT EXISTS idx_rooms_status ON rooms(status);

CREATE TABLE IF NOT EXISTS residents (
    id bigserial PRIMARY KEY,
    room_id bigint NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    full_name varchar(200) NOT NULL,
    birth_date date NOT NULL,
    passport_series varchar(10),
    passport_number varchar(20),
    email varchar(100),
    phone varchar(20),
    move_in_date date NOT NULL,
    move_out_date date,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_residents_room_id ON residents(room_id);
CREATE INDEX IF NOT EXISTS idx_residents_full_name ON residents(full_name);
CREATE INDEX IF NOT EXISTS idx_residents_phone ON residents(phone);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS residents;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS buildings;
-- +goose StatementEnd
