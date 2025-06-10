package mocks

import (
	"context"
	"mtg-report/internal/core/domain"

	"github.com/stretchr/testify/mock"
)

type CardsRepositoryMock struct {
	mock.Mock
}

func NewCardsRepositoryMock() *CardsRepositoryMock {
	return &CardsRepositoryMock{}
}

func (c *CardsRepositoryMock) InsertCard(ctx context.Context, card domain.Cards) (domain.Cards, error) {
	args := c.Called(ctx, card)
	if args.Get(0) == nil {
		return domain.Cards{}, args.Error(1)
	}
	return args.Get(0).(domain.Cards), args.Error(1)
}

func (c *CardsRepositoryMock) InsertCards(ctx context.Context, cards []domain.Cards) error {
	args := c.Called(ctx, cards)
	return args.Error(0)
}

func (c *CardsRepositoryMock) GetCardbyID(ctx context.Context, id string) (domain.Cards, error) {
	args := c.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.Cards{}, args.Error(1)
	}
	return args.Get(0).(domain.Cards), args.Error(1)
}

func (c *CardsRepositoryMock) GetCards(ctx context.Context, filters map[string]string) ([]domain.Cards, error) {
	args := c.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cards), args.Error(1)
}

func (c *CardsRepositoryMock) DeleteCard(ctx context.Context, id string) error {
	args := c.Called(ctx, id)
	return args.Error(0)
}

func (c *CardsRepositoryMock) GetCardHistory(ctx context.Context, id string) ([]domain.Cards, error) {
	args := c.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cards), args.Error(1)
}

func (c *CardsRepositoryMock) UpdateCard(ctx context.Context, card domain.UpdateCard) (domain.Cards, error) {
	args := c.Called(ctx, card)
	if args.Get(0) == nil {
		return domain.Cards{}, args.Error(1)
	}
	return args.Get(0).(domain.Cards), args.Error(1)
}

func (c *CardsRepositoryMock) GetCardsPaginated(ctx context.Context, filters map[string]string, offset, limit int) ([]domain.Cards, error) {
	args := c.Called(ctx, filters, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cards), args.Error(1)
}

func (c *CardsRepositoryMock) GetCardsCount(ctx context.Context, filters map[string]string) (int64, error) {
	args := c.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (c *CardsRepositoryMock) GetCardHistoryPaginated(ctx context.Context, id string, offset, limit int) ([]domain.Cards, error) {
	args := c.Called(ctx, id, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cards), args.Error(1)
}

func (c *CardsRepositoryMock) GetCardHistoryCount(ctx context.Context, id string) (int64, error) {
	args := c.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (c *CardsRepositoryMock) GetCollectionStats(ctx context.Context) (domain.CollectionStats, error) {
	args := c.Called(ctx)
	return args.Get(0).(domain.CollectionStats), args.Error(1)
}
