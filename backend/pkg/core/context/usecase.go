package context

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/usecase/auth"
	"github.com/jmoiron/sqlx"
)

func NewAuthenticateUsecase(db *sqlx.DB) auth.Usecase {
	return auth.NewUsecase(UserRepository(db))
}
