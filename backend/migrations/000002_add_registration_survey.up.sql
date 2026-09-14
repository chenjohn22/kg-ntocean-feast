CREATE TABLE activity_restaurants (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    category VARCHAR(100) NOT NULL,
    name VARCHAR(150) NOT NULL,
    location VARCHAR(255) NOT NULL DEFAULT '',
    sort_order INT UNSIGNED NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_activity_restaurants_identity (category, name, location),
    KEY idx_activity_restaurants_display (active, category, sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE registrations
    ADD COLUMN ticket_source VARCHAR(32) NULL AFTER email,
    ADD COLUMN activity_restaurant_id BIGINT UNSIGNED NULL AFTER ticket_source,
    ADD COLUMN satisfaction TINYINT UNSIGNED NULL AFTER activity_restaurant_id,
    ADD COLUMN suggestion TEXT NULL AFTER satisfaction,
    ADD KEY idx_registrations_ticket_source (ticket_source),
    ADD KEY idx_registrations_activity_restaurant (activity_restaurant_id),
    ADD KEY idx_registrations_satisfaction (satisfaction),
    ADD CONSTRAINT fk_registrations_activity_restaurant
        FOREIGN KEY (activity_restaurant_id) REFERENCES activity_restaurants (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    ADD CONSTRAINT chk_registrations_satisfaction
        CHECK (satisfaction IS NULL OR satisfaction BETWEEN 1 AND 5);
