package bonneFete

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type defaultHandler struct {
	usecase interface{}
	validator          validator.Validator
}

// Handler implements the analytics-related controllers
type Handler interface {
	Post(ctx echo.Context) error
	Get(ctx echo.Context) error
}

// NewHandler returns a new instance of a handler supporting analytics
func NewHandler(db *sqlx.DB) Handler {
	return defaultHandler{
		validator:          validator.Get(),
	}
}

