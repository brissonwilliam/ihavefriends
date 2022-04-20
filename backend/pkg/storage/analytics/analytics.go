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

	// TODO: split BF and billing analytics
	GetBillTotalsAllUsers() ([]models.BillUserTotal, error)
	GetUserBillTotalsByTime(userId uuid.OrderedUUID) (models.BillUserTotalsByTime, error)
	CreateBill(update models.CreateBill) error
	UpdateUserAgBillTotal(newBill models.CreateBill) error
	DeleteLastUserBill(userId uuid.OrderedUUID) error
	RecomputeUserBillAggregates(userId uuid.OrderedUUID) error

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
