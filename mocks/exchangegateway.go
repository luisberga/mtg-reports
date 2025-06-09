package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ExchangeGatewayMock struct {
	mock.Mock
}

func NewExchangeGatewayMock() *ExchangeGatewayMock {
	return &ExchangeGatewayMock{}
}

func (m *ExchangeGatewayMock) GetUSD(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}
