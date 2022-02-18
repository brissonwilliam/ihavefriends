package logger

import (
	"github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

// Log contains the logger unique instance
var log *logrus.Logger

// Get Returns the main logger instance
func Get() *logrus.Logger {
	if log == nil {
		cfg := config.GetLogger()

		log = logrus.New()
		log.SetLevel(cfg.Level)
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	return log
}

func WithStack() *logrus.Entry {
	return Get().WithField("stack", string(debug.Stack()))
}
