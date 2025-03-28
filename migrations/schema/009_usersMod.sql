-- +goose Up
ALTER TABLE users
ADD COLUMN image_url TEXT;

-- +goose Down
ALTER TABLE users
DROP COLUMN image_url;