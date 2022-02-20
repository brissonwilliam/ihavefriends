package analytics

import (
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
)

type AnalyticsRepository interface {
	GetTotalByUsers() ([]models.BFTotalByUser, error)
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
