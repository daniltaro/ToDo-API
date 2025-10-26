-- +goose Up
-- +goose StatementBegin
ALTER TABLE "tasks" ADD COLUMN "login" text DEFAULT 'unknown';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "tasks" DROP COLUMN "login";
-- +goose StatementEnd
