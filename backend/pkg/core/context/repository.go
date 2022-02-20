package context

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/analytics"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/user"
	"github.com/jmoiron/sqlx"
)

func UserRepository(db *sqlx.DB) user.UserRepository {
	return user.NewUserRepository(db)
}

func AnalyticsRepository(db *sqlx.DB) analytics.AnalyticsRepository {
	return analytics.NewAnlyticsRepository(db)
}
