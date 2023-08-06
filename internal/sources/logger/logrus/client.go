package logrus

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Level map[string]uint32 = map[string]uint32{
	"debug": 5,
	"info":  4,
	"warn":  3,
	"error": 2,
	"fatal": 1,
	"panic": 0,
}

type logrus struct {
	*log.Logger
}

func New(level string) *logrus {
	value, ok := Level[level]
	if !ok {
		value = Level["debug"]
	}

	return &logrus{
		Logger: &log.Logger{
			Out:       os.Stderr,
			Formatter: new(log.TextFormatter),
			Hooks:     make(log.LevelHooks),
			Level:     log.Level(value),
		},
	}
}

func (l *logrus) WithError(err error) CustomEntry {
	return l.Logger.WithError(err)
}

func (l *logrus) WithFields(fields Fields) CustomEntry {
	return l.Logger.WithFields(log.Fields(fields))
}

func (l *logrus) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *logrus) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *logrus) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}
