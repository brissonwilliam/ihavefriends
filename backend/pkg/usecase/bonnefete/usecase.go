package bonnefete

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/analytics"
)

type Usecase interface {
	GetAnalytics() (*models.BonneFeteAnalytics, error)
}

func NewUsecase(analyticsRepo analytics.AnalyticsRepository) Usecase {
	return defaultUsecase{
		repo: analyticsRepo,
	}
}

type defaultUsecase struct {
	repo analytics.AnalyticsRepository
}

func (u defaultUsecase) GetAnalytics() (*models.BonneFeteAnalytics, error) {
	bfTotalByUsers, err := u.repo.GetTotalByUsers()
	if err != nil {
		return nil, err
	}

	// compute global total
	total := uint(0)
	for _, t := range bfTotalByUsers {
		total += t.Count
	}

	return &models.BonneFeteAnalytics{
		Total:        total,
		TotalByUsers: bfTotalByUsers,
	}, nil
}
