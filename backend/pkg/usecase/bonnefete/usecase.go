package bonnefete

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/analytics"
)

type Usecase interface {
	GetAnalytics() (*models.BonneFeteAnalytics, error)
	Increment(userId uuid.OrderedUUID) (*models.BonneFeteAnalytics, error)
}

func NewUsecase(txProvider storage.TxProvider, analyticsRepo analytics.AnalyticsRepository) Usecase {
	return defaultUsecase{
		repo:       analyticsRepo,
		txProvider: txProvider,
	}
}

type defaultUsecase struct {
	repo       analytics.AnalyticsRepository
	txProvider storage.TxProvider
}

func (u defaultUsecase) GetAnalytics() (*models.BonneFeteAnalytics, error) {
	bfTotalByUsers, err := u.repo.GetTotalByUsers()
	if err != nil {
		return nil, err
	}

	return computeAnalytics(bfTotalByUsers), nil
}

func (u defaultUsecase) Increment(userId uuid.OrderedUUID) (analytics *models.BonneFeteAnalytics, err error) {
	var uow storage.UnitOfWork
	if uow, err = u.txProvider.Begin(); err != nil {
		return nil, err
	}

	defer u.txProvider.Close(uow, &err)

	err = u.repo.WithUnitOfWork(uow).IncrementBF(userId)
	if err != nil {
		return nil, err
	}

	bfTotalByUsers, err := u.repo.WithUnitOfWork(uow).GetTotalByUsers()
	if err != nil {
		return nil, err
	}

	return computeAnalytics(bfTotalByUsers), nil
}

func computeAnalytics(bfTotalByUsers []models.BFTotalByUser) *models.BonneFeteAnalytics {
	total := uint(0)
	for _, t := range bfTotalByUsers {
		total += t.Count
	}

	return &models.BonneFeteAnalytics{
		Total:        total,
		TotalByUsers: bfTotalByUsers,
	}
}
