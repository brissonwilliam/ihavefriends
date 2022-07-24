/* Bill highest totals, when indexes are recomputed, are messed up because of the adjustment.
   When bills got segmented, only 1 bill was created with the cumulative total.

   ** CONSIDERING ALL HIGHEST TOTALS IN bill_ag_user ARE THE CORRECT VALUE**
   This marks the adjustment bill with a flag (so it can be ignored in queries),
   removes the highest total from the bill amount and creates a new bill with the highest total
 */

ALTER TABLE bill
    ADD column temp_adjust_highest_tot FLOAT(2) DEFAULT NULL,
    ADD COLUMN is_adjustment BOOLEAN NOT NULL DEFAULT FALSE;

INSERT INTO bill (id, total, user_id, temp_adjust_highest_tot)
SELECT t.id, t.total, t.user_id, t.highest_total FROM (
         SELECT bill.id AS id, bill.total AS total, bill.user_id AS user_id, bagu.highest_total AS highest_total
         FROM bill
         JOIN bill_ag_user bagu ON bagu.user_id = bill.user_id
         WHERE (bill.user_id, bill.created) IN
               (SELECT user_id, min(created) FROM bill GROUP BY user_id)
           AND bill.created < '2022-04-07'
     ) t
ON DUPLICATE KEY UPDATE
    temp_adjust_highest_tot = VALUES(temp_adjust_highest_tot),
    is_adjustment = true,
    total = VALUES(total) - VALUES(temp_adjust_highest_tot);

INSERT INTO bill (id, total, user_id, created)
SELECT UNHEX(REPLACE(UUID(),'-','')), bill.temp_adjust_highest_tot, bill.user_id, bill.created
FROM bill
WHERE bill.temp_adjust_highest_tot IS NOT NULL;

ALTER TABLE bill
    DROP COLUMN temp_adjust_highest_tot;