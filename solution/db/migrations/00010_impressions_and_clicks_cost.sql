-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE impressions ADD COLUMN cost float NOT NULL DEFAULT 0;
ALTER TABLE clicks ADD COLUMN cost float NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE impressions DROP COLUMN cost;
ALTER TABLE clicks DROP COLUMN cost;
-- +goose StatementEnd
