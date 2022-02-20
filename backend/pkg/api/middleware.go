package api

import (
	"github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func GlobalMiddlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{middleware.Logger(), middleware.CORS(), middleware.BodyLimit("2M")}
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey:    auth.CTX_KEY,
		Claims:        &auth.JWTClaims{},
		SigningKey:    []byte(config.GetWeb().JwtKey),
		SigningMethod: middleware.AlgorithmHS256,
	})
}

func Permissions(requiredPermissions ...string) echo.MiddlewareFunc {
	return func(nextHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user := auth.GetJWTClaimsFromContext(ctx)

			if user == nil {
				return ctx.String(http.StatusUnauthorized, "Unauthorized")
			}

			if auth.HasSuperAdminPermission(user.Permissions) {
				return nextHandler(ctx)
			}

			for _, p := range requiredPermissions {
				if !auth.HasPermission(p, user.Permissions) {
					return ctx.String(http.StatusForbidden, "Forbidden")
				}
			}

			return nextHandler(ctx)
		}
	}
}
