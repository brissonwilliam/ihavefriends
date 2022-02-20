package http

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/bonneFete"
	user "github.com/brissonwilliam/ihavefriends/backend/pkg/http/users"
	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	Auth      auth.Handler
	BonneFete bonneFete.Handler
	User      user.Handler
}

func NewHandlers(db *sqlx.DB) Handlers {
	return Handlers{
		Auth:      auth.NewHandler(db),
		BonneFete: bonneFete.NewHandler(db),
		User:      user.NewHandler(db),
	}
}
