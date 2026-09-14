ALTER TABLE registrations
    DROP FOREIGN KEY fk_registrations_activity_restaurant,
    DROP CHECK chk_registrations_satisfaction,
    DROP INDEX idx_registrations_ticket_source,
    DROP INDEX idx_registrations_activity_restaurant,
    DROP INDEX idx_registrations_satisfaction,
    DROP COLUMN ticket_source,
    DROP COLUMN activity_restaurant_id,
    DROP COLUMN satisfaction,
    DROP COLUMN suggestion;

DROP TABLE IF EXISTS activity_restaurants;
