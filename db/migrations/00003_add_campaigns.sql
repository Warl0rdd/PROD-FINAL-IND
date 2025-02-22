-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE campaign_gender AS ENUM ('MALE', 'FEMALE', 'ALL');
CREATE TABLE campaigns
(
    id                  uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    advertiser_id       uuid             NOT NULL,
    impression_limit    int              NOT NULL,
    clicks_limit        int              NOT NULL,
    cost_per_impression float            NOT NULL,
    cost_per_click      float            NOT NULL,
    ad_title            text             NOT NULL,
    ad_text             text             NOT NULL,
    start_date          int              NOT NULL,
    end_date            int              NOT NULL,
    gender              campaign_gender  NOT NULL DEFAULT 'ALL',
    age_from            int,
    age_to              int,
    location            text,
    FOREIGN KEY (advertiser_id) REFERENCES advertisers (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE campaigns;
DROP TYPE campaign_gender;
-- +goose StatementEnd
