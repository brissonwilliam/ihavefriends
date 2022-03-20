package analytics

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
	"github.com/jmoiron/sqlx"
)

type AnalyticsRepository interface {
	GetTotalBFByUsers() ([]models.BFTotalByUser, error)
	IncrementBF(userId uuid.OrderedUUID) error
	ResetCount(userId uuid.OrderedUUID) error

	GetTotalBillsByUser() ([]models.BillUserTotal, error)
	UpdateUserBill(update models.BillUpdate) error
	UndoLastUserBill(userId uuid.OrderedUUID) error

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
