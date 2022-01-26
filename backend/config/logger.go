package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Logger struct {
	Level logrus.Level
}

func GetLogger() Logger {
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	return Logger {
		Level: level,
	}
}