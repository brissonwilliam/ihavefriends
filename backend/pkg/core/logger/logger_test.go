package logger

import (
	"errors"
	"testing"
)

func TestLoggerWithStack(t *testing.T) {
	WithStack().Error(errors.New("some error"))
}
