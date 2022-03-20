package bill

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/analytics"
)

type Usecase interface {
	GetAnalytics(userId uuid.OrderedUUID) (*models.BillAnalytics, error)
	UpdateUserBill(update models.BillUpdate) (*models.BillAnalytics, error)
	UndoLastUserBill(userId uuid.OrderedUUID) (*models.BillAnalytics, error)
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

func (u defaultUsecase) GetAnalytics(userId uuid.OrderedUUID) (*models.BillAnalytics, error) {
	totalBillsByUser, err := u.repo.GetTotalBillsByUser()
	if err != nil {
		return nil, err
	}

	return computeAnalytics(totalBillsByUser, userId), nil
}

func (u defaultUsecase) UndoLastUserBill(userId uuid.OrderedUUID) (analytics *models.BillAnalytics, err error) {
	var uow storage.UnitOfWork
	if uow, err = u.txProvider.Begin(); err != nil {
		return nil, err
	}

	defer u.txProvider.Close(uow, &err)

	err = u.repo.WithUnitOfWork(uow).UndoLastUserBill(userId)
	if err != nil {
		return nil, err
	}

	var totalBillsByUser []models.BillUserTotal
	totalBillsByUser, err = u.repo.WithUnitOfWork(uow).GetTotalBillsByUser()
	if err != nil {
		return nil, err
	}

	return computeAnalytics(totalBillsByUser, userId), nil
}

func (u defaultUsecase) UpdateUserBill(update models.BillUpdate) (analytics *models.BillAnalytics, err error) {
	var uow storage.UnitOfWork
	if uow, err = u.txProvider.Begin(); err != nil {
		return nil, err
	}

	defer u.txProvider.Close(uow, &err)

	err = u.repo.WithUnitOfWork(uow).UpdateUserBill(update)
	if err != nil {
		return nil, err
	}

	var totalBillsByUser []models.BillUserTotal
	totalBillsByUser, err = u.repo.WithUnitOfWork(uow).GetTotalBillsByUser()
	if err != nil {
		return nil, err
	}

	return computeAnalytics(totalBillsByUser, update.UserId), nil
}

func computeAnalytics(billUserTotals []models.BillUserTotal, userId uuid.OrderedUUID) *models.BillAnalytics {
	var userTotal models.BillUserTotal
	grandTotal := models.Amount(0.0)
	for i, ut := range billUserTotals {
		grandTotal += ut.CumulativeTotal

		if ut.UserId == userId {
			userTotal = billUserTotals[i]
		}
	}

	return &models.BillAnalytics{
		GrandTotal:   grandTotal,
		TotalByUsers: billUserTotals,
		UserTotal:    userTotal,
	}
}
