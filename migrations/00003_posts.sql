-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Posts
(
    id      SERIAL PRIMARY KEY,
    title   VARCHAR(255) NOT NULL,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES Users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Posts;
-- +goose StatementEnd
