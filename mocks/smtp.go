package mocks

import (
	"net/smtp"

	"github.com/stretchr/testify/mock"
)

type SMTPAuthMock struct {
	mock.Mock
}

func NewSMTPAuthMock() *SMTPAuthMock {
	return &SMTPAuthMock{}
}

func (s *SMTPAuthMock) Start(server *smtp.ServerInfo) (string, []byte, error) {
	args := s.Called(server)
	return args.String(0), args.Get(1).([]byte), args.Error(2)
}

func (s *SMTPAuthMock) Next(fromServer []byte, more bool) ([]byte, error) {
	args := s.Called(fromServer, more)
	return args.Get(0).([]byte), args.Error(1)
}
