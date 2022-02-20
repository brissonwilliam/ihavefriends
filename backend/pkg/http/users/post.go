package user

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h defaultHandler) Post(ctx echo.Context) error {
	var form models.CreaterUserForm
	if err := ctx.Bind(&form); err != nil {
		return httperr.UnreadableForm(ctx, err)
	}

	if err := h.validator.Struct(form); err != nil {
		return httperr.InvalidRequest(ctx, err)
	}

	user, err := h.usecase.CreateUser(form)
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, user)
}
