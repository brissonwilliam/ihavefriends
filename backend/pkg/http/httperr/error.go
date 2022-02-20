package httperr

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func FromCoreErr(ctx echo.Context, err error) error {
	switch err.(type) {
	case core.ErrNotFound:
		ctx.String(http.StatusNotFound, err.Error())
		return err
	case core.ErrConflict:
		ctx.String(http.StatusConflict, err.Error())
		return err
	default:
		// defaults to internal error. Log but only write a generic error message in client response
		logger.Get().Error(err)
		ctx.String(http.StatusInternalServerError, "An internal error has occurred")
	}

	return err
}

func UnreadableForm(ctx echo.Context, err error) error {
	ctx.String(http.StatusBadRequest, "Unreadable form")
	return err
}

func InvalidRequest(ctx echo.Context, err error) error {
	ctx.String(http.StatusBadRequest, "Invalid request. "+err.Error())
	return err
}
