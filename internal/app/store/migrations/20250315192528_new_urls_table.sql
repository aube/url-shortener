-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
        id serial PRIMARY KEY,
        short_url CHAR(10) UNIQUE,
        original_url TEXT);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists urls;
-- +goose StatementEnd
