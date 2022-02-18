package http

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/bonneFete"
	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	Auth      auth.Handler
	BonneFete bonneFete.Handler
}

func NewHandlers(db *sqlx.DB) Handlers {
	return Handlers{
		Auth:      auth.NewHandler(db),
		BonneFete: bonneFete.NewHandler(db),
	}
}
