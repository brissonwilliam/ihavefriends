package api

import (
	"errors"
	"github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func GlobalMiddlewares() []echo.MiddlewareFunc {
	rl := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      10,
		Burst:     20,
		ExpiresIn: time.Minute * 3,
	})
	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	})
	return []echo.MiddlewareFunc{middleware.Logger(), cors, middleware.BodyLimit("2M"), middleware.RateLimiter(rl)}
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey:    auth.CTX_KEY,
		Claims:        &auth.JWTClaims{},
		SigningKey:    []byte(config.GetWeb().JwtKey),
		SigningMethod: middleware.AlgorithmHS256,
	})
}

func JWTParamsMiddleware() echo.MiddlewareFunc {
	return func(nextHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			tokenStr := ctx.QueryParam("token")
			if tokenStr == "" {
				return ctx.String(http.StatusUnauthorized, "Unauthorized")
			}

			t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, errors.New("invalid token signature")
				}
				return []byte(config.GetWeb().JwtKey), nil
			})

			if err != nil {
				logger.Get().Error(err)
				return ctx.String(http.StatusUnauthorized, "Unauthorized")
			}

			if !t.Valid {
				return ctx.String(http.StatusUnauthorized, "Unauthorized")
			}

			return nextHandler(ctx)
		}
	}
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
