package analytics

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	endOfDayHour = 9 // UTC, so that's 4AM EST or 5AM EDT (east time with daylight's saving)
)

const (
	queryGetBillTotalsAllUsers = `
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

	// furthest we're aggregating the data is 1 month ago, so add that in WHERE to optimize query (created is indexed)
	queryGetUserBillTotalByTime = `
		SELECT 
			user.id AS user_id,
			SUM(IF(bill.created >= DATE_SUB(now(), INTERVAL 2 DAY), bill.total, 0)) AS last_48h_total,
			SUM(IF(bill.created >= :start_of_last_week AND bill.created < :start_of_current_week, bill.total, 0)) AS last_week_total,
			SUM(IF(bill.created >= :start_of_current_week AND bill.created <= :end_of_day, bill.total, 0)) AS this_week_total,
			SUM(IF(bill.created >= :start_of_month AND bill.created <= :end_of_day, bill.total, 0)) AS this_month_total
		FROM bill
		JOIN user ON user.id = bill.user_id
		WHERE user.id = :user_id AND bill.created >= :furthest_created_lookup
		GROUP BY user.id
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

func (r defaultUserRepository) GetBillTotalsAllUsers() ([]models.BillUserTotal, error) {
	totalByUsers := []models.BillUserTotal{}
	err := sqlx.Select(r.db, &totalByUsers, queryGetBillTotalsAllUsers)
	if err != nil {
		return nil, err
	}
	return totalByUsers, nil
}

func (r defaultUserRepository) GetUserBillTotalsByTime(userId uuid.OrderedUUID) (models.BillUserTotalsByTime, error) {
	currentTime := time.Now().UTC()

	startOfWeek := getStartOfWeek(currentTime)
	startOfMonth := getStartOfMonth(currentTime)
	currentTimeEndOfDay := getEndOfDay(currentTime)
	queryCriteria := queryGetUserBillTotalsByTimeCriteria{
		StartOfWeek:           startOfWeek,
		EndOfDay:              currentTimeEndOfDay,
		StartOfLastWeek:       startOfWeek.Add(-1 * time.Hour * 24 * 7),
		StartOfMonth:          startOfMonth,
		FurthestCreatedLookup: currentTime.Add(-1 * time.Hour * 24 * 32),
		UserId:                userId,
	}
	query, args, _ := r.db.BindNamed(queryGetUserBillTotalByTime, queryCriteria)

	totalsByTime := []models.BillUserTotalsByTime{}
	err := sqlx.Select(r.db, &totalsByTime, query, args...)
	if err != nil {
		return models.BillUserTotalsByTime{}, err
	}
	if len(totalsByTime) < 1 {
		return models.BillUserTotalsByTime{UserId: userId}, nil
	}

	return totalsByTime[0], nil
}

type queryGetUserBillTotalsByTimeCriteria struct {
	StartOfWeek           time.Time        `db:"start_of_current_week"`
	EndOfDay              time.Time        `db:"end_of_day"`
	StartOfLastWeek       time.Time        `db:"start_of_last_week"`
	StartOfMonth          time.Time        `db:"start_of_month"`
	FurthestCreatedLookup time.Time        `db:"furthest_created_lookup"`
	UserId                uuid.OrderedUUID `db:"user_id"`
}

// returns the datetime for the last monday, or today if today is monday
// Also aligns to endOfDayHour for better accounting (usually should be 4 or 5AM EST)
func getStartOfWeek(t time.Time) time.Time {
	year, month, day := t.Date()
	currentZeroDay := time.Date(year, month, day, endOfDayHour, 0, 0, 0, t.Location())

	weekday := time.Duration(t.Weekday())
	// offset sundays by 7 days since sunday is 0
	if weekday == 0 {
		weekday += 7
	}

	// subtract by whatever weekday it is minus 1 to be aligned with monday (weekday 1)
	monday := currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
	return monday
}

// aligns current date to first day of month
func getStartOfMonth(clientTime time.Time) time.Time {
	year, month, _ := clientTime.Date()
	return time.Date(year, month, 01, endOfDayHour, 0, 0, 0, clientTime.Location())
}

// aligns the current date to endOfDayHour if it has not been reached yet for the day, or else enfOfDayHour for tomorrow
func getEndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	endOfDay := time.Date(year, month, day, endOfDayHour, 0, 0, 0, t.Location())
	if t.Hour() >= endOfDayHour {
		endOfDay = endOfDay.Add(time.Hour * 24)
	}
	return endOfDay
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
