-- +goose Down
-- +goose StatementBegin
DROP TABLE link_shortener.links 

DROP SCHEMA link_shortener;
-- +goose StatementEnd