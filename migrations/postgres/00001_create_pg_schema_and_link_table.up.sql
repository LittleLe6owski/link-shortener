-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS link_shortener;

CREATE TABLE IF NOT EXISTS link_shortener.links (
    id SERIAL PRIMARY KEY,
    full_uri TEXT NOT NULL,
    short_uri VARCHAR(16) UNIQUE NOT NULL,
    create_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    is_deleted bool NOT NULL DEFAULT false
);

-- +goose StatementEnd