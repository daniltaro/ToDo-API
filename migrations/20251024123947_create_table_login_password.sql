-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    login text NOT NULL UNIQUE,
    password text NOT NULL,
    created_at timestamptz DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS users;
-- +goose StatementEnd
