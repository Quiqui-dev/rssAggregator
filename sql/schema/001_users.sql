-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    create_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    display_name TEXT NOT NULL,
    email_address TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL);

-- +goose Down

DROP TABLE users;