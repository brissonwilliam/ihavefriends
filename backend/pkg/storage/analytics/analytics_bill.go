package analytics

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/jmoiron/sqlx"
)

const (
	queryGetBillTotalByUser = `
		SELECT
			user.id AS user_id,
			user.username AS username,
			IFNULL(bill_ag_user.cumulative_total, 0.0) AS cumulative_total,
			IFNULL(bill_ag_user.highest_total, 0.0) AS highest_total,
			IFNULL(bill_ag_user.last_total, 0.0) AS last_total,
			bill_ag_user.last_updated AS last_updated
		FROM user
		LEFT JOIN bill_ag_user ON bill_ag_user.user_id = user.id
		ORDER BY cumulative_total DESC
	`

	queryUpdateBillAgUser = `
		UPDATE bill_ag_user SET 
			last_total = :new_bill_total,
			cumulative_total = cumulative_total + :new_bill_total,
			highest_total = GREATEST(:new_bill_total, highest_total),
			last_updated = now()
		WHERE user_id = :user_id
	`

	queryInsertBill = `
		INSERT INTO bill (id, user_id, total) VALUES (:id, :user_id, :new_bill_total)
	`

	queryRecomputeUserAgBill = `
		INSERT INTO bill_ag_user (user_id, last_total, highest_total, cumulative_total)
			SELECT 
				user.id AS user_id,
				IFNULL(lt.total, 0) AS last_total,
				IFNULL(nht.highest_total, 0) AS highest_total,
			    IFNULL(ct.total, 0) AS cumulative_total
			FROM user
			LEFT JOIN (
				SELECT user_id, total FROM bill WHERE user_id = :user_id ORDER BY created DESC LIMIT 1
			)  lt ON lt.user_id = user.id
			LEFT JOIN (
				SELECT user_id, max(total)  AS highest_total FROM bill WHERE user_id = :user_id GROUP BY user_id
			) nht ON nht.user_id = user.id
			LEFT JOIN (
			    SELECT user_id, sum(total) AS total FROM bill WHERE user_id = :user_id GROUP BY user_id
			) ct ON ct.user_id = user.id
			WHERE user.id = :user_id
		ON DUPLICATE KEY UPDATE
			last_updated = NOW(),
		    cumulative_total = VALUES(cumulative_total),
			highest_total = VALUES(highest_total),
			last_total = VALUES(last_total)
	`

	queryDeleteUserLastBill = `
		DELETE FROM bill WHERE id = (SELECT id FROM (SELECT id FROM bill WHERE user_id = ? ORDER BY created DESC LIMIT 1) t)
	`

	// only run this before updating so that last_updated keeps its right value
	queryEnsureBillAgUserEntryExists = `
		INSERT INTO bill_ag_user(user_id) VALUES(?)
		ON DUPLICATE KEY UPDATE last_updated = NOW()
	`
)

func (r defaultUserRepository) GetTotalBillsByUser() ([]models.BillUserTotal, error) {
	totalByUsers := []models.BillUserTotal{}
	err := sqlx.Select(r.db, &totalByUsers, queryGetBillTotalByUser)
	if err != nil {
		return nil, err
	}
	return totalByUsers, nil
}

func (r defaultUserRepository) CreateBill(bill models.CreateBill) error {
	// insert the new bill in its raw form
	_, err := sqlx.NamedExec(r.db, queryInsertBill, &bill)
	if err != nil {
		return err
	}

	return nil
}

func (r defaultUserRepository) UpdateUserAgBillTotal(newBill models.CreateBill) error {
	_, err := r.db.Exec(queryEnsureBillAgUserEntryExists, newBill.UserId)
	if err != nil {
		return err
	}

	rowsUpdated, err := sqlx.NamedExec(r.db, queryUpdateBillAgUser, &newBill)
	if err != nil {
		return err
	}
	n, errRows := rowsUpdated.RowsAffected()
	if errRows != nil {
		logger.Get().Error(err)
		return errRows
	}
	if n < 1 {
		return core.NewErrNotFound("user " + newBill.UserId.String())
	}

	return nil
}

func (r defaultUserRepository) DeleteLastUserBill(userId uuid.OrderedUUID) error {
	_, err := r.db.Exec(queryDeleteUserLastBill, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r defaultUserRepository) RecomputeUserBillAggregates(userId uuid.OrderedUUID) error {
	arg := struct {
		UserId uuid.OrderedUUID `db:"user_id"`
	}{UserId: userId}
	_, err := sqlx.NamedExec(r.db, queryRecomputeUserAgBill, arg)
	if err != nil {
		return err
	}

	return nil
}
