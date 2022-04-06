package bill

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *defaultHandler) Post(ctx echo.Context) error {
	var form models.CreateBill
	if err := ctx.Bind(&form); err != nil {
		return httperr.UnreadableForm(ctx, err)
	}

	jwtClaims := auth.GetJWTClaimsFromContext(ctx)
	form.UserId = jwtClaims.Id
	form.ID = uuid.NewOrderedUUID()

	if err := h.validator.Struct(form); err != nil {
		return httperr.InvalidRequest(ctx, err)
	}

	billAnalytics, err := h.usecase.CreateBill(form)
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}

	return ctx.JSON(http.StatusOK, billAnalytics)
}

func (h *defaultHandler) PostUndoLastBill(ctx echo.Context) error {
	jwtClaims := auth.GetJWTClaimsFromContext(ctx)
	userId := jwtClaims.Id

	billAnalytics, err := h.usecase.UndoLastUserBill(userId)
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}

	return ctx.JSON(http.StatusOK, billAnalytics)
}
