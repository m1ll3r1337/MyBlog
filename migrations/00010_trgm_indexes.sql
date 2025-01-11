-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX posts_title_trgm_idx ON posts USING gist (title gist_trgm_ops);
CREATE INDEX posts_description_trgm_idx ON posts USING gist (description gist_trgm_ops);
CREATE INDEX tags_name_trgm_idx ON tags USING gist (name gist_trgm_ops);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX posts_title_trgm_idx;
DROP INDEX posts_description_trgm_idx;
DROP INDEX tags_name_trgm_idx;

-- +goose StatementEnd
