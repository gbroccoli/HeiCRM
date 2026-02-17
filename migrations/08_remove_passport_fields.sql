-- +goose Up
ALTER TABLE residents DROP COLUMN IF EXISTS passport_series;
ALTER TABLE residents DROP COLUMN IF EXISTS passport_number;

-- +goose Down
ALTER TABLE residents ADD COLUMN passport_series VARCHAR(10);
ALTER TABLE residents ADD COLUMN passport_number VARCHAR(20);
