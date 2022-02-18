package auth

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/labstack/echo/v4"
	nethttp "net/http"
)

func (h defaultHandler) Post(ctx echo.Context) error {
	var form models.AuthForm
	if err := ctx.Bind(&form); err != nil {
		return httperr.UnreadableForm(ctx, err)
	}

	userWithCreds, err := h.usecase.Authenticate(form)
	if err != nil {
		return httperr.GetHttpError(ctx, err)
	}

	return ctx.JSON(nethttp.StatusOK, userWithCreds)
}
