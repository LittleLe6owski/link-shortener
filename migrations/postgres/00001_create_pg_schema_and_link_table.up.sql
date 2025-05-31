-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS link_shortener;

CREATE TABLE IF NOT EXISTS link_shortener.links (
    id SERIAL PRIMARY KEY,
    uri TEXT NOT NULL,
    short_uri VARCHAR(16) UNIQUE NOT NULL,
    create_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    call_counter INTEGER NOT NULL DEFAULT 0
);

COMMENT ON COLUMN link_shortener.links.call_counter IS 'Счётчик, сколько раз мы обратились к данной ссылке';
-- +goose StatementEnd