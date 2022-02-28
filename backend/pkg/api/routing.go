package api

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http"
	"github.com/labstack/echo/v4"
)

type Router interface {
	Configure(e *echo.Echo)
}

func NewRouter(handlers http.Handlers) Router {
	return defaultRouter{
		h: handlers,
	}
}

type defaultRouter struct {
	h http.Handlers
}

func (r defaultRouter) Configure(e *echo.Echo) {
	e.Use(GlobalMiddlewares()...)

	g := e.Group("api/")
	r.configureAuthRoutes(g)
	r.configureUserRoutes(g)
	r.configureBonneFeteRoutes(g)
}

func (r defaultRouter) configureAuthRoutes(g *echo.Group) {
	g.POST("auth", r.h.Auth.Post)
}

func (r defaultRouter) configureUserRoutes(g *echo.Group) {
	g.GET("publicUsers", r.h.User.GetPublicUsers)
	g.POST("users", r.h.User.Post, JWTMiddleware(), Permissions(auth.PERM_ADD_USER))
}

func (r defaultRouter) configureBonneFeteRoutes(g *echo.Group) {
	g.POST("bonneFete", r.h.BonneFete.Post, JWTMiddleware())
	g.GET("bonneFete", r.h.BonneFete.Get, JWTMiddleware())
	g.GET("bonneFete/ws", r.h.BonneFete.GetWebSocket, JWTParamsMiddleware())
}
