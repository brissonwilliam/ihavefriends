package server

import (
	config "github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/db"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/http"
	"github.com/labstack/echo/v4"
)

type WebServer interface {
	Start() error
}

func NewWebServer() WebServer {
	return &defaultServer{}
}

type defaultServer struct{}

func (d defaultServer) Start() error {
	cfg := config.GetConfig()

	logger.Get().WithField("cfg", cfg).Info("Starting web server")

	db, err := db.Connect(cfg.DB)
	if err != nil {
		logger.Get().WithError(err).Error("Could not connect to db")
		return err
	}

	e := echo.New()

	handlers := http.NewHandlers(db)
	api.NewRouter(handlers).Configure(e)

	return e.Start(cfg.Web.Address)
}
