-- +goose Up
CREATE table users (
    user_id UUID PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    notification_methods JSONB NOT NULL DEFAULT '[]'::jsonb
);

-- +goose Down
drop table users