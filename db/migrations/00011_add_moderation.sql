-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE campaigns ADD COLUMN approved boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE campaigns DROP COLUMN approved;
-- +goose StatementEnd
