CREATE TABLE bill (
    user_id BINARY(16) NOT NULL PRIMARY KEY,
    cumulative_total FLOAT(2) NOT NULL DEFAULT 0,
    highest_total FLOAT(2) NOT NULL DEFAULT 0,
    second_highest_total FLOAT(2) NOT NULL DEFAULT 0,
    last_total FLOAT(2) NOT NULL DEFAULT 0,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY fk_bill_user_id (user_id) REFERENCES user(id)
);