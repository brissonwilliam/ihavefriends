package bill

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/context"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/validator"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/usecase/bill"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"time"
)

type defaultHandler struct {
	usecase   bill.Usecase
	validator validator.Validator
}

// Handler implements the analytics-related controllers
type Handler interface {
	Post(ctx echo.Context) error
	PostUndoLastBill(ctx echo.Context) error
	Get(ctx echo.Context) error
}

// NewHandler returns a new instance of a handler supporting analytics
func NewHandler(db *sqlx.DB) Handler {
	return &defaultHandler{
		validator: validator.Get(),
		usecase:   context.NewBillUsecase(db),
	}
}

// returns the client time when given in query params, defaults to server time
func getClientTime(ctx echo.Context) time.Time {
	param := ctx.Param("time")
	if param == "" {
		return time.Now()
	}

	t, err := time.Parse(time.RFC3339, param)
	if err != nil {
		logger.Get().Error(err)
		return time.Now()
	}

	return t
}
