CREATE TABLE lottery_codes (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    code VARCHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_lottery_codes_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE registrations (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    lottery_code_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    email VARCHAR(254) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_registrations_lottery_code (lottery_code_id),
    KEY idx_registrations_name (name),
    KEY idx_registrations_phone (phone),
    KEY idx_registrations_email (email),
    KEY idx_registrations_registered_at (registered_at),
    CONSTRAINT fk_registrations_lottery_code
        FOREIGN KEY (lottery_code_id) REFERENCES lottery_codes (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
