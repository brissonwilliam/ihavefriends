package bonneFete

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http/httperr"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

var (
	// the http request upgrader to a websocket
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if strings.Contains(origin, "127.0.0.1") {
				return true
			}
			if strings.Contains(origin, "sourpusss.com") {
				return true
			}
			if strings.Contains(origin, "localhost") {
				return true
			}
			return false
		},
	}
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
