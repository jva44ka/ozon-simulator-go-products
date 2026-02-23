-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
    ADD COLUMN count SMALLINT NOT NULL DEFAULT 10;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
    DROP COLUMN count;
-- +goose StatementEnd