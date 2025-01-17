-- +goose Up
CREATE TABLE password_tokens(
    id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at  TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE   
);

-- +goose Down
DROP TABLE password_tokens;