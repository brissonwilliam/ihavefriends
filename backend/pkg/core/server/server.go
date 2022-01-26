package server

import (
	"encoding/json"
	config "github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/labstack/echo/v4"
)

type WebServer interface{
	Start() error
}

func NewWebServer() WebServer {
	return &defaultServer{}
}

type defaultServer struct {

}

func (d defaultServer) Start() error {
	cfg := config.GetConfig()

	b, err := json.Marshal(cfg)
	if err == nil {
		logger.Get().WithField("config", string(b)).Info("Starting web server")
	}

	// TODO: connect to db

	// TODO: configure routes

	e := echo.New()

	return e.Start(cfg.Web.Address)
}




