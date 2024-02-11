-- +goose Up
-- +goose StatementBegin
ALTER TABLE todo
ALTER COLUMN description TYPE VARCHAR(1000),
ALTER COLUMN description SET DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE todo
ALTER COLUMN description TYPE TEXT
ALTER COLUMN description DROP DEFAULT;
-- +goose StatementEnd
