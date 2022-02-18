package auth

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/context"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/validator"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/usecase/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type defaultHandler struct {
	usecase   auth.Usecase
	validator validator.Validator
}

// Handler implements the analytics-related controllers
type Handler interface {
	Post(ctx echo.Context) error
}

// NewHandler returns a new instance of a handler supporting analytics
func NewHandler(db *sqlx.DB) Handler {
	return defaultHandler{
		usecase:   context.NewAuthenticateUsecase(db),
		validator: validator.Get(),
	}
}
