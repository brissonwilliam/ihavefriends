package auth

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/labstack/echo/v4"
	nethttp "net/http"
)

func (h defaultHandler) GetPublicUsers(ctx echo.Context) error {
	users, err := h.usecase.GetPublicUsers()
	if err != nil {
		return httperr.GetHttpError(ctx, err)
	}

	return ctx.JSON(nethttp.StatusOK, users)
}
