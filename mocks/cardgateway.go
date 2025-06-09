package mocks

import (
	"context"
	"mtg-report/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type CardGatewayMock struct {
	mock.Mock
}

func NewCardGatewayMock() *CardGatewayMock {
	return &CardGatewayMock{}
}

func (m *CardGatewayMock) GetCardPrice(ctx context.Context, card domain.Cards) (float64, error) {
	args := m.Called(ctx, card)
	return args.Get(0).(float64), args.Error(1)
}
