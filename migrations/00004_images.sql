-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Images
(
    id   SERIAL PRIMARY KEY,
    post_id int,
    data BYTEA NOT NULL,
    FOREIGN KEY (post_id) REFERENCES Posts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Images;
-- +goose StatementEnd
