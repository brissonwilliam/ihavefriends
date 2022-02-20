package analytics

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
	"github.com/jmoiron/sqlx"
)

const (
	queryGetBFTotalByUser = `
		SELECT
			user.username AS username,
			IFNULL(bonnefete.total, 0) AS total
		FROM user
		LEFT JOIN bonnefete ON bonnefete.user_id = user.id
	`

	queryUpdateIncrementBF = `
		UPDATE bonnefete SET total = total + 1 WHERE user_id = ?
	`

	queryEnsureBFEntryExists = `
		INSERT INTO bonnefete (user_id, total) VALUES(?, 0)
		ON DUPLICATE KEY UPDATE
			last_updated = NOW()
	`
)

type AnalyticsRepository interface {
	GetTotalByUsers() ([]models.BFTotalByUser, error)
	IncrementBF(userId uuid.OrderedUUID) error
	WithUnitOfWork(uow storage.UnitOfWork) AnalyticsRepository
}

type defaultUserRepository struct {
	db sqlx.Ext
}

func NewAnlyticsRepository(db sqlx.Ext) AnalyticsRepository {
	return defaultUserRepository{
		db: db,
	}
}

func (r defaultUserRepository) WithUnitOfWork(uow storage.UnitOfWork) AnalyticsRepository {
	tx := storage.UnitAsTransaction(uow)
	return NewAnlyticsRepository(tx)
}

func (r defaultUserRepository) GetTotalByUsers() ([]models.BFTotalByUser, error) {
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
	if n, _ := rowsUpdated.RowsAffected(); n < 1 {
		return core.NewErrNotFound("user " + userId.String())
	}

	return nil
}
