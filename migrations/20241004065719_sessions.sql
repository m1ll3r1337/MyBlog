-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Sessions
(
    id         SERIAL PRIMARY KEY,
    user_id    int UNIQUE references users (id) on delete cascade,
    token_hash text unique not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Sessions;
-- +goose StatementEnd
