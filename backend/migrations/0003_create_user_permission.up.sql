CREATE TABLE user_permission (
    user_id BINARY(16) NOT NULL,
    permission VARCHAR(50) NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY pk_user_permission (user_id, permission),
    FOREIGN KEY fk_user_permission_user_id (user_id) REFERENCES user(id)
);