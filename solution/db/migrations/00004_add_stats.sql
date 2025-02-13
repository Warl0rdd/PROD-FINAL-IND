-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE impressions
(
    id          uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    campaign_id uuid             NOT NULL,
    client_id   uuid             NOT NULL,
    day         int              NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaigns (id),
    FOREIGN KEY (client_id) REFERENCES clients (id),
    UNIQUE (campaign_id, client_id)
);
CREATE TABLE clicks
(
    id          uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    campaign_id uuid             NOT NULL,
    client_id   uuid             NOT NULL,
    day         int              NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaigns (id),
    FOREIGN KEY (client_id) REFERENCES clients (id),
    UNIQUE (campaign_id, client_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE clicks;
DROP TABLE impressions;
-- +goose StatementEnd
