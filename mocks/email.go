package mocks

import (
	"github.com/stretchr/testify/mock"
)

type EmailMock struct {
	mock.Mock
}

func NewEmailMock() *EmailMock {
	return &EmailMock{}
}

func (m *EmailMock) SendEmail(cardsTable, cardsPriceFormatted string) error {
	args := m.Called(cardsTable, cardsPriceFormatted)
	return args.Error(0)
}
