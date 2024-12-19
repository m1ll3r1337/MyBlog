-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags
(
    id      SERIAL PRIMARY KEY,
    name   VARCHAR(20) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tags;
-- +goose StatementEnd
