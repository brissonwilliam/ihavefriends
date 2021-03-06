CREATE TABLE bonnefete(
    user_id BINARY(16) NOT NULL PRIMARY KEY,
    total BIGINT NOT NULL DEFAULT 0,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY fk_bonnefete_user_id (user_id) REFERENCES user(id)
);