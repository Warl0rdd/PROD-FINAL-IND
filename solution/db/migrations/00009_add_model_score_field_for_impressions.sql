-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE impressions ADD COLUMN model_score float;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE impressions DROP COLUMN model_score;
-- +goose StatementEnd
