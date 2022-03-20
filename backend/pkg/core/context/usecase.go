package context

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/usecase/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/usecase/bill"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/usecase/bonnefete"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/usecase/user"
	"github.com/jmoiron/sqlx"
)

func NewAuthUsecase(db *sqlx.DB) auth.Usecase {
	return auth.NewUsecase(UserRepository(db))
}

func NewUserUsecase(db *sqlx.DB) user.Usecase {
	return user.NewUsecase(storage.NewTxProvider(db), UserRepository(db))
}

func NewBonneFeteUsecase(db *sqlx.DB) bonnefete.Usecase {
	return bonnefete.NewUsecase(storage.NewTxProvider(db), AnalyticsRepository(db))
}

func NewBillUsecase(db *sqlx.DB) bill.Usecase {
	return bill.NewUsecase(storage.NewTxProvider(db), AnalyticsRepository(db))
}
