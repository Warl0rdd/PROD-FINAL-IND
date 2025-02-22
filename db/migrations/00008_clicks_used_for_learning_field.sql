-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE impressions ADD COLUMN used_for_learning boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE impressions DROP COLUMN used_for_learning;
-- +goose StatementEnd
