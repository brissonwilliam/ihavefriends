package bonneFete

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	// the http request upgrader to a websocket
	upgrader = websocket.Upgrader{}
)

func (h *defaultHandler) GetWebSocket(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	h.wsHub.Subscribe <- &WSClient{
		ws:  ws,
		hub: h.wsHub,
		msg: make(chan []byte, maxMessageSize),
	}

	return nil
}

func (h defaultHandler) Get(ctx echo.Context) error {
	bfAnalytics, err := h.usecase.GetAnalytics()
	if err != nil {
		return httperr.FromCoreErr(ctx, err)
	}
	return ctx.JSON(http.StatusOK, bfAnalytics)
}
