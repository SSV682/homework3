package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *logrus.Logger
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(l)
	logger.Out = os.Stdout

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == logrus.ErrorLevel {
		l.Debug(message, args...)
	} else {
		l.msg("error", message, args...)
	}
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)
	os.Exit(1)
}

func (l *Logger) log(message string, args ...interface{}) {
	l.logger.Info(message, args)
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
