package bonneFete

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func(h defaultHandler) Get(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, `{"content": "test!"}`)
}
