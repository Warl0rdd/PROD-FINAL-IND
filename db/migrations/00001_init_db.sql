-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE gender AS ENUM ('MALE', 'FEMALE');

CREATE TABLE clients
(
    id       uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    login    text             NOT NULL UNIQUE,
    age      int              NOT NULL,
    location text             NOT NULL,
    gender   gender           NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE clients;
DROP TYPE gender;
-- +goose StatementEnd
