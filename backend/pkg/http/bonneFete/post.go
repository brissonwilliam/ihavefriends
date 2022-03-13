package bonneFete

import (
	"encoding/json"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *defaultHandler) Post(ctx echo.Context) error {
	jwtClaims := auth.GetJWTClaimsFromContext(ctx)

	userId := jwtClaims.Id

	bfAnalytics, err := h.usecase.Increment(userId)
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}

	analyticsJson, err := json.Marshal(bfAnalytics)
	if err != nil {
		logger.Get().WithError(err).Error("Could not parse analytics to json. Will not send to websocket broadcast")
	} else {
		h.wsHub.BroadcastMsg <- analyticsJson
	}

	return ctx.JSON(http.StatusOK, bfAnalytics)
}

func (h *defaultHandler) ResetCount(ctx echo.Context) error {
	jwtClaims := auth.GetJWTClaimsFromContext(ctx)

	userId := jwtClaims.Id

	bfAnalytics, err := h.usecase.ResetCount(userId)
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}

	analyticsJson, err := json.Marshal(bfAnalytics)
	if err != nil {
		logger.Get().WithError(err).Error("Could not parse analytics to json. Will not send to websocket broadcast")
	} else {
		h.wsHub.BroadcastMsg <- analyticsJson
	}

	return ctx.JSON(http.StatusOK, bfAnalytics)
}
