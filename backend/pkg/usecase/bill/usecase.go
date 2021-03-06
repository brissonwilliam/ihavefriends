package bill

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/analytics"
)

type Usecase interface {
	GetAnalytics(userId uuid.OrderedUUID) (*models.BillAnalytics, error)
	CreateBill(newBill models.CreateBill) (*models.BillAnalytics, error)
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
	var uow storage.UnitOfWork
	var err error
	if uow, err = u.txProvider.Begin(); err != nil {
		return nil, err
	}

	defer u.txProvider.Close(uow, &err)

	return u.getAnalytics(&uow, userId)
}

func (u defaultUsecase) getAnalytics(uow *storage.UnitOfWork, userId uuid.OrderedUUID) (*models.BillAnalytics, error) {
	totalBillsByUser, err := u.repo.WithUnitOfWork(*uow).GetBillTotalsAllUsers()
	if err != nil {
		return nil, err
	}

	totalsByTime, err := u.repo.WithUnitOfWork(*uow).GetUserBillTotalsByTime(userId)
	if err != nil {
		return nil, err
	}

	return computeAnalytics(totalBillsByUser, totalsByTime, userId), nil
}

func (u defaultUsecase) UndoLastUserBill(userId uuid.OrderedUUID) (analytics *models.BillAnalytics, err error) {
	var uow storage.UnitOfWork
	if uow, err = u.txProvider.Begin(); err != nil {
		return nil, err
	}

	defer u.txProvider.Close(uow, &err)

	err = u.repo.WithUnitOfWork(uow).DeleteLastUserBill(userId)
	if err != nil {
		return nil, err
	}

	err = u.repo.WithUnitOfWork(uow).RecomputeUserBillAggregates(userId)
	if err != nil {
		return nil, err
	}

	return u.getAnalytics(&uow, userId)
}

func (u defaultUsecase) CreateBill(newBill models.CreateBill) (analytics *models.BillAnalytics, err error) {
	var uow storage.UnitOfWork
	if uow, err = u.txProvider.Begin(); err != nil {
		return nil, err
	}

	defer u.txProvider.Close(uow, &err)

	err = u.repo.WithUnitOfWork(uow).CreateBill(newBill)
	if err != nil {
		return nil, err
	}

	err = u.repo.WithUnitOfWork(uow).UpdateUserAgBillTotal(newBill)
	if err != nil {
		return nil, err
	}

	return u.getAnalytics(&uow, newBill.UserId)
}

func computeAnalytics(billUserTotals []models.BillUserTotal, totalsByTime models.BillUserTotalsByTime, userId uuid.OrderedUUID) *models.BillAnalytics {
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
		TotalsByTime: totalsByTime,
	}
}
