-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Comments
(
    id      SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    user_id INT,
    post_id INT,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Posts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Comments;
-- +goose StatementEnd
