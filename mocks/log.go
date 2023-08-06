package mocks

import (
	"mtg-report/internal/sources/logger/logrus"

	"github.com/stretchr/testify/mock"
)

type CustomMock struct {
	mock.Mock
}

func NewCustomMock() *CustomMock {
	return &CustomMock{}
}

func (c *CustomMock) Fatal(args ...interface{}) {
	c.Called(args)
}

func (c *CustomMock) Warn(args ...interface{}) {
	c.Called(args)
}

func (c *CustomMock) Error(args ...interface{}) {
	c.Called(args)
}

func (c *CustomMock) Info(args ...interface{}) {
	c.Called(args)
}

type LogMock struct {
	mock.Mock
}

func NewLogMock() *LogMock {
	return &LogMock{}
}

func (l *LogMock) Info(args ...interface{}) {
	l.Called(args)
}

func (l *LogMock) WithFields(fields logrus.Fields) logrus.CustomEntry {
	args := l.Called(fields)
	return args.Get(0).(logrus.CustomEntry)
}

func (l *LogMock) WithError(err error) logrus.CustomEntry {
	args := l.Called(err)
	return args.Get(0).(logrus.CustomEntry)
}

func (l *LogMock) Error(args ...interface{}) {
	l.Called(args)
}

func (l *LogMock) Warn(args ...interface{}) {
	l.Called(args)
}
