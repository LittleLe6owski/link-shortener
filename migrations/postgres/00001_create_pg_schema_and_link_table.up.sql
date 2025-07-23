-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS link_shortener;
CREATE TABLE IF NOT EXISTS link_shortener.links (
    id uuid PRIMARY KEY,
    full_uri TEXT NOT NULL,
    short_uri BIGINT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    is_deleted bool NOT NULL DEFAULT false
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE link_shortener.links;
DROP SCHEMA link_shortener;
-- +goose StatementEnd