-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_tags
(
    tag_id INT,
    post_id INT,
    FOREIGN KEY (post_id) REFERENCES Posts (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES Tags (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post_tags;
-- +goose StatementEnd
