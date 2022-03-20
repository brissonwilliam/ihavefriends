package analytics

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/jmoiron/sqlx"
)

const (
	queryGetBFTotalByUser = `
		SELECT
			user.username AS username,
			IFNULL(bonnefete.total, 0) AS total
		FROM user
		LEFT JOIN bonnefete ON bonnefete.user_id = user.id
		ORDER BY total DESC
	`

	queryUpdateIncrementBF = `
		UPDATE bonnefete SET total = total + 1, last_updated = now() WHERE user_id = ?
	`

	queryUpdateResetBFCount = `
		UPDATE bonnefete SET total = 0, last_updated = now() WHERE user_id = ?
	`

	// only run this before updating so that last_updated keeps its right value
	queryEnsureBFEntryExists = `
		INSERT INTO bonnefete (user_id, total) VALUES(?, 0)
		ON DUPLICATE KEY UPDATE
			last_updated = NOW()
	`
)

func (r defaultUserRepository) GetTotalBFByUsers() ([]models.BFTotalByUser, error) {
	totalByUsers := []models.BFTotalByUser{}
	err := sqlx.Select(r.db, &totalByUsers, queryGetBFTotalByUser)
	if err != nil {
		return nil, err
	}
	return totalByUsers, nil
}

func (r defaultUserRepository) IncrementBF(userId uuid.OrderedUUID) error {
	_, err := r.db.Exec(queryEnsureBFEntryExists, userId)
	if err != nil {
		return err
	}

	rowsUpdated, err := r.db.Exec(queryUpdateIncrementBF, userId)
	if err != nil {
		return err
	}
	n, errRows := rowsUpdated.RowsAffected()
	if errRows != nil {
		logger.Get().Error(err)
		return errRows
	}
	if n < 1 {
		return core.NewErrNotFound("user " + userId.String())
	}

	return nil
}

func (r defaultUserRepository) ResetCount(userId uuid.OrderedUUID) error {
	_, err := r.db.Exec(queryEnsureBFEntryExists, userId)
	if err != nil {
		return err
	}

	rowsUpdated, err := r.db.Exec(queryUpdateResetBFCount, userId)
	if err != nil {
		return err
	}
	n, errRows := rowsUpdated.RowsAffected()
	if errRows != nil {
		logger.Get().Error(err)
		return errRows
	}
	if n < 1 {
		return core.NewErrNotFound("user " + userId.String())
	}

	return nil
}
