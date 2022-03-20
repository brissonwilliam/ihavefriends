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
			IFNULL(bill.cumulative_total, 0.0) AS cumulative_total,
			IFNULL(bill.highest_total, 0.0) AS highest_total,
			IFNULL(bill.last_total, 0.0) AS last_total,
			bill.last_updated AS last_updated
		FROM user
		LEFT JOIN bill ON bill.user_id = user.id
		ORDER BY cumulative_total DESC
	`

	queryUpdateBill = `
		UPDATE bill SET 
			last_total = :new_bill_total,
			cumulative_total = cumulative_total + :new_bill_total,
			second_highest_total = IF(:new_bill_total >= highest_total,highest_total,second_highest_total),
			highest_total = GREATEST(:new_bill_total, highest_total),
			last_updated = now()
		WHERE user_id = :user_id
	`

	queryUndoLastUserBill = `
		UPDATE bill SET 
			cumulative_total = cumulative_total - last_total,
			highest_total = IF(highest_total = last_total, second_highest_total, highest_total),
			second_highest_total = IF(highest_total = last_total, 0.0, second_highest_total),
			last_total = 0.0,
			last_updated = now()
		WHERE user_id = ?
	`

	// only run this before updating so that last_updated keeps its right value
	queryEnsureBillEntryExists = `
		INSERT INTO bill(user_id) VALUES(?)
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

func (r defaultUserRepository) UpdateUserBill(update models.BillUpdate) error {
	_, err := r.db.Exec(queryEnsureBillEntryExists, update.UserId)
	if err != nil {
		return err
	}

	rowsUpdated, err := sqlx.NamedExec(r.db, queryUpdateBill, &update)
	if err != nil {
		return err
	}
	n, errRows := rowsUpdated.RowsAffected()
	if errRows != nil {
		logger.Get().Error(err)
		return errRows
	}
	if n < 1 {
		return core.NewErrNotFound("user " + update.UserId.String())
	}

	return nil
}

func (r defaultUserRepository) UndoLastUserBill(userId uuid.OrderedUUID) error {
	_, err := r.db.Exec(queryEnsureBillEntryExists, userId)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(queryUndoLastUserBill, userId)
	if err != nil {
		return err
	}

	return nil
}
