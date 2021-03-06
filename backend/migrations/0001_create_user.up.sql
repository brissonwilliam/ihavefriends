CREATE TABLE user (
    id BINARY(16) NOT NULL PRIMARY KEY,
    username VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    email VARCHAR(320) DEFAULT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_ut BIGINT UNSIGNED NOT NULL DEFAULT 0,

    UNIQUE KEY uk_user_email (email),
    UNIQUE KEY uk_username (username)
);