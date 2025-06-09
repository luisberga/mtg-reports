package mocks

import (
	"context"
	"mtg-report/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type ConciliateRepositoryMock struct {
	mock.Mock
}

func NewConciliateRepositoryMock() *ConciliateRepositoryMock {
	return &ConciliateRepositoryMock{}
}

func (m *ConciliateRepositoryMock) InsertCardDetails(ctx context.Context, cardsDetails []domain.CardsDetails) error {
	args := m.Called(ctx, cardsDetails)
	return args.Error(0)
}

func (m *ConciliateRepositoryMock) GetCardsForUpdate(ctx context.Context, offset, limit int) ([]domain.Cards, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]domain.Cards), args.Error(1)
}
