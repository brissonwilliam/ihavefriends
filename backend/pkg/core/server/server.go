package server

import (
	"errors"
	"fmt"
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

	if err := checkConfig(cfg); err != nil {
		logger.Get().Error(err.Error())
		return err
	}

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

func checkConfig(cfg config.CompleteConfig) error {
	baseErrMsg := "invalid JwtKey. Make sure %s env var is set"
	if cfg.Web.JwtKey == "" {
		return errors.New(fmt.Sprintf(baseErrMsg, "web_jwtkey"))
	}
	if cfg.DB.Username == "" {
		return errors.New(fmt.Sprintf(baseErrMsg, "db_username"))
	}
	if cfg.DB.Host == "" {
		return errors.New(fmt.Sprintf(baseErrMsg, "db_host"))
	}
	if cfg.DB.Host == "" {
		return errors.New(fmt.Sprintf(baseErrMsg, "db_name"))
	}
	return nil
}
