-- +goose Up 
ALTER TABLE users
ADD COLUMN access_level VARCHAR(30) NOT NULL DEFAULT 'user';

-- +goose Down
ALTER TABLE users
DROP COLUMN access_level;