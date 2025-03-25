-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
        id serial PRIMARY KEY,
        user_id int,
        short_url CHAR(10) UNIQUE,
        original_url TEXT,
        deleted boolean default false
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists urls;
-- +goose StatementEnd
