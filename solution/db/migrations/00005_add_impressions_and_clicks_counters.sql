-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE campaigns
    ADD COLUMN
        impressions_count int NOT NULL DEFAULT 0,
    ADD COLUMN
        clicks_count      int NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE campaigns
    DROP COLUMN impressions_count,
    DROP COLUMN clicks_count;
-- +goose StatementEnd
