package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type PriceServiceMock struct {
	mock.Mock
}

func NewPriceServiceMock() *PriceServiceMock {
	return &PriceServiceMock{}
}

func (m *PriceServiceMock) Conciliate(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
