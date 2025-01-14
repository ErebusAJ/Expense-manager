-- +goose Up
ALTER TABLE users
ADD COLUMN token_version INT DEFAULT 0,
ADD COLUMN last_logged_in TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN token_version,
DROP COLUMN last_logged_in;
