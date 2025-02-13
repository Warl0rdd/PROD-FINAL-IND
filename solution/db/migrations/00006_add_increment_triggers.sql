-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE OR REPLACE FUNCTION update_campaign_impressions_count()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE campaigns
    SET impressions_count = impressions_count + 1
    WHERE id = NEW.campaign_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_campaign_clicks_count()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE campaigns
    SET clicks_count = clicks_count + 1
    WHERE id = NEW.campaign_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_after_insert_impression
    AFTER INSERT ON impressions
    FOR EACH ROW
EXECUTE FUNCTION update_campaign_impressions_count();

CREATE TRIGGER trg_after_insert_click
    AFTER INSERT ON clicks
    FOR EACH ROW
EXECUTE FUNCTION update_campaign_clicks_count();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TRIGGER trg_after_insert_impression ON impressions;
DROP TRIGGER trg_after_insert_click ON clicks;
DROP FUNCTION update_campaign_impressions_count();
DROP FUNCTION update_campaign_clicks_count();
-- +goose StatementEnd
