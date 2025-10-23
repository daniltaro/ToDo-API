-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id          uuid PRIMARY KEY,
    title       text NOT NULL,
    description text,
    is_done     boolean NOT NULL DEFAULT false,
    deadline    timestamptz,
    created_at  timestamptz DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  IF EXISTS tasks;
-- +goose StatementEnd
