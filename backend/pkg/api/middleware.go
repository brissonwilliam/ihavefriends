package api

import (
	"github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GlobalMiddlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{middleware.CORS()}
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWT(config.GetWeb().JwtKey)
}
