package mocks

import "github.com/stretchr/testify/mock"

type timerMock struct {
	mock.Mock
}

func NewTimerMock() *timerMock {
	return &timerMock{}
}

func (t *timerMock) Now() string {
	argsMock := t.Called()
	return argsMock.Get(0).(string)
}
