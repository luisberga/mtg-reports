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
