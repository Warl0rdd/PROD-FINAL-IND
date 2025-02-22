-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE advertisers
(
    id   uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name text             NOT NULL
);
CREATE TABLE ml_scores
(
    id            uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    client_id     uuid             NOT NULL,
    advertiser_id uuid             NOT NULL,
    score         int              NOT NULL,
    FOREIGN KEY (client_id) REFERENCES clients (id),
    FOREIGN KEY (advertiser_id) REFERENCES advertisers (id),
    UNIQUE (client_id, advertiser_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE ml_scores;
DROP TABLE advertisers;
-- +goose StatementEnd
