package bonneFete

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h defaultHandler) Get(ctx echo.Context) error {
	bfAnalytics, err := h.usecase.GetAnalytics()
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, bfAnalytics)
}
