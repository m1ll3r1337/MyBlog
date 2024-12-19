-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN description TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN description;
-- +goose StatementEnd
