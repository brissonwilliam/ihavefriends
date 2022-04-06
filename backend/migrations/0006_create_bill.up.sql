ALTER TABLE bill RENAME bill_ag_user;
ALTER TABLE bill_ag_user
    DROP COLUMN second_highest_total;

CREATE TABLE bill (
    id BINARY(16) NOT NULL PRIMARY KEY,
    user_id BINARY(16) NOT NULL,
    total FLOAT(2) NOT NULL,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX (created),
    FOREIGN KEY `fk_bill_user_id_2` (user_id) REFERENCES user(id)
);

INSERT INTO bill (id, user_id, total, created)
    SELECT UNHEX(REPLACE(UUID(),'-','')), user_id, cumulative_total, created
    FROM bill_ag_user;

UPDATE bill_ag_user SET last_total = cumulative_total;