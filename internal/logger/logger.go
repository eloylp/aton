package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
}

type BasicLogger struct {
	l *logrus.Logger
}

func NewBasicLogger(w io.Writer) *BasicLogger { //nolint:golint
	logger := logrus.New()
	logger.SetOutput(w)
	return &BasicLogger{
		l: logger,
	}
}

func (b *BasicLogger) Infof(format string, args ...interface{}) {
	b.l.Infof(format, args...)
}

func (b *BasicLogger) Info(args ...interface{}) {
	b.l.Info(args...)
}

func (b *BasicLogger) Errorf(format string, args ...interface{}) {
	b.l.Errorf(format, args...)
}

func (b *BasicLogger) Error(args ...interface{}) {
	b.l.Error(args...)
}
