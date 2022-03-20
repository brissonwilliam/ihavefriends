package bill

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h defaultHandler) Get(ctx echo.Context) error {
	jwtClaims := auth.GetJWTClaimsFromContext(ctx)
	userId := jwtClaims.Id

	bfAnalytics, err := h.usecase.GetAnalytics(userId)
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, bfAnalytics)
}
