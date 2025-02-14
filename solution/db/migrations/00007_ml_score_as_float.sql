-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE ml_scores ALTER COLUMN score TYPE float;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE ml_scores ALTER COLUMN score TYPE int;
-- +goose StatementEnd
